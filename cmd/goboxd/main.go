package main

import (
	"log"
	"net/http"
	"os"

	"github.com/thesouldev/goboxd/internal/runner"
	"github.com/thesouldev/goboxd/internal/server"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("GOBXD_ADDR"); v != "" {
		addr = v
	}

	srv := server.New(runner.New())
	h := srv.Handler()

	log.Printf("goboxd listening on %s", addr)
	if err := http.ListenAndServe(addr, h); err != nil {
		log.Fatal(err)
	}
}
