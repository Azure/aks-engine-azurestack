// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package common

import "strings"

// RemoveFromCommaSeparatedList excludes from the input comma-separated list the strings toRemove
func RemoveFromCommaSeparatedList(input string, toRemove ...string) string {
	removeMap := map[string]bool{}
	for _, remove := range toRemove {
		removeKey := strings.ToLower(strings.TrimSpace(remove))
		removeMap[removeKey] = true
	}
	ret := []string{}
	for _, value := range strings.Split(input, ",") {
		key := strings.TrimSpace(value)
		if !removeMap[strings.ToLower(key)] {
			ret = append(ret, key)
		}
	}
	return strings.Join(ret, ",")
}
