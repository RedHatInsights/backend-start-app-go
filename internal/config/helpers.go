package config

import (
	"os"
	"path"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

func BinaryName() string {
	return path.Base(os.Args[0])
}

func InClowder() bool {
	return clowder.IsClowderEnabled()
}
