package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"sort"

	"github.com/jehiah/toomanycars.nyc/data"
	"golang.org/x/text/message"
)

type Data struct {
	OnStreet       data.Changes
	DCA            data.DCALicenses
	ParkingLot     data.ParkingLots
	PrivateGarages data.Garages
}

func (d Data) RecentChanges() data.Changes {
	var o data.Changes
	o = d.OnStreet
	o = append(o, d.DCA.RecentChanges()...)
	sort.Slice(o, func(i, j int) bool { return o[i].EffectiveDate.After(o[j].EffectiveDate) })

	return o
}

func (d Data) ParkingSpaces() int {
	return d.OnStreet.EstimateSpaces() + d.DCA.Spaces() + d.ParkingLot.EstimateSpaces() + d.PrivateGarages.EstimateSpaces()
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
	p := message.NewPrinter(message.MatchLanguage("en"))

	funcMap := template.FuncMap{
		"tokenString": tokenString,
		"millionsqft": func(n float64) float64 {
			return n / 1000000.0
		},
		"commify": func(v interface{}) string {
			switch v.(type) {
			case int:
				return p.Sprintf("%d", v)
			case float64:
				return p.Sprintf("%.1f", v)
			default:
				panic(fmt.Sprintf("unknown type %T for %#v", v, v))
			}
		},
	}

	t, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	curbParking, err := data.ParseCurbChangesFromFile("data/curb_changes.csv")
	if err != nil {
		log.Fatal(err)
	}
	dca, err := data.ParseDCAFromFile("data/dca_licenses.json")
	if err != nil {
		log.Fatal(err)
	}
	doittParkingLot, err := data.ParseDOITTParkingLotFromFile("data/DOITT_planimetrics_parking_lot.json")
	if err != nil {
		log.Fatal(err)
	}
	doittPrivateGarages, err := data.ParseDOITTGaragesFromFile("data/DOITT_planimetrics_building_garages.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("DOT OnStreet changes %d (Estimated %d spaces)", curbParking.DeltaSpaces(), curbParking.EstimateSpaces())
	log.Printf("DCA managed spaces %d", dca.Spaces())
	log.Printf("DOITT planimetrics Parking Lots %d lots covering %.f sqft. Estimated %d spaces", len(doittParkingLot), doittParkingLot.SurfaceArea(), doittParkingLot.EstimateSpaces())
	log.Printf("DOITT planimetrics Private Garages %d covering %.f sqft. Estimated %d spaces", len(doittPrivateGarages), doittPrivateGarages.SurfaceArea(), doittPrivateGarages.EstimateSpaces())

	w, err := os.Create("www/index.html")
	defer w.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = t.ExecuteTemplate(w, "index.html", Data{
		OnStreet:       curbParking,
		DCA:            dca,
		ParkingLot:     doittParkingLot,
		PrivateGarages: doittPrivateGarages,
	})
	if err != nil {
		log.Fatal(err)
	}
}
