package data

import (
	"encoding/json"
	"io"
	"os"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type Borough struct {
	Name    string
	Polygon []orb.Polygon
}

type BoroughGeoJSON struct {
	Type       string           `json:"type"`
	Geometry   geojson.Geometry `json:"geometry"`
	Properties struct {
		Name string `json:"borough"`
	} `json:"properties"`
}

var (
	Manhattan    = &Borough{Name: "Manhattan"}
	Brooklyn     = &Borough{Name: "Brooklyn"}
	Bronx        = &Borough{Name: "Bronx"}
	Queens       = &Borough{Name: "Queens"}
	StatenIsland = &Borough{Name: "Staten Island"}
	Boroughs     = []*Borough{Manhattan, Brooklyn, Bronx, Queens, StatenIsland}
)

func LoadBoroughGeoJSON(r io.Reader) error {
	var o []BoroughGeoJSON
	err := json.NewDecoder(r).Decode(&o)
	if err != nil {
		return err
	}
	for _, oo := range o {
		for _, b := range Boroughs {
			if b.Name == oo.Properties.Name {
				b.Polygon = append(b.Polygon, oo.Geometry.Coordinates.(orb.Polygon))
			}
		}
	}
	return nil
}

func LoadBoroughGeoJSONFromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return LoadBoroughGeoJSON(f)
}
