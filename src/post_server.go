package main

import (
	"fmt"
	//"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received")
	fmt.Println("Method: " + r.Method)
	var contentLength int64 = r.ContentLength
	fmt.Println("Content length: " + strconv.FormatInt(contentLength, 10))

	if contentLength != 0 {
		body := r.Body
		image_data, _ := ioutil.ReadAll(body)

		ioutil.WriteFile("body_data", image_data, 0777)
	}

}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
