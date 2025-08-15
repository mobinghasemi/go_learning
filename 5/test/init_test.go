package main

import (
	"fmt"
	"os"
	"pip/commands"
	"pip/fs"
	"testing"
)

var (
	dir  *fs.Dir
	gopi commands.GOPI
	pip  *commands.PIP
)

type LocalGOPI struct {
	data map[string]*fs.Dir
}

func (gopi *LocalGOPI) Get(pkgName string) (*fs.Dir, error) {
	dir, prs := gopi.data[pkgName]
	if !prs {
		return nil, fmt.Errorf("404: package %s cannot found in GOPI", pkgName)
	}

	return dir.Clone(), nil
}

func TestMain(m *testing.M) {
	Setup()
	os.Exit(m.Run())
}

func generateProject(deps ...string) *fs.Dir {
	req := "requirements.txt"
	dir := fs.MkDir()
	dir.CreateFile(req)
	for _, dep := range deps {
		dir.AppendToFile(req, fmt.Sprintf("%s\n", dep))
	}

	dir.CreateDir("src")
	dir.CreateFile("src/main.go")
	dir.AppendToFile("src/main.go", `// TODO: implement`)

	return dir
}

func Setup() {
	dir = fs.MkDir()

	err := dir.CreateFile("README.md")
	if err != nil {
		panic(err)
	}

	err = dir.CreateDir("src")
	if err != nil {
		panic(err)
	}

	err = dir.CreateFile("src/main.go")
	if err != nil {
		panic(err)
	}

	err = dir.CreateDir("src/fs")
	if err != nil {
		panic(err)
	}

	err = dir.CreateFile("src/fs/file1.go")
	if err != nil {
		panic(err)
	}

	err = dir.CreateFile("src/fs/file2.go")
	if err != nil {
		panic(err)
	}

	err = dir.WriteToFile("README.md", "### MY PIP")
	if err != nil {
		panic(err)
	}

	gopi = &LocalGOPI{
		data: map[string]*fs.Dir{
			"echo":                          generateProject("jwt", "testify", "fasttemplate"),
			"jwt":                           generateProject(),
			"testify":                       generateProject("go-spew", "go-difflib"),
			"fasttemplate":                  generateProject("bytebufferpool"),
			"go-spew":                       generateProject(),
			"bytebufferpool":                generateProject(),
			"go-difflib":                    generateProject(),
			"prj-with-indirect-invalid-dep": generateProject("prj-with-invalid-dep"),
			"prj-with-invalid-dep":          generateProject("invalid-dep"),
		},
	}

	pip = commands.NewPIP(fs.MkDir(), gopi)
}
