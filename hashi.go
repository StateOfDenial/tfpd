package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
)

var baseUrl string

func discoverLockFile() (string, error) {
	d, _ := os.Getwd()
	for d != "/" {
		files, _ := os.ReadDir(d)
		for _, file := range files {
			if file.Name() == ".terraform.lock.hcl" {
				return filepath.Join(d, file.Name()), nil
			}
		}
		d = filepath.Dir(d)
	}
	return "", errors.New("could not find a lock file")
}

func getProvidersFromLockFile(filePath string) []string {
	providers := make([]string, 0)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	r, err := regexp.Compile("registry.terraform.io/(.*)\"\\s")
	if err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		if r.MatchString(scanner.Text()) {
			providers = append(providers, r.FindStringSubmatch(scanner.Text())[1])
		}
	}
	return providers
}

type Properties struct {
	Type  string
	Alias string
}

type Version struct {
	Attributes struct {
		Version string
	}
	Id string
}

func (v Version) String() string {
	return v.Attributes.Version
}

type Provider struct {
	Type       string
	Id         string
	Properties Properties
}

type ProviderRes struct {
	Data     Provider
	Links    map[string]interface{}
	Included []Version
}

type Resource struct {
	Type       string
	Id         string
	Attributes struct {
		Category string
		Slug     string
	}
	Links map[string]interface{}
}

func (r Resource) String() string {
	return fmt.Sprintf("%s: %s", r.Attributes.Category, r.Attributes.Slug)
}

type ProviderVersionRes struct {
	Data     Version
	Links    map[string]interface{}
	Included []Resource
}

type ResourceDoc struct {
	Attributes struct {
		Content string
	}
}

type ProviderDocRes struct {
	Data     ResourceDoc
	Links    map[string]interface{}
	Included []Resource
}

func callRegistryUrl(path string) []byte {
	url, err := url.Parse(baseUrl + path)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func getProviderId(provider string) string {
	body := callRegistryUrl("providers/" + provider)
	var result ProviderRes
	json.Unmarshal(body, &result)
	return string(result.Data.Id)
}

func getProviderVersions(provider string) ProviderRes {
	body := callRegistryUrl("providers/" + provider + "?include=provider-versions")
	var result ProviderRes
	json.Unmarshal(body, &result)
	return result
}

func getProviderVersionResources(version string) ProviderVersionRes {
	body := callRegistryUrl("provider-versions/" + version + "?include=provider-docs")
	var result ProviderVersionRes
	json.Unmarshal(body, &result)
	return result
}

func getResourceDoc(resource string) string {
	body := callRegistryUrl("provider-docs/" + resource)
	var result ProviderDocRes
	json.Unmarshal(body, &result)
	return result.Data.Attributes.Content
}
