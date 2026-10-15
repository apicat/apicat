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
	mux, err := newHandler(getDoc)
	if err != nil {
		return err
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Serving API docs at http://%s\n", addr)
	return http.ListenAndServe(addr, mux)
}

func newHandler(getDoc func() *model.APIDoc) (http.Handler, error) {
	mux := http.NewServeMux()

	staticFS, err := fs.Sub(web.Content, "static")
	if err != nil {
		return nil, fmt.Errorf("cannot create static filesystem: %w", err)
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

	mux.HandleFunc("/components", func(w http.ResponseWriter, r *http.Request) {
		data := &pageData{APIDoc: getDoc(), Page: "components"}
		data.ComponentType = r.URL.Query().Get("type")
		data.ComponentName = r.URL.Query().Get("name")
		if data.ComponentType != "" || data.ComponentName != "" {
			for _, group := range data.ComponentGroups {
				if group.Kind != data.ComponentType {
					continue
				}
				for i := range group.Definitions {
					if group.Definitions[i].Name == data.ComponentName {
						data.Component = &group.Definitions[i]
						data.ComponentTitle = group.Title
						break
					}
				}
			}
			if data.Component == nil {
				http.NotFound(w, r)
				return
			}
		}
		renderPage(w, data)
	})
	return mux, nil
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
	Page           string
	ComponentType  string
	ComponentName  string
	ComponentTitle string
	Component      *model.Component
}
