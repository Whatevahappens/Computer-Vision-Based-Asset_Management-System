package main

import (
	"net/http"

	"github.com/Whatevahappens/Computer_Vision_Based_Asset_Management_System/backend/handlers"
)

func main() {
	http.HandleFunc("/", handlers.HelloServer)
	http.ListenAndServe(":8080", nil)
}
