package main

import (
	"bytes"
	"errors"
	"flag"
	"github.com/zeddy-go/zeddy/slicex"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

var src = flag.String("s", ".", "source path")
var dest = flag.String("d", "", "dest path")
var modName = flag.String("n", "", "module name, default dest dir name")
var toTemplate = flag.Bool("toTemplate", false, "project convert to template")
var fromTemplate = flag.Bool("fromTemplate", false, "template convert to project")

var ignores = []string{".git", "go.sum", "go.mod"}

func main() {
	flag.Parse()

	if *src == "" {
		panic(errors.New("source is required"))
	}
	src, err := filepath.Abs(*src)
	if err != nil {
		panic(err)
	}

	if *dest == "" {
		panic(errors.New("dest is required"))
	}
	dest, err := filepath.Abs(*dest)
	if err != nil {
		panic(err)
	}

	if *toTemplate {
		toTemp(src, dest)
	} else if *fromTemplate {
		fromTemp(src, dest)
	} else {
		normalCopy(src, dest)
	}
}

func fromTemp(src string, dest string) {
	destName := filepath.Base(dest)
	if *modName != "" {
		destName = *modName
	}
	createEmptyModule(dest, destName)

	params := struct {
		Module string
	}{
		Module: destName,
	}

	filepath.Walk(src, func(path string, info fs.FileInfo, e error) (err error) {
		if e != nil {
			panic(e)
		}
		if slicex.Contains(filepath.Base(path), ignores) {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return
			}
		}

		dest := strings.Replace(path, src, dest, 1)
		if info.IsDir() {
			err = os.MkdirAll(dest, os.ModeDir)
			if err != nil {
				panic(err)
			}
		} else {
			content := readFile(path)
			t, err := template.New(src).Parse(string(content))
			if err != nil {
				panic(err)
			}
			buffer := new(bytes.Buffer)
			err = t.Execute(buffer, params)
			if err != nil {
				panic(err)
			}

			content = buffer.Bytes()
			dest = strings.Replace(dest, ".tpl", "", 1)
			writeFile(dest, content)
		}

		return
	})
}

func toTemp(src string, dest string) {
	srcGoMod := readModFile(src + "/go.mod")
	srcName := srcGoMod.Module.Mod.Path

	filepath.Walk(src, func(path string, info fs.FileInfo, e error) (err error) {
		if e != nil {
			panic(e)
		}
		if slicex.Contains(filepath.Base(path), ignores) {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return
			}
		}

		dest := strings.Replace(path, src, dest, 1)
		if info.IsDir() {
			err = os.MkdirAll(dest, os.ModeDir)
			if err != nil {
				panic(err)
			}
		} else {
			content := readFile(path)
			switch filepath.Ext(path) {
			case ".go":
				content = replaceGoFile(path, content, "{{.Module}}", srcName)
			}
			dest += ".tpl"
			writeFile(dest, content)
		}

		return
	})
}

func normalCopy(src string, dest string) {
	srcGoMod := readModFile(src + "/go.mod")
	srcName := srcGoMod.Module.Mod.Path
	destName := filepath.Base(dest)
	if *modName != "" {
		destName = *modName
	}

	createEmptyModule(dest, destName)

	filepath.Walk(src, func(path string, info fs.FileInfo, e error) (err error) {
		if e != nil {
			panic(e)
		}
		if slicex.Contains(filepath.Base(path), ignores) {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return
			}
		}

		dest := strings.Replace(path, src, dest, 1)
		if info.IsDir() {
			err = os.MkdirAll(dest, os.ModeDir)
			if err != nil {
				panic(err)
			}
		} else {
			content := readFile(path)
			switch filepath.Ext(path) {
			case ".go":
				content = replaceGoFile(path, content, destName, srcName)
			}
			writeFile(dest, content)
		}

		return
	})
}

func readFile(path string) (content []byte) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	content, err = io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return
}

func readModFile(path string) (f *modfile.File) {
	content := readFile(path)
	f, err := modfile.Parse(path, content, nil)
	if err != nil {
		panic(err)
	}
	return
}

func writeFile(path string, content []byte) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		panic(err)
	}
}

func createEmptyModule(path string, name string) {
	err := os.MkdirAll(path, os.ModeDir)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("go", "mod", "init", name)
	cmd.Dir = path
	err = cmd.Run()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err != nil {
		panic(err)
	}
}

func replaceGoFile(path string, content []byte, destName string, srcName string) []byte {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, content, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	for _, item := range f.Imports {
		if *toTemplate {
			item.Path.Value = strings.Replace(item.Path.Value, srcName, "{{.Module}}", 1)
		} else {
			item.Path.Value = strings.Replace(item.Path.Value, srcName, destName, 1)
		}
	}

	var tmp []byte
	buffer := bytes.NewBuffer(tmp)
	err = format.Node(buffer, fset, f)
	return buffer.Bytes()
}
