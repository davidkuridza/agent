// +build rhel centos

package main

import (
	"fmt"
	"regexp"
	"strings"
)

type mappedPackages map[string]*Package

func GetPackages() ([]*Package, error) {
	return readPackages()
}

// Fetches packages with `yum info installed`
func readPackages() ([]*Package, error) {
	results := []*Package{}

	b, err := execCommand("yum", "info", "installed")
	if err != nil {
		return results, err
	}

	// split response by packages
	packages := strings.Split(string(b), "\n\n")

	// loop through packages
	reKV := regexp.MustCompile("(Name|Arch|Version|Release|From repo|Epoch)\\s+: (.*)")
	for _, pkg := range packages {
		m := reKV.FindAllStringSubmatch(pkg, -1)

		// put results into a key value map
		matches := make(map[string]string)
		for _, v := range m {
			matches[v[1]] = v[2]
		}

		// skip empty lines
		if matches["Name"] == "" {
			continue
		}

		results = append(results, buildPackage(matches))
	}

	return results, nil
}

func buildVersion(matches map[string]string) string {
	version := ""

	if epoch, ok := matches["Epoch"]; ok && epoch != "" {
		version = fmt.Sprintf("%s:", epoch)
	}

	return fmt.Sprintf("%s%s-%s", version, matches["Version"], matches["Release"])
}