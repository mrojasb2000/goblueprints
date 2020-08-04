package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/mrojasb2000/goblueprints/chapter2/trace"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Host: %s\n", r.Host)
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", "9000", "The addr of the application.")
	flag.Parse() // parse the arguments and extracts the appropriate flags
	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	// root
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)

	// get the room going
	// separate goroutine
	go r.run()

	address := fmt.Sprintf("%s:%s", "0.0.0.0", *addr)
	//fmt.Printf("Starting web server on: %d\n", *addr)
	log.Println("Starting web server on", *addr)

	// start the web server
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal("ListeringAndServe:", err)
	}
}
