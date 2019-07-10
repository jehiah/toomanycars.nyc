package data

import (
	"encoding/json"
	"io"
	"math"
	"os"
	"strconv"
)

type ParkingLot struct {
	ID          string  `json:"source_id"`
	Status      string  `json:"status"`
	ShapeLength float64 `json:"shape_length"`
	ShapeArea   float64 `json:"shape_area"`
}
type ParkingLots []ParkingLot

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

// UnmarshalJSON converts shape_leng, shape_area into float64
func (p *ParkingLot) UnmarshalJSON(b []byte) error {
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
		p = &ParkingLot{}
	}
	p.ID = data.SourceID
	p.Status = data.Status
	p.ShapeLength, err = strconv.ParseFloat(data.ShapeLength, 64)
	if err != nil {
		return err
	}
	p.ShapeArea, err = strconv.ParseFloat(data.ShapeArea, 64)
	if err != nil {
		return err
	}
	return nil
}

func ParseDOITTParkingLot(r io.Reader) (ParkingLots, error) {
	var o []ParkingLot
	err := json.NewDecoder(r).Decode(&o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func ParseDOITTParkingLotFromFile(file string) (ParkingLots, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseDOITTParkingLot(f)
}
