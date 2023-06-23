package main

import (
  "embed"
  "io/fs"
  "log"
  "net/http"
)

//go:embed public
var public embed.FS

func main() {
  // We want to serve static content from the root of the 'public' directory,
  // but go:embed will create a FS where all the paths start with 'public/...'.
  // Using fs.Sub we "cd" into 'public' and can serve files relative to it.
  publicFS, err := fs.Sub(public, "public")
  if err != nil {
    log.Fatal(err)
  }

  http.Handle("/", http.FileServer(http.FS(publicFS)))

  port := ":9223"
  log.Fatal(http.ListenAndServe(port, nil))
}
