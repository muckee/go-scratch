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

func main() {

  debug, _ := os.LookupEnv("GOLANG_DEBUG")

  basepath, basepathIsSet := os.LookupEnv("GOLANG_BASEPATH")
  if !basepathIsSet {
    basepath = "/"
  }

  // Determine which port the application should be served on
  port, portIsSet := os.LookupEnv("GOLANG_PORT")
  if !portIsSet {
    port = "9123"
  }

  // Get the path of the static content directory from the `GOLANG_STATIC_CONTENT_DIRECTORY` environment variable or use `/static`
  staticContentDirectory, staticContentDirectoryIsSet := os.LookupEnv("GOLANG_STATIC_CONTENT_DIRECTORY")
  if !staticContentDirectoryIsSet {
    staticContentDirectory = "/static"
  }

  publicFS, err := fs.Sub(public, "public")
  httpFS := http.FileServer(http.FS(publicFS))

  // If the static content directory exists, assign it as the public directory
  staticContentDirectoryExists, err := exists(staticContentDirectory)
  if staticContentDirectoryExists {
      httpFS = http.FileServer(http.Dir(fmt.Sprintf("%s", staticContentDirectory)))
  }

  handleRequest := func(w http.ResponseWriter, r *http.Request) {

    fmt.Fprintf(os.Stderr, "Request received: %s", r.URL.Path)

    if debug == "true" {
      fmt.Fprintf(os.Stderr, "Request received: %s", r.URL.Path)
    }

    // if basepath == "/" {
    	httpFS.ServeHTTP(w, r)
    // } else {
    //   http.StripPrefix(fmt.Sprintf("%s/", basepath), httpFS).ServeHTTP(w, r)
    // }
  }

  http.HandleFunc(basepath, handleRequest)
  // http.Handle("/", httpFS)

  // Check for handling errors
  if err != nil {
    log.Fatal(err)
  }

  // Serve the static content
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
