//go:build ignore

package main

import (
    "fmt"
    "os"
    "log"
    "net/http"
    "os/exec"
)

type Page struct {
    Title string
    Body  []byte
}

func handler(w http.ResponseWriter, r *http.Request) {
    body, _ := os.ReadFile("../web-editor-2025-05-16-23h13m.html")
    fmt.Fprintf(w,"%s", body)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    param_path := query.Get("path")
    fmt.Printf("view %s", param_path)
    fmt.Println("")
    p, _ := loadPage(param_path)
    fmt.Fprintf(w, "%s", p.Body)
}

func commandHandler(w http.ResponseWriter, r *http.Request) {
    cmd := exec.Command("ls", "-la")
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Fprintf(w, "[web-shell] Error happened:\n%s\n[web-shell] output:\n", err, output)
        return
    }
    fmt.Fprintf(w, "%s", output)
}


func main() {
    fmt.Println(len(os.Args), os.Args)
    port := "8080"
    if len(os.Args) > 1 {
        port = os.Args[1]
    }
    fmt.Printf("web-editor started at http://localhost:%s/\n", port)
    http.HandleFunc("/command/", commandHandler)
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/save/", saveHandler)
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":" + port, nil))
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