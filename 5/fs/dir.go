package fs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	cp "github.com/otiai10/copy"
)

type Dir struct {
	d string
}

func MkDir() *Dir {
	dir, err := os.MkdirTemp("", "vc")
	if err != nil {
		panic(err)
	}
	return &Dir{
		d: dir + string(filepath.Separator),
	}
}

func (d *Dir) CreateFile(filename string) error {
	f, err := os.Create(d.d + filename)
	if err != nil {
		return errors.New("cannot create a file")
	}
	defer f.Close()
	return nil
}

func (d *Dir) CreateDir(dirname string) error {
	err := os.Mkdir(d.d+dirname, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dir) ListFilesIn(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(d.d+dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			myPath := strings.TrimPrefix(path, d.d)
			myPath = strings.ReplaceAll(myPath, "\\", "/") // Normalize path separators
			files = append(files, myPath)
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("cannot list files in dir")
	}
	return files, nil
}

func (d *Dir) ListFilesRoot() []string {
	res, err := d.ListFilesIn("")
	if err != nil {
		panic(err)
	}
	return res
}

func (d *Dir) CatFile(file string) (string, error) {
	content, err := os.ReadFile(d.d + file)
	if err != nil {
		return "", errors.New("cannot read the file")
	}
	return string(content), nil
}

func Contains(list []string, elem string) bool {
	for _, v := range list {
		if v == elem {
			return true
		}
	}
	return false
}

func (d *Dir) WriteToFile(file string, content string) error {
	if !Contains(d.ListFilesRoot(), file) {
		return errors.New("files does not exist")
	}
	err := os.WriteFile(d.d+file, []byte(content), os.ModePerm)
	if err != nil {
		panic(err)
	}
	return nil
}

func (d *Dir) AppendToFile(file, content string) error {
	if !Contains(d.ListFilesRoot(), file) {
		return errors.New("file does not exist")
	}
	f, err := os.OpenFile(d.d+file, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		panic(err)
	}

	return nil
}

func (d *Dir) Clone() *Dir {
	new_Dir, err := os.MkdirTemp("", "vc")
	if err != nil {
		panic(err)
	}
	cwd := &Dir{
		d: new_Dir + "/",
	}
	cp.Copy(d.d, cwd.d)
	return cwd
}

func (d *Dir) Mount(path string, other *Dir) error {
	if err := d.CreateDir(path); err != nil {
		return err
	}
	if err := cp.Copy(other.d, d.d+path); err != nil {
		return err
	}
	return nil

}

func (d *Dir) Remove(path string) error {
	return os.RemoveAll(d.d + path)
}
