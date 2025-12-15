package main

import "net/http"

func main() {
	http.HandleFunc("/", Index)
	http.ListenAndServe(":8080", nil)
}

func Index(writer http.ResponseWriter, request *http.Request) {

}
