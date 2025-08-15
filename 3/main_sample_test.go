package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TLDRSample struct {
	data map[string]string
}

func (t *TLDRSample) Retrieve(key string) string {
	return t.data[key]
}

func (t *TLDRSample) List() []string {
	res := make([]string, 0, len(t.data))
	for k := range t.data {
		res = append(res, k)
	}
	return res
}

func NewTLDRSample() TLDRProvider {
	return &TLDRSample{
		data: map[string]string{
			"ls":   "ls is good",
			"bash": "bash os also good",
			"lsd":  "even better than ls",
			"zsh":  "a posix complaint shell with cool features",
			"go":   "the famous go compiler from Google inc",
		},
	}
}

func BeforeEachSample() {
	GetConnection().Migrator().DropTable(&TLDREntity{})
	GetConnection().AutoMigrate(&TLDREntity{})
}

func TestRetrieveSample(t *testing.T) {
	BeforeEachSample()
	imd := NewTLDRSample()
	cached := NewTLDRDBCached(imd)
	assert.Equal(t, imd.Retrieve("bash"), cached.Retrieve("bash")) // slow
	assert.Equal(t, imd.Retrieve("bash"), cached.Retrieve("bash")) // fast(cached)
}

func TestListSample(t *testing.T) {
	BeforeEachSample()
	imd := NewTLDRSample()
	cached := NewTLDRDBCached(imd)
	assert.ElementsMatch(t, imd.List(), cached.List())
}
