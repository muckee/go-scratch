package main

import (
  "bytes"
  "embed"
  "fmt"
  "io/fs"
  "log"
  "net/http"
  "os"
  "os/exec"
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

func isGolangApplication(path string) bool {
    file, err := os.Open(path)
    if err != nil {
        return false
    }
    defer file.Close()

    // Define the Golang binary signature (magic number)
    golangMagicNumber := []byte{0x7f, 'E', 'L', 'F'}

    // Read the first 4 bytes of the file
    buf := make([]byte, 4)
    _, err = file.Read(buf)
    if err != nil {
        return false
    }

    // Compare the read bytes with the Golang binary signature
    return bytes.Equal(buf, golangMagicNumber)
}

func getStaticContentDirectory() string {

  // Get the path of the static content directory from the `GOLANG_STATIC_CONTENT_DIRECTORY` environment variable or use `/static`
  staticContentDirectory, staticContentDirectoryIsSet := os.LookupEnv("GOLANG_STATIC_CONTENT_DIRECTORY")
  if !staticContentDirectoryIsSet {
    staticContentDirectory = "/static"
  }

  return staticContentDirectory
}

func directoryIsValid(path string) bool {

  // If the static content directory exists, assign it as the public directory
  directoryExists, err := exists(path)

  if err != nil {
    return false
  }

  if !directoryExists {
    return false
  }

  return true
}

func main() {

  // Determine whether application is in debug mode
  debug, _ := os.LookupEnv("GOLANG_DEBUG")

  // Determine URL basepath to be used as HTTP route
  basepath, basepathIsSet := os.LookupEnv("GOLANG_URL_PREFIX")
  if !basepathIsSet {
    basepath = "/"
  }

  // Determine which port the application should be served on
  port, portIsSet := os.LookupEnv("GOLANG_PORT")
  if !portIsSet {
    port = "9123"
  }

  // Get the path of the static content directory from the `GOLANG_STATIC_CONTENT_DIRECTORY` environment variable or use `/static`
  staticContentDirectory := getStaticContentDirectory()

  publicFS, err := fs.Sub(public, "public")

  httpFS := http.FileServer(http.FS(publicFS))

  if directoryIsValid(staticContentDirectory) {
      httpFS = http.FileServer(http.Dir(fmt.Sprintf("%s", staticContentDirectory)))
  }

  handleRequest := func(w http.ResponseWriter, r *http.Request) {
    if debug == "true" {
        fmt.Fprintf(os.Stderr, "Request received: %s\n", r.URL.Path)
    }

    // Check if the request is for the Go application
    if r.URL.Path == "/goapp/app" {
        if isGolangApplication("/goapp/app") {
            // Execute the Golang application as a separate process
            cmd := exec.Command("/goapp/app")
            cmd.Stdout = w
            cmd.Stderr = w
            err := cmd.Run()
            if err != nil {
                // Handle the error if needed
                http.Error(w, "Error running the Golang application", http.StatusInternalServerError)
            }
        } else {
            // Return an error response if it's not a valid executable
            http.Error(w, "Failed to execute: not a valid executable", http.StatusNotFound)
        }
        return
    }

    // Attempt to serve static files first
    requestedPath := fmt.Sprintf("%s%s", staticContentDirectory, r.URL.Path)
    fileExists, err := exists(requestedPath)

    if err == nil && fileExists {
        // If the requested file exists, serve it
        if basepath == "/" {
            httpFS.ServeHTTP(w, r)
        } else {
            http.StripPrefix(fmt.Sprintf("%s", basepath), httpFS).ServeHTTP(w, r)
        }
        return
    }

    // If the file doesn't exist and the request is not for a file or API, serve the index.html
    if r.URL.Path == "/" || !fileExists {
        indexPath := fmt.Sprintf("%s/index.html", staticContentDirectory)
        http.ServeFile(w, r, indexPath)
        return
    }

    // If none of the above conditions are met, serve the request normally
    if basepath == "/" {
        httpFS.ServeHTTP(w, r)
    } else {
        http.StripPrefix(fmt.Sprintf("%s", basepath), httpFS).ServeHTTP(w, r)
    }
  }


  http.HandleFunc("/", handleRequest)

  // Check for handling errors
  if err != nil {
    log.Fatal(err)
  }

  // Serve the static content
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
