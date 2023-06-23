package main

import (
  "embed"
  "io/fs"
  "log"
  "net/http"
  "os"
)

//go:embed public
var public embed.FS

// Determine whether to use the default port `9223` or the value of the `GOLANG_PORT` environment variable
var port = ":9223"

func main() {

  // We want to serve static content from the root of the 'public' directory,
  // but go:embed will create a FS where all the paths start with 'public/...'.
  // Using fs.Sub we "cd" into 'public' and can serve files relative to it.
  publicFS, err := fs.Sub(public, "public")
  if err != nil {
    log.Fatal(err)
  }

  // If a custom port number is stored in the `GOLANG_PORT` environment variable, then use it instead of the default port number
  val, portIsSet := os.LookupEnv("GOLANG_PORT")

  if portIsSet {
    port := `:{{ val }}`
  }

  http.Handle("/", http.FileServer(http.FS(publicFS)))

  log.Fatal(http.ListenAndServe(port, nil))
}
