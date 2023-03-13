//go:build test
// +build test

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"log"
	"os/exec"

	"github.com/Azure/aks-engine-azurestack/test/e2e/kubernetes/util"
)

// Generate will run aks-engine generate on a given cluster definition
func (e *Engine) Generate() error {
	cmd := exec.Command("./bin/aks-engine-azurestack", "generate", e.Config.ClusterDefinitionTemplate, "--output-directory", e.Config.GeneratedDefinitionPath)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to generate aks-engine template with cluster definition - %s: %s\n", e.Config.ClusterDefinitionTemplate, err)
		log.Printf("Command:./bin/aks-engine-azurestack generate %s --output-directory %s\n", e.Config.ClusterDefinitionTemplate, e.Config.GeneratedDefinitionPath)
		log.Printf("Output:%s\n", out)
		return err
	}
	return nil
}

// Deploy will run aks-engine deploy on a given cluster definition
func (e *Engine) Deploy(location string) error {
	cmd := exec.Command("./bin/aks-engine-azurestack", "deploy",
		"--location", location,
		"--api-model", e.Config.ClusterDefinitionPath,
		"--dns-prefix", e.Config.DefinitionName,
		"--output-directory", e.Config.GeneratedDefinitionPath,
		"--resource-group", e.Config.DefinitionName,
	)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to deploy aks-engine template with cluster definition - %s: %s\n", e.Config.ClusterDefinitionTemplate, err)
		log.Printf("Output:%s\n", out)
		return err
	}
	return nil
}

// Upgrade will run aks-engine upgrade on a given cluster definition
func (e *Engine) Upgrade(location string, upgradeVersion string) error {
	cmd := exec.Command("./bin/aks-engine-azurestack", "upgrade",
		"--location", location,
		"--api-model", e.Config.GeneratedApiModelPath,
		"--resource-group", e.Config.DefinitionName,
		"--upgrade-version", upgradeVersion,
		"--vm-timeout", "20",
	)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to upgrade aks-engine template with cluster definition - %s: %s\n", e.Config.ClusterDefinitionTemplate, err)
		log.Printf("Output:%s\n", out)
		return err
	}
	return nil
}
