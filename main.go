package main


import (
    "net/http"
    "fmt"
    "log"
    "io/ioutil"
)

func main() {

    err := http.ListenAndServe(":9999", requestHandler{})
    
    log.Fatal(err)
    
}


type requestHandler struct{}


func (rh requestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    // PUTTING a document
    if r.Method == "PUT" {

        body, err := ioutil.ReadAll(r.Body)

        // Error reading the request body
        if err != nil {

            fmt.Fprintf(w, "Could not read request body")

        // Request body received
        } else {
            
            fmt.Fprintf(w, "PUT document to %s: %s\n", r.URL.Path, body)

        }

    }

    

}