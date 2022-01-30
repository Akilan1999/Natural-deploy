package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// content holds our static web server content.
////go:embed css/*
//var content embed.FS

//go:embed index.html
var f string

func main() {

	fmt.Println(f)

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, "")
		//fmt.Fprintf(w, "Hello!")
	})

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
