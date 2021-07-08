package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port, avail := os.LookupEnv("PORT")
	if !avail {
		panic("no port number provided, server cannot start")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Hello!, %s!", r.URL.Path[1:])
	})
	http.ListenAndServe(fmt.Sprintf(":"+port), nil)
}
