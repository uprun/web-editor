//go:build ignore

package main

import (
    "fmt"
    "os"
    "log"
    "net/http"
)

type Page struct {
    Title string
    Body  []byte
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("view %s", r.URL.Path)
    title := r.URL.Path[len("/view/"):]
    p, _ := loadPage(title)
    fmt.Fprintf(w, "%s", p.Body)
}


func main() {
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/save/", saveHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadPage(title string) (*Page, error) {
    filename := title
    body, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}


func (p *Page) save() error {
    filename := p.Title
    return os.WriteFile(filename, p.Body, 0600)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
    title := r.FormValue("title")
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    p.save()
}
