// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package kubernetesupgrade

import (
	"testing"

	policyv1beta1 "k8s.io/api/policy/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidateUserCreatedPodSecurityPolices(t *testing.T) {
	cases := []struct {
		name          string
		curr          string
		next          string
		policies      []policyv1beta1.PodSecurityPolicy
		errorExpected bool
	}{
		{
			"Upgrade from v1.23 to v1.24, no not validate",
			"1.23.0",
			"1.24.0",
			[]policyv1beta1.PodSecurityPolicy{},
			false,
		},
		{
			"Upgrade from v1.25, no not validate",
			"1.25.0",
			"1.26.0",
			[]policyv1beta1.PodSecurityPolicy{},
			false,
		},
		{
			"Upgrade to v1.25, only default policies present",
			"1.24.0",
			"1.25.0",
			[]policyv1beta1.PodSecurityPolicy{
				{ObjectMeta: v1.ObjectMeta{Name: "privileged"}},
				{ObjectMeta: v1.ObjectMeta{Name: "restricted"}},
			},
			false,
		},
		{
			"Upgrade to v1.25, default policies present plus extra",
			"1.24.0",
			"1.25.0",
			[]policyv1beta1.PodSecurityPolicy{
				{ObjectMeta: v1.ObjectMeta{Name: "privileged"}},
				{ObjectMeta: v1.ObjectMeta{Name: "restricted"}},
				{ObjectMeta: v1.ObjectMeta{Name: "extra"}},
			},
			true,
		},
		{
			"Upgrade to v1.25, non default policies present",
			"1.24.0",
			"1.25.0",
			[]policyv1beta1.PodSecurityPolicy{
				{ObjectMeta: v1.ObjectMeta{Name: "privileged"}},
				{ObjectMeta: v1.ObjectMeta{Name: "extra"}},
			},
			true,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			err := validateUserCreatedPodSecurityPolices(c.curr, c.next, c.policies)
			if err == nil && c.errorExpected {
				t.Fatal("expected validateUserCreatedPodSecurityPolices to return an error but it did not")
			} else if err != nil && !c.errorExpected {
				t.Fatalf("validateUserCreatedPodSecurityPolices not expected to return an error but it returned '%s'", err)
			}
		})
	}
}
