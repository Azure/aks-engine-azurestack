// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import (
	"testing"

	"github.com/Azure/aks-engine-azurestack/pkg/api/common"
	"github.com/blang/semver"
	. "github.com/onsi/gomega"
)

func TestInvalidVersion(t *testing.T) {
	RegisterTestingT(t)

	invalid := []string{
		"invalid number",
		"invalid.number",
		"a4.b7.c3",
		"31.29.",
		".17.02",
		"43.156.89.",
		"1.2.a"}

	for _, v := range invalid {
		_, e := semver.Make(v)
		Expect(e).NotTo(BeNil())
	}
}

func TestVersionCompare(t *testing.T) {
	RegisterTestingT(t)

	type record struct {
		v1, v2    string
		isGreater bool
	}
	records := []record{
		{"37.48.59", "37.48.59", false},
		{"17.4.5", "3.1.1", true},
		{"9.6.5", "9.45.5", false},
		{"2.3.8", "2.3.24", false}}

	for _, r := range records {
		ver, e := semver.Make(r.v1)
		Expect(e).To(BeNil())
		constraint, e := semver.Make(r.v2)
		Expect(e).To(BeNil())
		Expect(r.isGreater).To(Equal(ver.GT(constraint)))
	}
}

func TestOrchestratorUpgradeInfo(t *testing.T) {
	RegisterTestingT(t)
	testVersions := []string{"1.6.9", "1.7.0", "1.7.15", "1.8.4", "1.9.6", "1.10.0-beta.2", "1.11.0", "1.12.0"}
	for _, deployedVersion := range testVersions {
		csOrch := &OrchestratorProfile{
			OrchestratorType:    Kubernetes,
			OrchestratorVersion: deployedVersion,
		}
		v, e := getKubernetesAvailableUpgradeVersions(deployedVersion, common.GetAllSupportedKubernetesVersions(false, false, false))
		Expect(e).To(BeNil())
		orch, e := GetOrchestratorVersionProfile(csOrch, false, false)
		Expect(e).To(BeNil())
		Expect(len(orch.Upgrades)).To(Equal(len(v)))
	}

	// The latest version is not upgradable
	csOrch := &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: common.GetMaxVersion(common.GetAllSupportedKubernetesVersions(false, false, false), true),
	}
	orch, e := GetOrchestratorVersionProfile(csOrch, false, false)
	Expect(e).To(BeNil())
	Expect(len(orch.Upgrades)).To(Equal(0))
}

func TestGetOrchestratorVersionProfileList(t *testing.T) {
	RegisterTestingT(t)
	// kubernetes only
	list, e := GetOrchestratorVersionProfileList(common.Kubernetes, "", false, "")
	Expect(e).To(BeNil())
	Expect(len(list)).To(Equal(len(common.GetAllSupportedKubernetesVersions(false, false, false))))
	for _, v := range list {
		Expect(common.GetAllSupportedKubernetesVersions(false, false, false)).Should(ContainElement(v.OrchestratorProfile.OrchestratorVersion))
	}
}

func TestKubernetesInfo(t *testing.T) {
	RegisterTestingT(t)

	invalid := []string{
		"invalid number",
		"invalid.number",
		"a4.b7.c3",
		"31.29.",
		".17.02",
		"43.156.89.",
		"1.2.a",
		"1.5.9",
		"1.6.8"}

	for _, v := range invalid {
		csOrch := &OrchestratorProfile{
			OrchestratorType:    Kubernetes,
			OrchestratorVersion: v,
		}

		_, e := kubernetesInfo(csOrch, false, false)
		Expect(e).NotTo(BeNil())
	}

}

func TestGetKubernetesAvailableUpgradeVersions(t *testing.T) {
	RegisterTestingT(t)
	cases := []struct {
		name             string
		version          string
		versions         []string
		expectedUpgrades []string
	}{
		{
			name:             "from 1.7.15",
			version:          "1.7.15",
			versions:         []string{"1.9.10", "1.9.11", "1.10.3", "1.10.4", "1.11.3", "1.11.4", "1.12.0-alpha.1"},
			expectedUpgrades: []string{"1.9.10", "1.9.11"},
		},
		{
			name:             "from 1.8.14",
			version:          "1.8.14",
			versions:         []string{"1.7.15", "1.8.14", "1.8.15", "1.9.10", "1.9.11", "1.10.3", "1.10.4"},
			expectedUpgrades: []string{"1.8.15", "1.9.10", "1.9.11"},
		},
		{
			name:             "from 1.8.14 with alternate supported versions",
			version:          "1.8.14",
			versions:         []string{"1.9.10", "1.9.11", "1.10.3", "1.10.4", "1.11.3", "1.11.4", "1.12.0-alpha.1"},
			expectedUpgrades: []string{"1.9.10", "1.9.11"},
		},
		{
			name:             "from 1.9.10",
			version:          "1.9.10",
			versions:         []string{"1.9.10", "1.9.11", "1.10.3", "1.10.4", "1.11.3", "1.11.4", "1.12.0-alpha.1"},
			expectedUpgrades: []string{"1.9.11", "1.10.3", "1.10.4"},
		},
		{
			name:             "from 1.10.4",
			version:          "1.10.4",
			versions:         []string{"1.9.10", "1.9.11", "1.10.3", "1.10.4", "1.11.3", "1.11.4", "1.12.0-alpha.1"},
			expectedUpgrades: []string{"1.11.3", "1.11.4"},
		},
		{
			name:             "from 1.12.1",
			version:          "1.12.1",
			versions:         []string{"1.9.10", "1.9.11", "1.10.3", "1.10.4", "1.11.3", "1.11.4", "1.12.1", "1.12.2"},
			expectedUpgrades: []string{"1.12.2"},
		},
		{
			name:             "no supported upgrades",
			version:          "1.12.2",
			versions:         []string{"1.9.10", "1.9.11", "1.10.3", "1.10.4", "1.11.3", "1.11.4", "1.12.1", "1.12.2"},
			expectedUpgrades: []string{},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			upgrades, err := getKubernetesAvailableUpgradeVersions(c.version, c.versions)
			Expect(err).To(BeNil())
			Expect(upgrades).To(Equal(c.expectedUpgrades))
		})
	}
}
