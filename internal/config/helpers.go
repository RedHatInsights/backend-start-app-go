package config

import (
	"os"
	"path"
)

func BinaryName() string {
	return path.Base(os.Args[0])
}
