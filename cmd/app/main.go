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

func main() {

  // Attempt to get the port number from the `GOLANG_PORT` environment variable
  port, portIsSet := os.LookupEnv("GOLANG_PORT")

  // If the `GOLANG_PORT` environment variable is not set, use the default port
  if !portIsSet {
    // Declare the default port
    port = ":9223"
  }

  // We want to serve static content from the root of the 'public' directory,
  // but go:embed will create a FS where all the paths start with 'public/...'.
  // Using fs.Sub we "cd" into 'public' and can serve files relative to it.
  publicFS, err := fs.Sub(public, "public")

  // Throw an error if the filesystem cannot be created
  if err != nil {
    log.Fatal(err)
  }

  // Serve the filesystem under the root path
  http.Handle("/", http.FileServer(http.FS(publicFS)))

  // If the service fails, log the service
  log.Fatal(http.ListenAndServe(port, nil))
}
