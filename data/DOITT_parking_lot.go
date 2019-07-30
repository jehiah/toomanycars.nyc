package data

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

type ParkingLot struct {
	Type                 string            `json:"type"`
	Geometry             *geojson.Geometry `json:"geometry"`
	ParkingLotProperties `json:"properties"`
}
type ParkingLotProperties struct {
	ID          string  `json:"source_id"`
	Status      string  `json:"status"`
	ShapeLength float64 `json:"shape_length"`
	ShapeArea   float64 `json:"shape_area"`
}
type ParkingLots []*ParkingLot

func (p ParkingLots) Filter(b Borough) ParkingLots {
	var o ParkingLots
Loop:
	for _, pp := range p {
		if pp.Geometry == nil {
			log.Fatalf("nil Geometry %#v", pp)
		}
		center := pp.Geometry.Geometry().Bound().Center()
		for _, subshape := range b.Polygon {
			if planar.PolygonContains(subshape, center) {
				// log.Printf("%s[%d] contains %s at %#v", b.Name, i, pp.ID, center)
				o = append(o, pp)
				continue Loop
			}
		}
		// log.Printf("%s not contained in %s", pp.ID, b.Name)
	}
	return o
}

func (p ParkingLots) SurfaceArea() (total float64) {
	for _, pp := range p {
		total += pp.ShapeArea
	}
	return
}

// EstimateSpaces excluding the provided number of spaces
// exclusion allows accounting for double count of DCA licensed lots
func (p ParkingLots) EstimateSpaces(exclude int) (spaces int) {
	spaces -= exclude
	for _, pp := range p {
		spaces += pp.EstimateSpaces()
	}
	return
}

func (p ParkingLot) EstimateSpaces() int {
	estimate := p.ShapeArea / 350
	return int(math.Floor(estimate))
}

func (g *ParkingLot) UnmarshalJSON(b []byte) error {
	type tempType struct {
		Type                 string            `json:"type"`
		Geometry             *geojson.Geometry `json:"geometry"`
		ParkingLotProperties json.RawMessage   `json:"properties"`
	}
	var data tempType
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	if g == nil {
		g = &ParkingLot{}
	}
	(*g).Type = data.Type
	(*g).Geometry = data.Geometry
	return json.Unmarshal(data.ParkingLotProperties, &g.ParkingLotProperties)
}

// UnmarshalJSON converts shape_leng, shape_area into float64
func (p *ParkingLotProperties) UnmarshalJSON(b []byte) error {
	type tempType struct {
		SourceID    string `json:"source_id"`
		Status      string `json:"status"`
		ShapeLength string `json:"shape_leng"`
		ShapeArea   string `json:"shape_area"`
	}
	var data tempType
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	if p == nil {
		p = &ParkingLotProperties{}
	}
	p.ID = data.SourceID
	p.Status = data.Status
	p.ShapeLength, err = strconv.ParseFloat(data.ShapeLength, 64)
	if err != nil {
		return fmt.Errorf("%s for %#v %s", err, data, string(b))
	}
	p.ShapeArea, err = strconv.ParseFloat(data.ShapeArea, 64)
	if err != nil {
		return fmt.Errorf("%s for %#v %s", err, data, string(b))
	}
	return nil
}

func ParseDOITTParkingLot(r io.Reader) (ParkingLots, error) {
	type FeatureCollection struct {
		Features []*ParkingLot `json:"features"`
	}
	var o FeatureCollection
	err := json.NewDecoder(r).Decode(&o)
	if err != nil {
		return nil, err
	}
	log.Printf("Lot: %#v", o.Features[0])
	return o.Features, nil
}

func ParseDOITTParkingLotFromFile(file string) (ParkingLots, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseDOITTParkingLot(f)
}
