package main

import (
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
)

type Entry struct {
    Title string
    Body  []byte
}

func (p *Entry) save() error {
    filename := p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadEntry(title string) (*Entry, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Entry{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Entry) {
    t, _ := template.ParseFiles(tmpl + ".html")
    t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, _ := loadEntry(title)
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadEntry(title)
    if err != nil {
        p = &Entry{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
}

func main() {
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/save/", saveHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}