package main

import (
	"errors"
	"fmt"
	"os"
)

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return ret
}

func promptForInput(prompt string) string {
	var ret string
	fmt.Println(prompt)
	fmt.Scanln(&ret)
	return ret
}

func processLockFile() ([]TerraformProvider, error) {
	var lockProviders []TerraformProvider
	lockFile, err := discoverLockFile()
	if err != nil {
		return lockProviders, errors.New("Could not find a lock file")
	}
	lockProviders = getProvidersFromLockFile(lockFile)
	if len(lockProviders) == 0 {
		return lockProviders, errors.New("Found no searchable providers in lock file")
	}
	return lockProviders, nil
}

func processProviders(providers []TerraformProvider) (string, string, error) {
	switch {
	case len(providers) > 1:
		provs := make([]string, len(providers))
		for i := range providers {
			provs[i] = providers[i].Name + " : " + providers[i].Version
		}
		ff := NewFuzzyFinder()
		ff.SetFuzzyItems(provs)
		providerIdx := ff.FuzzyFind()
		return providers[providerIdx].Name, providers[providerIdx].Version, nil
	case len(providers) == 1:
		return providers[0].Name, providers[0].Version, nil
	default:
		return "", "", errors.New("No providers list to process")
	}
}

func main() {
	var providerName, providerVersion string
	baseUrl = "http://registry.terraform.io/v2/"
	providers, err := processLockFile()
	if err != nil {
		fmt.Println()
		providerName = promptForInput("Enter in a provider to look for: e.g. 'hashicorp/google'")
	} else {
		providerName, providerVersion, err = processProviders(providers)
		if err != nil {
			panic("tried to process 0 providers from lock file, shouldn't happen")
		}
	}

	versions := getProviderVersions(providerName)
	var providerVersionResources ProviderVersionRes
	if providerVersion == "" {
		vers := make([]string, len(versions.Included))
		for i := range versions.Included {
			vers[i] = versions.Included[i].String()
		}
		ff := NewFuzzyFinder()
		ff.SetFuzzyItems(vers)
		versionIdx := ff.FuzzyFind()
		providerVersionResources = getProviderVersionResources(versions.Included[versionIdx].Id)
	} else {
		filterTest := func(v Version) bool { return v.Attributes.Version == providerVersion }
		version := filter(versions.Included, filterTest)
		providerVersionResources = getProviderVersionResources(version[0].Id)
	}

	resources := make([]string, len(providerVersionResources.Included))
	for i := range providerVersionResources.Included {
		resources[i] = providerVersionResources.Included[i].String()
	}
	ff := NewFuzzyFinder()
	ff.SetFuzzyItems(resources)
	resourceIdx := ff.FuzzyFind()

	fmt.Fprint(os.Stdout, getResourceDoc(providerVersionResources.Included[resourceIdx].Id))
}
