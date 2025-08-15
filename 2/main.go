package library

import (
	"strings"
)

const (
	msgNotBorrowed        = "The book has not been borrowed"
	msgAlreadyBorrowed    = "The book is already borrowed by"
	msgAlreadyExist       = "The book is already in the library"
	msgNotDefindInLibrary = "The book is not defined in the library"
	msgNotEnoughCapacity  = "Not enough capacity"
	msgOK                 = "OK"
)

type Book struct {
	Title    string
	Borrower string
}
type Library struct {
	capacity int
	books    map[string]*Book
}

func NewLibrary(capacity int) *Library {
	if capacity < 0 {
		capacity = 0
	}
	return &Library{
		capacity: capacity,
		books:    make(map[string]*Book, capacity),
	}
}

func (library *Library) AddBook(name string) string {
	bookNameLower := strings.ToLower(strings.TrimSpace(name))
	if _, found := library.books[bookNameLower]; found {
		return msgAlreadyExist
	}
	if len(library.books) >= library.capacity {
		return msgNotEnoughCapacity
	}
	library.books[bookNameLower] = &Book{Title: name, Borrower: ""}
	return msgOK
}

func (library *Library) BorrowBook(bookName, personName string) string {
	bookNameLower := strings.ToLower(strings.TrimSpace(bookName))
	b, found := library.books[bookNameLower]
	if !found {
		return msgNotDefindInLibrary
	}
	if b.Borrower != "" {
		return msgAlreadyBorrowed + " " + b.Borrower
	}
	b.Borrower = strings.TrimSpace(personName)
	return msgOK
}
func (library *Library) ReturnBook(bookName string) string {
	bookNameLower := strings.ToLower(strings.TrimSpace(bookName))
	b, found := library.books[bookNameLower]
	if !found {
		return msgNotDefindInLibrary
	}
	if b.Borrower == "" {
		return msgNotBorrowed
	}
	b.Borrower = ""
	return msgOK
}
