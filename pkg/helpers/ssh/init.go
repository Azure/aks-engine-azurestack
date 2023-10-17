// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package ssh

import (
	"io"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

var (
	khsource string
	khpath   string
)

const (
	knownHostFileName = "known_hosts"
)

// copyKnownHosts copies the SSH's known_hosts file to aks-engine's working directory.
// Errors are logged but not returned as we can still continue regardless.
func copyKnownHosts() {
	src, err := os.Open(khsource)
	if err != nil {
		log.Warnf("failed to open file %s", khsource)
		return
	}
	defer src.Close()

	dir := path.Dir(khpath)
	if err = os.MkdirAll(dir, 0700); err != nil {
		log.Warnf("failed to create %s directory", dir)
		return
	}

	dst, err := os.Create(khpath)
	if err != nil {
		log.Warnf("failed to open file %s", khpath)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		log.Warnf("failed to copy %s to %s", khsource, khpath)
		return
	}
}
