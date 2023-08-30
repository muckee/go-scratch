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

var publicFS fs.FS

var httpFS http.Handler

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func main() {

  // Step 1) Determine which port should be used to serve static content

  // Attempt to retrieve the port from env vars
  port, portIsSet := os.LookupEnv("GOLANG_PORT")

  // If no port has been set via env vars, use `9223` as the fallback port
  if !portIsSet {
    port = "9223"
  }

  // Step 2) Determine the location of the public filesystem
  publicFS, err := fs.Sub(public, "public")
  httpFS := http.FileServer(http.FS(publicFS))

  // Attempt to get the build directory from the `GOLANG_STATIC_CONTENT_DIRECTORY` environment variable
  staticContentDirectory, staticContentDirectoryIsSet := os.LookupEnv("GOLANG_STATIC_CONTENT_DIRECTORY")

  // If the `GOLANG_STATIC_CONTENT_DIRECTORY` environment variable is not set, use the default static content directory
  if !staticContentDirectoryIsSet {
    staticContentDirectory = "static"
  }

  // Check if the static content directory exists
  staticContentDirectoryExists, err := exists(staticContentDirectory)

  // If the static content directory exists, assign it as the value of `publicFS`
  if staticContentDirectoryExists {
      publicFS := http.FileServer(http.Dir(fmt.Sprintf(":%s", staticContentDirectory)))
      httpFS := http.StripPrefix(fmt.Sprintf("/:%s/", staticContentDirectory), publicFS)
  }

  // Point the root endpoint at the chosen filesystem
  http.Handle("/", httpFS)

  // Throw an error if the filesystem cannot be created
  if err != nil {
    log.Fatal(err)
  }

  // Serve the static content
  // The return value of the `http.ListenAndServe()` command is always logged as a fatal error
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
