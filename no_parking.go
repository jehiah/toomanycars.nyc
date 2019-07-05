package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"

	parkingdata "github.com/jehiah/🚫🚗.nyc/data"
)

type Data struct {
	ParkingSpaces int
	Changes       parkingdata.Changes
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
	f, err := os.Open("data/parking_spaces.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	changes, err := parkingdata.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	err = t.ExecuteTemplate(w, "index.html", Data{
		ParkingSpaces: changes.ProjectedTotal(3000000),
		Changes:       changes,
	})
	if err != nil {
		log.Fatal(err)
	}
}
