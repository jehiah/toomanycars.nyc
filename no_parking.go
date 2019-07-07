package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"sort"

	parkingdata "github.com/jehiah/🚫🚗.nyc/data"
)

type Data struct {
	InitialParkingSpaces int
	CurbParking          parkingdata.Changes
	DCA                  parkingdata.DCALicenses
	ParkingLot           parkingdata.ParkingLots
}

func (d Data) RecentChanges() parkingdata.Changes {
	var o parkingdata.Changes
	o = d.CurbParking
	o = append(o, d.DCA.RecentChanges()...)
	sort.Slice(o, func(i, j int) bool { return o[i].EffectiveDate.After(o[j].EffectiveDate) })

	return o
}

func (d Data) ParkingSpaces() int {
	return d.InitialParkingSpaces + d.CurbParking.DeltaSpaces() + d.DCA.Spaces() + d.ParkingLot.EstimateSpaces()
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
		"millionsqft": func(n float64) float64 {
			return n / 1000000.0
		},
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
	doittParkingLot, err := parkingdata.ParseDOITTParkingLotFromFile("data/DOITT_planimetrics_parking_lot.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("DOT changes %d", curbParking.DeltaSpaces())
	log.Printf("DCA managed spaces %d", dca.Spaces())
	log.Printf("DOITT planimetrics Parking Lots %d lots covering %.f sqft. Estimated %d spaces", len(doittParkingLot), doittParkingLot.SurfaceArea(), doittParkingLot.EstimateSpaces())

	w, err := os.Create("www/index.html")
	defer w.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = t.ExecuteTemplate(w, "index.html", Data{
		InitialParkingSpaces: 3000000 - 600000,
		CurbParking:          curbParking,
		DCA:                  dca,
		ParkingLot:           doittParkingLot,
	})
	if err != nil {
		log.Fatal(err)
	}
}
