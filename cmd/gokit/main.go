package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	gokithttp "github.com/codedbyhassan/gokit/api/http"
	"github.com/codedbyhassan/gokit/interpret/pipeline"
)

func main() {
	serve := flag.Bool("serve", false, "start the GoKit HTTP API")
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()

	if *serve {
		log.Printf("GoKit API listening on %s", *addr)
		log.Fatal(http.ListenAndServe(*addr, gokithttp.NewServer().Handler()))
	}

	input := flag.Args()
	if len(input) == 0 { fmt.Println("usage: gokit [--serve] [--addr :8080] <input>"); return }
	result, err := pipeline.Parse(fmt.Sprint(input...))
	if err != nil { log.Fatal(err) }
	fmt.Println(result.Value)
}
