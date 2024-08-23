
  handleRequest := func(w http.ResponseWriter, r *http.Request) {

    if debug == "true" {
      fmt.Fprintf(os.Stderr, "Request received: %s", r.URL.Path)
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
