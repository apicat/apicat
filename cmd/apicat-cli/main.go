package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/apicat/apicat/v3/internal/converter"
	"github.com/apicat/apicat/v3/internal/loader"
	"github.com/apicat/apicat/v3/internal/render"
	"github.com/apicat/apicat/v3/internal/server"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	host := flag.String("host", "127.0.0.1", "host to bind to")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: apicat-cli [flags] <spec-directory>\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	dir := flag.Arg(0)

	docModel, err := loader.Load(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	apiDoc, err := converter.Convert(docModel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := render.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading templates: %v\n", err)
		os.Exit(1)
	}

	if err := server.Start(apiDoc, *host, *port); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
