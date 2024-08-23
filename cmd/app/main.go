package main

import (
  "bytes"
  "embed"
  "fmt"
  "io"
  "io/fs"
  "log"
  "net/http"
  "os"
  "os/exec"
  "strings"
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


  // Define the service name and port for the Next.js API
  apiServiceName := "actions-thugnerdz"  // Replace with your service name
  apiServicePort := "3000"               // Port on which the API is exposed

  handleRequest := func(w http.ResponseWriter, r *http.Request) {

    if debug == "true" {
      fmt.Fprintf(os.Stderr, "Request received: %s", r.URL.Path)
    }

    // Handle API requests
    if r.URL.Path == "/api/" || r.URL.Path[:5] == "/api/" {
        apiProxy.ServeHTTP(w, r)
        return
    }

    // Check if the request is for the API
    if strings.HasPrefix(r.URL.Path, "/api/") {
        // Forward API requests to the Next.js API server
        apiURL := fmt.Sprintf("http://%s:%s", apiServiceName, apiServicePort)
        apiURL = apiURL + r.URL.Path
        resp, err := http.Get(apiURL)
        if err != nil {
            http.Error(w, "Error forwarding request to API", http.StatusInternalServerError)
            return
        }
        defer resp.Body.Close()
        w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
        w.WriteHeader(resp.StatusCode)
        io.Copy(w, resp.Body)
        return
    }

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
            // Return an error response or handle it as needed
            http.Error(w, "Failed to execute: not a valid executable", http.StatusNotFound)
        }
    } else if basepath == "/" {
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
