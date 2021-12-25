package constant

import (
	"os"
	"os/exec"
	"path/filepath"
)

var BasePath string

func initBasePath() {
	basePath, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}
	re, err := filepath.Abs(basePath)
	if err != nil {
		panic(err)
	}
	BasePath = re
}
