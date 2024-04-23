// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

// Package test provides a method to export Ginkgo test results to custom unit test reporters.
package test

import (
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// RunSpecsWithReporters bootstraps Ginkgo/Gomega tests to function and results go to the /test/junit directory and log output
func RunSpecsWithReporters(t *testing.T, junitprefix string, suitename string) {

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, suitename)
}
