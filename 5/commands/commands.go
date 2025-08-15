package commands

// you can import "github.com/lithammer/fuzzysearch/fuzzy"

import (
	"bufio"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"pip/fs"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type GOPI interface {
	Get(string) (*fs.Dir, error)
}

type PIP struct {
	installDir    *fs.Dir
	gopi          GOPI
	userInstalled []string
	allInstalled  []string
	onceInstalled map[string]bool
}

func NewPIP(dir *fs.Dir, gopi GOPI) *PIP {
	return &PIP{
		installDir:    dir,
		gopi:          gopi,
		userInstalled: nil,
		allInstalled:  nil,
		onceInstalled: make(map[string]bool),
	}
}

func parseRequirementsTXT(reqTXTContent string) []string {
	result := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(reqTXTContent))
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
	return result
}

func (pip *PIP) DirectDeps(pkgName string) ([]string, error) {
	pkDir, err := pip.gopi.Get(pkgName)
	if err != nil {
		return nil, err
	}
	reqs, err := pkDir.CatFile("requirements.txt")
	if err != nil {
		return nil, errors.New("invalid project dependencies")
	}
	return parseRequirementsTXT(reqs), nil
}

func AddAllToMap(m map[string]bool, slice []string) {
	for _, v := range slice {
		m[v] = true
	}
}

func MapToSlice(m map[string]bool) []string {
	res := make([]string, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}

func (pip *PIP) AllDeps(pkgName string) ([]string, error) {
	result := make(map[string]bool)
	direct, err := pip.DirectDeps(pkgName)
	if err != nil {
		return nil, err
	}
	AddAllToMap(result, direct)
	for _, dep := range direct {
		dep_deps, err := pip.AllDeps(dep)
		if err != nil {
			return nil, err
		}
		AddAllToMap(result, dep_deps)
	}
	return MapToSlice(result), nil
}

func Contains(list []string, elem string) bool {
	for _, v := range list {
		if v == elem {
			return true
		}
	}
	return false
}

func (pip *PIP) CopyFromGopi(pkgName string) error {
	pip.onceInstalled[pkgName] = true
	if !Contains(pip.allInstalled, pkgName) {
		pip.allInstalled = append(pip.allInstalled, pkgName)
	} else {
		return nil
	}
	dl, err := pip.gopi.Get(pkgName)
	if err != nil {
		return err
	}
	err = pip.installDir.Mount(pkgName, dl)
	if err != nil {
		return err
	}
	return nil
}

func (pip *PIP) Install(pkgNames ...string) error {
	for _, pkgName := range pkgNames {
		_, err := pip.AllDeps(pkgName)
		if err != nil {
			return err
		}
	}
	for _, pkgName := range pkgNames {
		allDeps, err := pip.AllDeps(pkgName)
		if err != nil {
			return err
		}
		err = pip.CopyFromGopi(pkgName)
		if err != nil {
			return err
		}
		for _, v := range allDeps {
			err := pip.CopyFromGopi(v)
			if err != nil {
				return err
			}
		}
		if !Contains(pip.userInstalled, pkgName) {
			pip.userInstalled = append(pip.userInstalled, pkgName)
		}
	}
	return nil
}

func (pip *PIP) InstallR(dir *fs.Dir, reqFile string) error {
	reqs, err := dir.CatFile(reqFile)
	if err != nil {
		return err
	}
	allPkgs := parseRequirementsTXT(reqs)
	err = pip.Install(allPkgs...)
	if err != nil {
		return err
	}
	return nil
}

func (pip *PIP) AllUserInstalledPackages() []string {
	result := make([]string, 0, len(pip.userInstalled))
	result = append(result, pip.userInstalled...)
	return result
}

func (pip *PIP) AllInstalledPackages() []string {
	result := make([]string, 0, len(pip.allInstalled))
	result = append(result, pip.allInstalled...)
	return result
}

func RemoveFromList(list []string, elem string) []string {
	for i, v := range list {
		if v == elem {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}
func (pip *PIP) AllNeededDepsForCurPkgs() []string {
	result := make([]string, 0)
	result = append(result, pip.AllUserInstalledPackages()...)
	for _, pkg := range pip.AllUserInstalledPackages() {
		allDeps, err := pip.AllDeps(pkg)
		if err != nil {
			panic(err)
		}
		for _, dep := range allDeps {
			if !Contains(result, dep) {
				result = append(result, dep)
			}
		}
	}
	return result
}

func (pip *PIP) Fix() {
	err := pip.CheckAndRemoveDanglings()
	if err != nil {
		panic(err)
	}
	if pip.Check() == nil {
		return
	}
	for _, pkg := range pip.AllNeededDepsForCurPkgs() {
		if Contains(pip.AllInstalledPackages(), pkg) {
			continue
		}
		allDeps, err := pip.AllDeps(pkg)
		if err != nil {
			panic(err)
		}
		err = pip.CopyFromGopi(pkg)
		if err != nil {
			panic(err)
		}
		for _, v := range allDeps {
			err := pip.CopyFromGopi(v)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (pip *PIP) FindDanglings() []string {
	allNeededDeps := pip.AllNeededDepsForCurPkgs()
	danglingPackages := make([]string, 0)

	for _, installedPkg := range pip.AllInstalledPackages() {
		if !Contains(allNeededDeps, installedPkg) && !Contains(danglingPackages, installedPkg) {
			danglingPackages = append(danglingPackages, installedPkg)
		}
	}
	return danglingPackages
}

func (pip *PIP) CheckAndRemoveDanglings() error {
	return pip.UninstallForce(pip.FindDanglings()...)
}

func (pip *PIP) Uninstall(pkgNamesToRemove ...string) error {
	for _, pkgName := range pkgNamesToRemove {
		if !Contains(pip.AllUserInstalledPackages(), pkgName) {
			return fmt.Errorf(
				"pkg %s is not installed explicitly, so it cannot removed",
				pkgName,
			)
		}
	}
	// Check if any other user-installed packages depend on the packages to be removed
	for _, pkgNameToRemove := range pkgNamesToRemove {
		neededByOtherUserPkgs := make([]string, 0)
		for _, userPkg := range pip.AllUserInstalledPackages() {
			if userPkg == pkgNameToRemove { // Skip checking against itself
				continue
			}
			userPkgDeps, err := pip.AllDeps(userPkg)
			if err != nil {
				panic(err) // This should ideally be handled more gracefully
			}
			if Contains(userPkgDeps, pkgNameToRemove) {
				neededByOtherUserPkgs = append(neededByOtherUserPkgs, userPkg)
			}
		}
		if len(neededByOtherUserPkgs) > 0 {
			return fmt.Errorf("cannot remove %s because pkgs %v need this", pkgNameToRemove, neededByOtherUserPkgs)
		}
	}
	// If all checks pass, proceed with forced uninstallation
	if err := pip.UninstallForce(pkgNamesToRemove...); err != nil {
		return err
	}
	// After uninstallation, check for and remove any newly created dangling dependencies
	err := pip.CheckAndRemoveDanglings()
	return err
}

func (pip *PIP) UninstallForce(pkgNames ...string) error {
	for _, pkgName := range pkgNames {
		if !Contains(pip.AllInstalledPackages(), pkgName) {
			return fmt.Errorf("pkg %s is not installed", pkgName)
		}
	}
	for _, pkgName := range pkgNames {
		pip.userInstalled = RemoveFromList(pip.userInstalled, pkgName)
		pip.allInstalled = RemoveFromList(pip.allInstalled, pkgName)
		if err := pip.installDir.Remove(pkgName); err != nil {
			return err
		}
	}
	return nil
}

func (pip *PIP) Check() error {
	needed := pip.AllNeededDepsForCurPkgs()
	for _, need := range needed {
		if !Contains(pip.AllInstalledPackages(), need) {
			return fmt.Errorf("%s should be installed but its not", need)
		}
	}
	return nil
}

func StdLib(lib string) bool {
	stdLibList := []string{
		"archive",
		"bufio",
		"builtin",
		"bytes",
		"compress",
		"container",
		"context",
		"crypto",
		"database",
		"debug",
		"embed",
		"encoding",
		"errors",
		"expvar",
		"flag",
		"fmt",
		"go",
		"hash",
		"html",
		"image",
		"io",
		"log",
		"math",
		"mime",
		"net",
		"os",
		"path",
		"plugin",
		"reflect",
		"regexp",
		"runtime",
		"sort",
		"strconv",
		"strings",
		"sync",
		"syscall",
		"testing",
		"text",
		"time",
		"unicode",
		"unsafe",
	}
	for _, v := range stdLibList {
		if v == lib || (strings.HasPrefix(lib, v) && strings.Contains(lib, "/")) {
			return true
		}
	}
	return false
}

func (pip *PIP) ImportCheck(src string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ImportsOnly)
	if err != nil {
		return err
	}
	for _, imp := range f.Imports {
		impStr := strings.Trim(imp.Path.Value, "\"")
		if !Contains(pip.allInstalled, impStr) && !StdLib(impStr) {
			return errors.New("unsatisfied import" + impStr)
		}
	}
	return nil
}

func (pip *PIP) OnceInstalledPackages() []string {
	res := make([]string, 0, 20)
	for k := range pip.onceInstalled {
		res = append(res, k)
	}
	return res
}

func (pip *PIP) LocalSearch(term string) []string {
	return fuzzy.Find(term, pip.OnceInstalledPackages())
}
