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

  // Attempt to get the port number from the `GOLANG_PORT` environment variable
  port, portIsSet := os.LookupEnv("GOLANG_PORT")

  // If the `GOLANG_PORT` environment variable is not set, use the default port
  if !portIsSet {
    port = "9223"
  }

  // Attempt to get the build directory from the `GOLANG_STATIC_CONTENT_DIRECTORY` environment variable
  staticContentDirectory, staticContentDirectoryIsSet := os.LookupEnv("GOLANG_STATIC_CONTENT_DIRECTORY")

  // If the `GOLANG_STATIC_CONTENT_DIRECTORY` environment variable is not set, use the default static content directory
  if !staticContentDirectoryIsSet {
    staticContentDirectory = "/static"
  }

  staticContentDirectoryExists, err := exists(staticContentDirectory)

  if staticContentDirectoryExists {

      publicFS := http.FileServer(http.Dir("static"))
      http.Handle("/static/", http.StripPrefix("/static/", fs))
    
  } else {

      // Using `fs.Sub()`, create a filesystem which uses the 'public' directory as its root
      publicFS, err := fs.Sub(public, "public")
  }

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
