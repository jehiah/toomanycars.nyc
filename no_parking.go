package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
)

type Data struct {
	ParkingSpaces int
}

func tokenString(s string) []string {
	var o []string
	for _, ss := range s {
		o = append(o, fmt.Sprintf("%c", ss))
	}
	return o
}

func main() {
	flag.Parse()

	funcMap := template.FuncMap{
		"tokenString": tokenString,
	}

	t, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	w, err := os.Create("www/index.html")
	defer w.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = t.ExecuteTemplate(w, "index.html", Data{3000000})
	if err != nil {
		log.Fatal(err)
	}
}
