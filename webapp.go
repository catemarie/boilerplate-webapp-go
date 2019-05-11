package main

import (
	"fmt"
	"io/ioutil"
)

type Entry struct {
    Title string
    Body  []byte
}

func (p *Entry) save() error {
    filename := p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Entry, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Entry{Title: title, Body: body}, nil
}

func main() {
    p1 := &Entry{Title: "Sample Entry", Body: []byte("Sample Text")}
    p1.save()
    p2, _ := loadPage("Sample Entry")
    fmt.Println(string(p2.Body))
}