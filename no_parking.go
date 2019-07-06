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
	InitialParkingSpaces int
	CurbParking          parkingdata.Changes
	DCA                  parkingdata.DCALicenses
}

func (d Data) RecentChanges() parkingdata.Changes {
	var o parkingdata.Changes
	o = d.CurbParking
	o = append(o, d.DCA.RecentChanges()...)
	return o
}

func (d Data) ParkingSpaces() int {
	return d.InitialParkingSpaces + d.CurbParking.DeltaSpaces() + d.DCA.Spaces()
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
	curbParking, err := parkingdata.ParseCurbChangesFromFile("data/curb_changes.csv")
	if err != nil {
		log.Fatal(err)
	}
	dca, err := parkingdata.ParseDCAFromFile("data/dca_licenses.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("DOT changes %d", curbParking.DeltaSpaces())
	log.Printf("DCA managed spaces %d", dca.Spaces())

	w, err := os.Create("www/index.html")
	defer w.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = t.ExecuteTemplate(w, "index.html", Data{
		InitialParkingSpaces: 3000000 - 600000,
		CurbParking:          curbParking,
		DCA:                  dca,
	})
	if err != nil {
		log.Fatal(err)
	}
}
