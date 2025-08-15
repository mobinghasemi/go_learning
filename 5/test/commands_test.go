package main

import (
	"pip/commands"
	"pip/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	assert.NotNil(t, pip)
}

func TestDirectDeps1(t *testing.T) {
	deps, err := pip.DirectDeps("echo")
	assert.NoError(t, err)
	assert.ElementsMatch(t,
		[]string{"jwt", "testify", "fasttemplate"},
		deps,
	)
}
func TestDirectDeps2(t *testing.T) {
	_, err := pip.DirectDeps("non-existing-project")
	assert.Error(t, err)
}

func TestDirectDeps3(t *testing.T) {
	deps, err := pip.DirectDeps("prj-with-invalid-dep")
	assert.NoError(t, err)
	assert.ElementsMatch(t,
		[]string{"invalid-dep"},
		deps,
	)
}
func TestDirectDeps4(t *testing.T) {
	deps, err := pip.DirectDeps("go-difflib")
	assert.NoError(t, err)
	assert.ElementsMatch(t,
		[]string{},
		deps,
	)
}

func TestAllDeps1(t *testing.T) {
	deps, err := pip.AllDeps("go-difflib")
	assert.NoError(t, err)
	assert.ElementsMatch(t,
		[]string{},
		deps,
	)
}
func TestAllDeps2(t *testing.T) {
	_, err := pip.AllDeps("non-existing-project")
	assert.Error(t, err)
}

func TestAllDeps3(t *testing.T) {
	_, err := pip.AllDeps("prj-with-invalid-dep")
	assert.Error(t, err)
}

func TestAllDeps4(t *testing.T) {
	deps, err := pip.AllDeps("echo")
	assert.NoError(t, err)
	assert.ElementsMatch(t,
		[]string{
			"jwt",
			"testify",
			"fasttemplate",
			"go-spew",
			"go-difflib",
			"bytebufferpool",
		},
		deps,
	)
}

func TestInstall1(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	err := pip.Install("echo")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"echo"}, pip.AllUserInstalledPackages())
}

func TestInstall2(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	err := pip.Install("echo")
	assert.NoError(t, err)
	assert.ElementsMatch(t,
		[]string{
			"echo",
			"jwt",
			"testify",
			"fasttemplate",
			"go-spew",
			"go-difflib",
			"bytebufferpool",
		},
		pip.AllInstalledPackages(),
	)
}

func TestInstall3(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	err := pip.Install("echo")
	assert.NoError(t, err)
	err = pip.Install("echo")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"echo"}, pip.AllUserInstalledPackages())
}

func TestInstall4(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	err := pip.Install("echo")
	assert.NoError(t, err)
	err = pip.Install("jwt")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"echo", "jwt"}, pip.AllUserInstalledPackages())
	assert.ElementsMatch(t,
		[]string{
			"echo",
			"jwt",
			"testify",
			"fasttemplate",
			"go-spew",
			"go-difflib",
			"bytebufferpool",
		},
		pip.AllInstalledPackages(),
	)
}

func TestInstall5(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	err := pip.Install("echo", "jwt")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"echo", "jwt"}, pip.AllUserInstalledPackages())
}

func TestInstall6(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	err := pip.Install("numpy")
	assert.Error(t, err)
	assert.NotContains(t, pip.AllUserInstalledPackages(), "numpy")
	assert.NotContains(t, pip.AllInstalledPackages(), "numpy")
}

func TestInstall7(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	err := pip.Install("jwt", "numpy", "go-spew")
	assert.Error(t, err)

	assert.NotContains(t, pip.AllUserInstalledPackages(), "numpy")
	assert.NotContains(t, pip.AllInstalledPackages(), "numpy")

	// transaction is canceled
	assert.NotContains(t, pip.AllUserInstalledPackages(), "jwt")
	assert.NotContains(t, pip.AllInstalledPackages(), "jwt")

	assert.NotContains(t, pip.AllUserInstalledPackages(), "go-spew")
	assert.NotContains(t, pip.AllInstalledPackages(), "go-spew")
}
func TestInstall8(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	err := pip.Install("jwt", "prj-with-indirect-invalid-dep", "go-spew")
	assert.Error(t, err)

	assert.NotContains(t, pip.AllUserInstalledPackages(), "prj-with-indirect-invalid-dep")
	assert.NotContains(t, pip.AllInstalledPackages(), "prj-with-indirect-invalid-dep")

	// transaction is canceled
	assert.NotContains(t, pip.AllUserInstalledPackages(), "jwt")
	assert.NotContains(t, pip.AllInstalledPackages(), "jwt")

	assert.NotContains(t, pip.AllUserInstalledPackages(), "go-spew")
	assert.NotContains(t, pip.AllInstalledPackages(), "go-spew")
}

func TestInstallR1(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	d := fs.MkDir()
	err := d.CreateFile("requirements.txt")
	assert.NoError(t, err)
	err = d.AppendToFile("requirements.txt", "echo\n")
	assert.NoError(t, err)
	err = d.AppendToFile("requirements.txt", "jwt\n")
	assert.NoError(t, err)
	err = pip.InstallR(d, "requirements.txt")
	assert.NoError(t, err)

	assert.ElementsMatch(t, []string{"echo", "jwt"}, pip.AllUserInstalledPackages())
	assert.Contains(t, pip.AllInstalledPackages(), "echo")
	assert.Contains(t, pip.AllInstalledPackages(), "jwt")
}

func TestInstallR2(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	d := fs.MkDir()
	err := d.CreateFile("requirements.txt")
	assert.NoError(t, err)
	err = d.AppendToFile("requirements.txt", "echo\n")
	assert.NoError(t, err)
	err = d.AppendToFile("requirements.txt", "prj-with-indirect-invalid-dep\n")
	assert.NoError(t, err)
	err = pip.InstallR(d, "requirements.txt")
	assert.Error(t, err)

	assert.Empty(t, pip.AllUserInstalledPackages())
	assert.NotContains(t, pip.AllInstalledPackages(), "echo")
	assert.NotContains(t, pip.AllInstalledPackages(), "jwt")
	assert.NotContains(t, pip.AllInstalledPackages(), "prj-with-indirect-invalid-dep")
}

func TestUninstall1(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)

	err := pip.Install("jwt")
	assert.NoError(t, err)

	err = pip.Uninstall("jwt")
	assert.NoError(t, err)

	assert.NotContains(t, pip.AllUserInstalledPackages(), "jwt")
}

func TestUninstall2(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)

	err := pip.Install("echo")
	assert.NoError(t, err)

	err = pip.Uninstall("jwt")
	assert.Error(t, err)

	assert.Contains(t, pip.AllInstalledPackages(), "jwt")
}
func TestUninstall3(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)

	err := pip.Install("echo", "jwt")
	assert.NoError(t, err)

	err = pip.Uninstall("jwt")
	assert.Error(t, err)

	assert.Contains(t, pip.AllInstalledPackages(), "jwt")
}

func TestUninstall4(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)

	err := pip.Install("echo")
	assert.NoError(t, err)
	err = pip.Install("jwt")
	assert.NoError(t, err)

	err = pip.Uninstall("jwt")
	assert.Error(t, err)

	assert.Contains(t, pip.AllInstalledPackages(), "jwt")
}

func TestUninstall5(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)

	err := pip.Install("echo")
	assert.NoError(t, err)
	err = pip.Install("jwt")
	assert.NoError(t, err)

	err = pip.Uninstall("echo")
	assert.NoError(t, err)

	assert.Contains(t, pip.AllUserInstalledPackages(), "jwt")
	assert.Contains(t, pip.AllInstalledPackages(), "jwt")
	assert.NotContains(t, pip.AllInstalledPackages(), "echo")
	assert.NotContains(t, pip.AllUserInstalledPackages(), "echo")
}

func TestUninstall6(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)

	err := pip.Install("echo")
	assert.NoError(t, err)
	err = pip.Uninstall("echo")
	assert.NoError(t, err)
	assert.Empty(t, pip.AllUserInstalledPackages())
	assert.Empty(t, pip.AllInstalledPackages())
}

func TestUninstall7(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)

	err := pip.Install("echo")
	assert.NoError(t, err)
	err = pip.Uninstall("echo", "jwt")
	assert.Error(t, err)

	assert.Contains(t, pip.AllUserInstalledPackages(), "echo")
	assert.Contains(t, pip.AllInstalledPackages(), "echo")
	assert.Contains(t, pip.AllInstalledPackages(), "jwt")
}
func TestUninstall8(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)

	err := pip.Install("echo", "jwt")
	assert.NoError(t, err)
	err = pip.Uninstall("jwt")
	assert.Error(t, err)

	assert.Contains(t, pip.AllUserInstalledPackages(), "echo")
	assert.Contains(t, pip.AllUserInstalledPackages(), "jwt")
	assert.Contains(t, pip.AllInstalledPackages(), "echo")
	assert.Contains(t, pip.AllInstalledPackages(), "jwt")
}

func TestForceUninstall1(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)

	err := pip.Install("echo")
	assert.NoError(t, err)

	err = pip.UninstallForce("jwt")
	assert.NoError(t, err)
}

func TestForceUninstall2(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)

	err := pip.Install("echo")
	assert.NoError(t, err)

	err = pip.UninstallForce("numpy")
	assert.Error(t, err)
}

func TestCheck1(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)

	err := pip.Install("echo")
	assert.NoError(t, err)

	err = pip.UninstallForce("jwt")
	assert.NoError(t, err)

	err = pip.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "jwt")
}

func TestCheck2(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)

	err := pip.Install("echo")
	assert.NoError(t, err)

	err = pip.UninstallForce("jwt")
	assert.NoError(t, err)

	err = pip.UninstallForce("testify")
	assert.NoError(t, err)

	err = pip.Check()
	assert.Error(t, err)

	err = pip.Install("jwt")
	assert.NoError(t, err)

	err = pip.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "testify")
}

func TestFix1(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)

	err := pip.Install("echo")
	assert.NoError(t, err)

	err = pip.UninstallForce("echo")
	assert.NoError(t, err)

	pip.Fix()

	assert.Empty(t, pip.AllInstalledPackages())
}

func TestFix2(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)

	err := pip.Install("echo")
	assert.NoError(t, err)

	err = pip.UninstallForce("jwt")
	assert.NoError(t, err)

	err = pip.UninstallForce("testify")
	assert.NoError(t, err)

	err = pip.Check()
	assert.Error(t, err)

	pip.Fix()

	err = pip.Check()
	assert.NoError(t, err)

	assert.Contains(t, pip.AllInstalledPackages(), "jwt")
	assert.Contains(t, pip.AllInstalledPackages(), "testify")
	assert.Contains(t, pip.AllInstalledPackages(), "go-difflib")

	assert.NotContains(t, pip.AllUserInstalledPackages(), "jwt")
	assert.NotContains(t, pip.AllUserInstalledPackages(), "testify")
	assert.NotContains(t, pip.AllUserInstalledPackages(), "go-difflib")
}

func TestImportCheck1(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	src := `package main
import (
	"echo"
	"jwt"
)
func main(){

}
`

	err := pip.ImportCheck(src)
	assert.Error(t, err)

	err = pip.Install("jwt")
	assert.NoError(t, err)

	err = pip.ImportCheck(src)
	assert.Error(t, err)

	err = pip.Install("echo")
	assert.NoError(t, err)

	err = pip.ImportCheck(src)
	assert.NoError(t, err)
}

func TestImportCheck2(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	src := `package main

import ( 	"fmt"
	"os" )
import "net/http"

func main() {
}
`
	err := pip.ImportCheck(src)
	assert.NoError(t, err)
}

func TestImportCheck3(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	src := `package main

import ( 	"fmt"
	"os"
)
import "net/http"
import "testify"
import "testing"
import "errors"
import "sync"
import "syscall"
import (
	"unicode"
)
import "mime"

import 			"crypto"
import ( 	"html"
	"image"
"log"
  "sort"
)
import "unsafe"
import     "unicode"
import ( "hash" )
   import        "bytes"

func main() {
}
`
	err := pip.ImportCheck(src)
	assert.Error(t, err)

	err = pip.Install("testify")
	assert.NoError(t, err)

	err = pip.ImportCheck(src)
	assert.NoError(t, err)
}

func TestImportCheck4(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	src := `package main

import 			"crypto"
import ( 	"html"
	"image"
"logging"
  "sort"
)
import "unsafe"
import     "unicode"
import ( "hash" )
   import        "bytes"

func main() {
}
`
	err := pip.ImportCheck(src)
	assert.Error(t, err)
}

func TestSearch1(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	err := pip.Install("echo")
	assert.NoError(t, err)
	result := pip.LocalSearch("ech")
	assert.Equal(t, []string{"echo"}, result)
}

func TestSearch2(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	err := pip.Install("echo")
	assert.NoError(t, err)
	result := pip.LocalSearch("cho")
	assert.Equal(t, []string{"echo"}, result)
}

func TestSearch3(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	err := pip.Install("echo")
	assert.NoError(t, err)
	result := pip.LocalSearch("e")
	assert.ElementsMatch(t,
		[]string{"echo", "testify", "fasttemplate", "go-spew", "bytebufferpool"},
		result,
	)
}

func TestSearch4(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	err := pip.Install("echo")
	assert.NoError(t, err)
	result := pip.LocalSearch("te")
	assert.ElementsMatch(t,
		[]string{"testify", "fasttemplate", "bytebufferpool"},
		result,
	)
}

func TestSearch5(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	err := pip.Install("echo")
	assert.NoError(t, err)
	result := pip.LocalSearch("tet")
	assert.ElementsMatch(t,
		[]string{"testify", "fasttemplate"},
		result,
	)
}

func TestSearch6(t *testing.T) {
	pip := commands.NewPIP(fs.MkDir(), gopi)
	err := pip.Install("echo")
	assert.NoError(t, err)

	// answer from cache
	err = pip.Uninstall("echo")
	assert.NoError(t, err)
	result := pip.LocalSearch("tet")
	assert.ElementsMatch(t,
		[]string{"testify", "fasttemplate"},
		result,
	)
}
