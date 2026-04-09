package hashicorp

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func DiscoverLockFile() (string, error) {
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

func GetProvidersFromLockFile(filePath string) []TerraformProvider {
	providers := make([]TerraformProvider, 0)
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

	r2, err := regexp.Compile("version\\s*=\\s*\"(.*)\"")
	if err != nil {
		log.Fatal(err)
	}

	provider := ""

	for scanner.Scan() {
		if provider == "" && r.MatchString(scanner.Text()) {
			provider = r.FindStringSubmatch(scanner.Text())[1]
		} else if provider != "" && r2.MatchString(scanner.Text()) {
			providers = append(providers, TerraformProvider{
				Name:    provider,
				Version: r2.FindStringSubmatch(scanner.Text())[1],
			})
			provider = ""
		}

	}

	return providers
}
