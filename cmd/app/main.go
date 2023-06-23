package main

import (
  "embed"
  "fmt"
  "io/fs"
  "log"
  "net/http"
  "os"
)

//go:embed public
var public embed.FS

func main() {

  // Attempt to get the port number from the `GOLANG_PORT` environment variable
  port, portIsSet := os.LookupEnv("GOLANG_PORT")

  // If the `GOLANG_PORT` environment variable is not set, use the default port
  if !portIsSet {
    port = "9223"
  }

  // Using `fs.Sub()`, create a filesystem which uses the 'public' directory as its root
  publicFS, err := fs.Sub(public, "public")

  // Throw an error if the filesystem cannot be created
  if err != nil {
    log.Fatal(err)
  }

  // Point the root endpoint at the filesystem created from the 'public' directory
  http.Handle("/", http.FileServer(http.FS(publicFS)))

  // Serve the static content
  // The return value of the `http.ListenAndServe()` command is always logged as a fatal error
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
