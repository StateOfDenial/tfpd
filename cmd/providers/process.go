package providers

import (
	"errors"
	"fmt"
	// "os"

	h "github.com/StateOfDenial/tfpd/internal/hashicorp"
	"github.com/StateOfDenial/tfpd/internal/tui"
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

func processLockFile() ([]h.TerraformProvider, error) {
	var lockProviders []h.TerraformProvider
	lockFile, err := h.DiscoverLockFile()
	if err != nil {
		return lockProviders, errors.New("Could not find a lock file")
	}
	lockProviders = h.GetProvidersFromLockFile(lockFile)
	if len(lockProviders) == 0 {
		return lockProviders, errors.New("Found no searchable providers in lock file")
	}
	return lockProviders, nil
}

func processProviders(providers []h.TerraformProvider, initProvider string) (string, string, error) {
	switch {
	case len(providers) > 1:
		provs := make([]string, len(providers))
		for i := range providers {
			provs[i] = providers[i].Name + " : " + providers[i].Version
		}
		ff := tui.NewFuzzyFinder()
		ff.SetFuzzyItems(provs)
		providerIdx := ff.FuzzyFindWithInput(initProvider)
		return providers[providerIdx].Name, providers[providerIdx].Version, nil
	case len(providers) == 1:
		return providers[0].Name, providers[0].Version, nil
	default:
		return "", "", errors.New("No providers list to process")
	}
}

func command(provider, version, resource string, isData bool) error {
	var providerName, providerVersion string
	hashiClient := h.NewClient().SetBaseUrl("http://registry.terraform.io/v2/")
	providers, err := processLockFile()
	if err != nil {
		if provider == "" {
			fmt.Println()
			providerName = promptForInput("Enter in a provider to look for: e.g. 'hashicorp/google'")
		} else {
			providerName = provider
		}
	} else {
		providerName, providerVersion, err = processProviders(providers, provider)
		if err != nil {
			panic("tried to process 0 providers from lock file, shouldn't happen")
		}
	}

	versions := hashiClient.GetProviderVersions(providerName)
	var providerVersionResources h.ProviderVersionRes
	if providerVersion == "" {
		vers := make([]string, len(versions.Included))
		for i := range versions.Included {
			vers[i] = versions.Included[i].String()
		}
		ff := tui.NewFuzzyFinder()
		ff.SetFuzzyItems(vers)
		versionIdx := ff.FuzzyFindWithInput(version)
		providerVersionResources = hashiClient.GetProviderVersionResources(versions.Included[versionIdx].Id)
	} else {
		filterTest := func(v h.Version) bool { return v.Attributes.Version == providerVersion }
		version := filter(versions.Included, filterTest)
		providerVersionResources = hashiClient.GetProviderVersionResources(version[0].Id)
	}

	resources := make([]string, len(providerVersionResources.Included))
	for i := range providerVersionResources.Included {
		resources[i] = providerVersionResources.Included[i].String()
	}
	ff := tui.NewFuzzyFinder()
	ff.SetFuzzyItems(resources)
	var resourceIdx int
	if isData {
		resourceIdx = ff.FuzzyFindWithInput("data-sources: " + resource)
	} else {
		resourceIdx = ff.FuzzyFindWithInput(resource)
	}

	m := tui.NewMDViewer(hashiClient.GetResourceDoc(providerVersionResources.Included[resourceIdx].Id))
	m.Display()
	return nil
}
