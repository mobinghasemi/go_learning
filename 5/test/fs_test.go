package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllFiles(t *testing.T) {
	assert.ElementsMatch(t,
		[]string{
			"README.md",
			"src/main.go",
			"src/fs/file1.go",
			"src/fs/file2.go"},
		dir.ListFilesRoot(),
	)
}

func TestDirFiles(t *testing.T) {
	actual, err := dir.ListFilesIn("src")
	assert.NoError(t, err)
	assert.ElementsMatch(t,
		[]string{
			"src/main.go",
			"src/fs/file1.go",
			"src/fs/file2.go",
		},
		actual,
	)
}

func TestCat(t *testing.T) {
	actual, err := dir.CatFile("README.md")
	assert.NoError(t, err)
	assert.Equal(t, "### MY PIP", actual)
}

func TestCatNotFound(t *testing.T) {
	_, err := dir.CatFile("src/non_existing_file")
	assert.Error(t, err)
}

func TestAppend(t *testing.T) {
	dir.AppendToFile("src/main.go", "package main")
	dir.AppendToFile("src/main.go", "\nfunc main(){}\n")
	actual, err := dir.CatFile("src/main.go")
	assert.NoError(t, err)
	assert.Equal(t, "package main\nfunc main(){}\n", actual)
}

func TestAppendNotFound(t *testing.T) {
	err := dir.AppendToFile("src/non_existing_file", "salam")
	assert.Error(t, err)
}

func TestWrite(t *testing.T) {
	dir.CreateFile(".gitignore")
	dir.WriteToFile(".gitignore", "*.out\n*.exe")
	actual, err := dir.CatFile(".gitignore")
	assert.NoError(t, err)
	assert.Equal(t, "*.out\n*.exe", actual)
}

func TestWriteNotFound(t *testing.T) {
	err := dir.WriteToFile("src/non_existing_file", "salam")
	assert.Error(t, err)
}

func TestClone(t *testing.T) {
	clonedWD := dir.Clone()
	clonedWD.AppendToFile("README.md", "junk content")
	clonedWD.CreateFile("LICENSE")

	actual, err := dir.CatFile("README.md")
	assert.NoError(t, err)
	assert.Equal(t, "### MY PIP", actual)
	assert.NotContains(t, dir.ListFilesRoot(), "LICENSE")
	assert.Contains(t, clonedWD.ListFilesRoot(), "LICENSE")
}
