package hashicorp

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

var baseUrl string

type Client struct {
	baseUrl string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) SetBaseUrl(baseUrl string) *Client {
	c.baseUrl = baseUrl
	return c
}

func (c Client) callRegistryUrl(path string) []byte {
	url, err := url.Parse(c.baseUrl + path)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	log.Print(res.Status)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func (c *Client) GetProviderId(provider string) string {
	body := c.callRegistryUrl("providers/" + provider)
	var result ProviderRes
	json.Unmarshal(body, &result)
	return string(result.Data.Id)
}

func (c *Client) GetProviderVersions(provider string) ProviderRes {
	body := c.callRegistryUrl("providers/" + provider + "?include=provider-versions")
	var result ProviderRes
	json.Unmarshal(body, &result)
	return result
}

func (c *Client) GetProviderVersionResources(version string) ProviderVersionRes {
	body := c.callRegistryUrl("provider-versions/" + version + "?include=provider-docs")
	var result ProviderVersionRes
	json.Unmarshal(body, &result)
	return result
}

func (c *Client) GetResourceDoc(resource string) string {
	body := c.callRegistryUrl("provider-docs/" + resource)
	var result ProviderDocRes
	json.Unmarshal(body, &result)
	return result.Data.Attributes.Content
}
