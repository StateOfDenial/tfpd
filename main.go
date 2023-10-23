package main

import (
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

func main() {
	baseUrl = "http://registry.terraform.io/v2/"
	lockFile, err := discoverLockFile()
	var providerName, providerVersion string
	var ff FuzzyFinder
	if err != nil {
		// Fall back to prompting for provider
		providerName = promptForInput("Enter in a provider to look for: e.g. 'hashicorp/google'")
	} else {
		lockProviders := getProvidersFromLockFile(lockFile)

		providers := make([]string, len(lockProviders))
		for i := range lockProviders {
			providers[i] = lockProviders[i].Name + " : " + lockProviders[i].Version
		}
		var providerIdx int
		switch {
		case len(providers) > 1:
			ff = NewFuzzyFinder()
			ff.SetFuzzyItems(providers)
			providerIdx = ff.FuzzyFind()
			providerName = lockProviders[providerIdx].Name
			providerVersion = lockProviders[providerIdx].Version
		case len(providers) == 1:
			providerIdx = 0
			providerName = lockProviders[providerIdx].Name
			providerVersion = lockProviders[providerIdx].Version
		default:
			providerName = promptForInput("Enter in a provider to look for: e.g. 'hashicorp/google'")
		}
	}

	versions := getProviderVersions(providerName)
	var providerVersionResources ProviderVersionRes
	if providerVersion == "" {
		vers := make([]string, len(versions.Included))
		for i := range versions.Included {
			vers[i] = versions.Included[i].String()
		}
		ff = NewFuzzyFinder()
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
	ff = NewFuzzyFinder()
	ff.SetFuzzyItems(resources)
	resourceIdx := ff.FuzzyFind()

	fmt.Fprint(os.Stdout, getResourceDoc(providerVersionResources.Included[resourceIdx].Id))
}
