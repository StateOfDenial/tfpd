package main

import (
	"fmt"
	"os"
)

func main() {
	baseUrl = "http://registry.terraform.io/v2/"
	lockFile, err := discoverLockFile()
	if err != nil {
		panic(err)
	}
	providers := getProvidersFromLockFile(lockFile)
	var providerIdx int
	var ff FuzzyFinder
	switch {
	case len(providers) > 1:
		ff = NewFuzzyFinder()
		ff.SetFuzzyItems(providers)
		providerIdx = ff.FuzzyFind()
	case len(providers) == 1:
		providerIdx = 0
	default:
		providerIdx = 0
		fmt.Println("oops")
	}
	versions := getProviderVersions(providers[providerIdx])
	ff = NewFuzzyFinder()
	vers := make([]string, len(versions.Included))
	for i := range versions.Included {
		vers[i] = versions.Included[i].String()
	}
	ff.SetFuzzyItems(vers)
	versionIdx := ff.FuzzyFind()

	providerVersionResources := getProviderVersionResources(versions.Included[versionIdx].Id)
	resources := make([]string, len(providerVersionResources.Included))
	for i := range providerVersionResources.Included {
		resources[i] = providerVersionResources.Included[i].String()
	}

	ff = NewFuzzyFinder()
	ff.SetFuzzyItems(resources)
	resourceIdx := ff.FuzzyFind()

	fmt.Fprint(os.Stdout, getResourceDoc(providerVersionResources.Included[resourceIdx].Id))
}
