package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/jehiah/toomanycars.nyc/data"
	"golang.org/x/text/message"
)

type Data struct {
	OnStreet         data.OnStreet
	DCA              data.DCALicenses
	ParkingLot       data.ParkingLots
	PrivateGarages   data.Garages
	MunicipalGarages data.MunicipalGarages
	Driveways        data.Driveways
	Boroughs         []*data.Borough
	Updated          time.Time
	BoroughCounter   Counter
	Timeframes       []time.Time
}

// func (d Data) RecentChanges() data.Changes {
// 	var o data.Changes
// 	o = d.OnStreet.Changes
// 	o = append(o, d.DCA.RecentChanges()...)
// 	sort.Slice(o, func(i, j int) bool {
// 		if o[i].EffectiveDate.Equal(o[j].EffectiveDate) {
// 			return strings.Compare(o[i].Name, o[j].Name) == -1
// 		}
// 		return o[i].EffectiveDate.After(o[j].EffectiveDate)
//
// 	})
//
// 	return o
// }

// RecentTimeframes returns the 1st of the month for the past 12 months
func RecentTimeframes() []time.Time {
	var o []time.Time
	y, m, _ := time.Now().Date()
	start := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 12; i++ {
		t := start.AddDate(0, -1*i, 0)
		o = append(o, t)
	}
	return o
}

func (d Data) ParkingSpaces() int {
	spaces := d.OnStreet.EstimateSpaces()
	spaces += d.DCA.Spaces()
	spaces += d.ParkingLot.EstimateSpaces(d.DCA.EstimateLotSpaces())
	spaces += d.PrivateGarages.EstimateSpaces()
	spaces += d.MunicipalGarages.Spaces()
	spaces += d.Driveways.GuessSpaces
	return spaces
}

type Counter map[string]int

func (c Counter) Add(b data.Borough, n int) string {
	c[b.Name] += n
	return ""
}
func (c Counter) Filter(b data.Borough) int {
	return c[b.Name]
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
		"iscomma": func(idx int, s []string) bool {
			l := len(s)
			if idx > 0 && idx != l && ((l-idx)%3 == 0) {
				return true
			}
			return false
		},
	}

	t, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html", "templates/dca_changes.html")
	if err != nil {
		log.Fatal(err)
	}

	if err = data.LoadBoroughGeoJSONFromFile("data/borough_boundaries.geojson"); err != nil {
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
	doittParkingLot, err := data.ParseDOITTParkingLotFromFile("data/DOITT_planimetrics_parking_lot.geojson")
	if err != nil {
		log.Fatal(err)
	}
	doittPrivateGarages, err := data.ParseDOITTGaragesFromFile("data/DOITT_planimetrics_building_garages.geojson")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("DOT OnStreet changes %d (Estimated %d spaces)", curbParking.DeltaSpaces(), curbParking.EstimateSpaces())
	log.Printf("DCA managed spaces %d, businesses: %d", dca.Spaces(), len(dca.Active()))
	log.Printf("DoITT planimetrics Parking Lots %d lots covering %.f sqft. Estimated %d spaces", len(doittParkingLot), doittParkingLot.SurfaceArea(), doittParkingLot.EstimateSpaces(dca.EstimateLotSpaces()))
	log.Printf("DoITT planimetrics Private Garages %d covering %.f sqft. Estimated %d spaces", len(doittPrivateGarages), doittPrivateGarages.SurfaceArea(), doittPrivateGarages.EstimateSpaces())
	log.Printf("%d Municipal Garages with %d spaces", len(data.AllMunicipalGarages), data.AllMunicipalGarages.Spaces())

	for _, b := range data.Boroughs {
		log.Printf("\n")
		d := dca.Filter(*b)
		log.Printf("%s DCA managed spaces %d, businesses: %d", b.Name, d.Spaces(), len(d.Active()))
		dp := doittParkingLot.Filter(*b)
		log.Printf("%s DoITT planimetrics Parking Lots %d lots covering %.f sqft. Estimated %d spaces", b.Name, len(dp), dp.SurfaceArea(), dp.EstimateSpaces(dca.Filter(*b).EstimateLotSpaces()))
		dpg := doittPrivateGarages.Filter(*b)
		log.Printf("%s DoITT planimetrics Private Garages %d covering %.f sqft. Estimated %d spaces", b.Name, len(dpg), dpg.SurfaceArea(), dpg.EstimateSpaces())
		mg := data.AllMunicipalGarages.Filter(*b)
		log.Printf("%s %d Municipal Garages with %d spaces", b.Name, len(mg), mg.Spaces())
	}

	est, _ := time.LoadLocation("America/New_York")

	g := dca.Group()
	for _, d := range RecentTimeframes() {
		filename := fmt.Sprintf("www/dca_%s.html", d.Format("200601"))
		log.Printf("rendering %s", filename)
		w, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		err = t.ExecuteTemplate(w, "dca_changes.html", struct {
			Changes []data.Change
			Date    time.Time
			Updated time.Time
		}{
			Changes: g.ChangesInMonth(d),
			Date:    d,
			Updated: time.Now().In(est),
		})
		if err != nil {
			log.Fatal(err)
		}
		w.Close()
	}

	log.Print("rendering www/index.html")
	w, err := os.Create("www/index.html")
	defer w.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = t.ExecuteTemplate(w, "index.html", Data{
		OnStreet:         curbParking,
		DCA:              dca,
		ParkingLot:       doittParkingLot,
		PrivateGarages:   doittPrivateGarages,
		MunicipalGarages: data.AllMunicipalGarages,
		Driveways:        data.DrivewayGuess,
		Updated:          time.Now().In(est),
		Boroughs:         data.Boroughs,
		BoroughCounter:   make(Counter),
		Timeframes:       RecentTimeframes(),
	})
	if err != nil {
		log.Fatal(err)
	}

}
