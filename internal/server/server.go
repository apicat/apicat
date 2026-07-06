package server

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/apicat/apicat/v3/internal/model"
	"github.com/apicat/apicat/v3/internal/render"
	"github.com/apicat/apicat/v3/web"
)

// Start serves the docs until the process exits. getDoc is called per
// request so -watch can swap in a freshly loaded spec.
func Start(getDoc func() *model.APIDoc, host string, port int) error {
	mux := http.NewServeMux()

	staticFS, err := fs.Sub(web.Content, "static")
	if err != nil {
		return fmt.Errorf("cannot create static filesystem: %w", err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		renderPage(w, &pageData{APIDoc: getDoc(), Page: "overview"})
	})

	mux.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, &pageData{APIDoc: getDoc(), Page: "endpoints"})
	})

	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Serving API docs at http://%s\n", addr)
	return http.ListenAndServe(addr, mux)
}

func renderPage(w http.ResponseWriter, data *pageData) {
	var buf bytes.Buffer
	if err := render.Execute(&buf, "layout", data); err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

type pageData struct {
	*model.APIDoc
	Page string
}
