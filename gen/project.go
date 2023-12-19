package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func createGoFile(path string, checkExists bool) (err error) {
	if checkExists {
		_, err = os.Stat(path)
		if err != nil {
			return
		}
	}

	path, _ = filepath.Abs(path)
	dirPath := filepath.Dir(path)
	pkgName := filepath.Base(dirPath)

	err = os.MkdirAll(dirPath, os.ModeDir)
	if err != nil {
		return
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	_, err = f.WriteString(fmt.Sprintf("package %s\n", pkgName))
	return
}

func createEmptyProject(name string, path ...string) (err error) {
	p, _ := filepath.Abs(".")
	if len(path) > 0 {
		p, err = filepath.Abs(path[0])
		if err != nil {
			return
		}
	}
	err = os.MkdirAll(p, os.ModeDir)
	if err != nil {
		return
	}
	f, err := os.OpenFile(p+"/go.mod", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("module %s\n\ngo %s\n\n", name, goVersion()))
	if err != nil {
		return
	}

	return
}

func goVersion() string {
	return strings.Replace(runtime.Version(), "go", "", 1)
}
