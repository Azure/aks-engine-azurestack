// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package ssh

import (
	"os"
	"path/filepath"
)

func init() {
	khsource = filepath.Join(os.Getenv("HOME"), ".ssh", knownHostFileName)
	khpath = filepath.Join(os.Getenv("HOME"), ".aks-engine-azurestack", knownHostFileName)
	copyKnownHosts()
}
