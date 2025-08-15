package main

import "gorm.io/gorm"

type TLDRProvider interface {
	Retrieve(string) string
	List() []string
}

type TLDREntity struct {
	gorm.Model
	Key string `gorm:"primaryKey;size:100"`
	Val string `gorm:"size:1000"`
}
