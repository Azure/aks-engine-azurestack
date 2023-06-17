package ssh

import (
	"os"
	"path/filepath"
)

func init() {
	khpath = filepath.Join(os.Getenv("UserProfile"), ".ssh", "known_hosts")
	lineBreak = "\r\n"
}
