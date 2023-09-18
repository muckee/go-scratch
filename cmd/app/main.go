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

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func logRequests(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(os.Stderr, "Request received: %s", r.URL.Path)
		f(w, r)
	}
}

func main() {

  // Step 1) Determine which port should be used to serve static content

  // Retrieve the port from the `GOLANG_PORT` environment variable or use `:9123`
  port, portIsSet := os.LookupEnv("GOLANG_PORT")

  // If no port has been set via env vars, use `9223` as the fallback port
  if !portIsSet {
    port = "9123"
  }

  // Step 2) Determine the location of the public filesystem
  publicFS, err := fs.Sub(public, "public")
  httpFS := http.FileServer(http.FS(publicFS))

  // Get the path of the static content directory from the `GOLANG_STATIC_CONTENT_DIRECTORY` environment variable or use `/static`
  staticContentDirectory, staticContentDirectoryIsSet := os.LookupEnv("GOLANG_STATIC_CONTENT_DIRECTORY")
  if !staticContentDirectoryIsSet {
    staticContentDirectory = "/static"
  }

  // If the static content directory exists, assign it as the public directory
  staticContentDirectoryExists, err := exists(staticContentDirectory)
  if staticContentDirectoryExists {
      httpFS = http.FileServer(http.Dir(fmt.Sprintf("%s", staticContentDirectory)))
  }

  // Retrieve the port from the `GOLANG_PORT` environment variable or use `:9123`
  debug, _ := os.LookupEnv("GOLANG_DEBUG")

  if debug == "true" {
    // Handle all requests
    http.Handle("/", logRequests(httpFS))
  } else {
    // Handle all requests
    http.Handle("/", logRequests(httpFS))
  }

  // Check for handling errors
  if err != nil {
    log.Fatal(err)
  }

  // Serve the static content
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
