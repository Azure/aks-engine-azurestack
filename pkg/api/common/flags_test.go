// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package common

import (
	"testing"
)

func TestRemoveFromCommaSeparatedList(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		toRemove []string
		expected string
	}{
		{
			"Single value to remove",
			"aks,engine,azure,stack,hub",
			[]string{"azure"},
			"aks,engine,stack,hub",
		},
		{
			"Multiple values to remove",
			"aks,engine,azure,stack,hub",
			[]string{"engine", "stack"},
			"aks,azure,hub",
		},
		{
			"Remove first value",
			"aks,engine,azure,stack,hub",
			[]string{"aks"},
			"engine,azure,stack,hub",
		},
		{
			"Remove last value",
			"aks,engine,azure,stack,hub",
			[]string{"hub"},
			"aks,engine,azure,stack",
		},
		{
			"Value to remove not in list",
			"aks,engine,azure,stack,hub",
			[]string{"foo", "bar"},
			"aks,engine,azure,stack,hub",
		},
		{
			"Input list has spaces",
			"aks, engine,azure , stack ,hub",
			[]string{"azure"},
			"aks,engine,stack,hub",
		},
		{
			"Input list has unexpected casing",
			"Aks,Engine,AZURE,Stack,Hub",
			[]string{"azure"},
			"Aks,Engine,Stack,Hub",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			actual := RemoveFromCommaSeparatedList(c.input, c.toRemove...)
			if actual != c.expected {
				t.Fatalf("expected removeFromCommaSeparatedList to return '%s', but instead got '%s", c.expected, actual)
			}
		})
	}
}
