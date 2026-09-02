package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	gokithttp "github.com/codedbyhassan/gokit/api/http"
	"github.com/codedbyhassan/gokit/interpret/pipeline"
	"github.com/codedbyhassan/gokit/web/frontend"
)

func main() {
	serve := flag.Bool("serve", false, "start the GoKit HTTP API")
	web := flag.Bool("web", false, "start the GoKit web frontend")
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()

	if *web {
		server, err := frontend.NewServer()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("GoKit web frontend listening on %s", *addr)
		log.Fatal(http.ListenAndServe(*addr, server.Handler()))
	}

	if *serve {
		log.Printf("GoKit API listening on %s", *addr)
		log.Fatal(http.ListenAndServe(*addr, gokithttp.NewServer().Handler()))
	}

	input := flag.Args()
	if len(input) == 0 {
		fmt.Println("usage: gokit [--serve|--web] [--addr :8080] <input>")
		return
	}
	result, err := pipeline.Parse(strings.Join(input, " "))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Value)
}
