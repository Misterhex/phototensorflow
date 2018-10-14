package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/gorilla/handlers"
)

func main() {

	r := http.NewServeMux()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html :=
			`
			<html>
			<body> 
			<h1>upload photo</h1> 
			<form action="/upload" method="post" enctype="multipart/form-data">
				<input type="file" name="image" accept="image/*" capture/>
				<br/>
				<input type="submit" value="Upload">
		  	</form>
			</body> 
			</html>
			`
		fmt.Fprintf(w, html)
	})

	r.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseMultipartForm(5 * 1024 * 1024)
			if err != nil {
				handleError(w, err)
			}

			f, _, err := r.FormFile("image")
			if err != nil {
				handleError(w, err)
			}

			defer f.Close()

			data, err := ioutil.ReadAll(f)
			if err != nil {
				handleError(w, err)
			}

			log.Println("classifying image ...")
			r := bytes.NewReader(data)
			out, err := ClassifyImage(r)
			if err != nil {
				log.Printf("error when classifying image, %s...\n", err.Error())
				handleError(w, err)
			}
			log.Printf("classify image result: %s\n", out)

			base64String := base64.StdEncoding.EncodeToString(data)

			templateData := struct {
				Base64              string
				ClassifyImageResult string
			}{
				base64String,
				out,
			}

			template := template.Must(template.ParseFiles("result.html"))

			template.Execute(w, templateData)

		} else {
			w.WriteHeader(404)
			fmt.Fprintf(w, "invalid")
		}
	})

	port := ":3001"
	log.Printf("server listening port %s\n", port)
	log.Fatal(http.ListenAndServe(":3001", handlers.LoggingHandler(os.Stdout, r)))
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	fmt.Fprintf(w, err.Error())
}
