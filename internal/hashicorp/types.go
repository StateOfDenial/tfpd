package hashicorp

import (
	"fmt"
)

type TerraformProvider struct {
	Name    string
	Version string
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
