package main

import (
    "errors"
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "regexp"
)

var templates = template.Must(template.ParseGlob("static/templates/*.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

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

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Entry")
    }
    return m[2], nil // The title is the second subexpression.
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Entry) {
    err := templates.ExecuteTemplate(w, tmpl + ".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadEntry(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadEntry(title)
    if err != nil {
        p = &Entry{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    p := &Entry{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    log.Fatal(http.ListenAndServe(":8080", nil))
}