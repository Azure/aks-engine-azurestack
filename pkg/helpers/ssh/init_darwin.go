package ssh

import (
	"os"
	"path/filepath"
)

func init() {
	khpath = filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	lineBreak = "\n"
}
