// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package main

import (
	"os"

	"github.com/Azure/aks-engine-azurestack/cmd"
	colorable "github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStderr())
	log.SetOutput(colorable.NewColorableStdout())
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
