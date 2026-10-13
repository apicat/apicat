package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/apicat/apicat/v3/internal/converter"
	"github.com/apicat/apicat/v3/internal/lint"
	"github.com/apicat/apicat/v3/internal/loader"
	"github.com/apicat/apicat/v3/internal/model"
	"github.com/apicat/apicat/v3/internal/project"
	"github.com/apicat/apicat/v3/internal/render"
	"github.com/apicat/apicat/v3/internal/server"
	"github.com/apicat/apicat/v3/internal/watcher"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	host := flag.String("host", "127.0.0.1", "host to bind to")
	watch := flag.Bool("watch", false, "reload when spec files in the directory change")
	flag.Parse()

	var dir string
	if flag.NArg() >= 1 {
		dir = flag.Arg(0)
	} else {
		// Project mode: resolve the apicat working directory for the
		// current project from ~/.apicat/projects/<encoded-cwd>.
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: cannot determine working directory: %v\n", err)
			os.Exit(1)
		}
		dir, err = project.DirFor(cwd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if _, statErr := os.Stat(dir); statErr != nil {
			fmt.Fprintf(os.Stderr, "Error: no apicat project for %s\n", cwd)
			fmt.Fprintf(os.Stderr, "Expected a spec at %s/openapi.yaml (or .yml/.json).\n", dir)
			fmt.Fprintf(os.Stderr, "Create it there, or pass a spec directory explicitly: apicat-cli [flags] <spec-directory>\n")
			os.Exit(1)
		}
		fmt.Printf("Project %s\nSpec dir %s\n", cwd, dir)
	}

	loadDoc := func() (*model.APIDoc, error) {
		docModel, err := loader.Load(dir)
		if err != nil {
			return nil, err
		}
		doc, err := converter.Convert(docModel)
		if err != nil {
			return nil, err
		}
		lint.Annotate(doc)
		return doc, nil
	}

	apiDoc, err := loadDoc()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var current atomic.Pointer[model.APIDoc]
	current.Store(apiDoc)

	if *watch {
		go func() {
			err := watcher.Watch(dir, 300*time.Millisecond, func() {
				doc, err := loadDoc()
				if err != nil {
					log.Printf("spec reload failed, keeping previous version: %v", err)
					return
				}
				current.Store(doc)
				log.Printf("spec reloaded")
			})
			if err != nil {
				log.Printf("watch disabled: %v", err)
			}
		}()
		fmt.Printf("Watching %s for changes\n", dir)
	}

	if err := render.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading templates: %v\n", err)
		os.Exit(1)
	}

	if err := server.Start(current.Load, *host, *port); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
