package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// Based on @Pollytrott testimony at City Council, later clarified as 2.990
// https://legistar.council.nyc.gov/View.ashx?M=F&ID=7538481&GUID=8A898206-4A26-40E1-B432-036D4B3112A5
const BaseOnStreetParking int = 2990000

type OnStreet struct {
	BaseEst int
	Changes
}

func (s OnStreet) Filter(b Borough) OnStreet {
	var o OnStreet
	// These ratios come from DOT
	//
	//
	// For comparison ratios from https://github.com/jehiah/analyze-nyc-parking-signs
	// While that process yields a different estimate than DOT publicly stated,
	// we can use it's per-borough breakdown
	//
	// 1595744 total car spaces
	// 207429 cars manhattan .129988895
	// 602270 cars broklyn .377422694
	// 244868 cars bronx .153450679
	// 436070 cars queens .27327065
	// 105107 cars Staten Island .065867081
	switch b.Name {
	case "Manhattan":
		o.BaseEst = int(float64(s.BaseEst) * 0.063545151)
	case "Bronx":
		o.BaseEst = int(float64(s.BaseEst) * 0.140468227)
	case "Brooklyn":
		o.BaseEst = int(float64(s.BaseEst) * 0.247491639)
	case "Queens":
		o.BaseEst = int(float64(s.BaseEst) * 0.401337793)
	case "Staten Island":
		o.BaseEst = int(float64(s.BaseEst) * 0.147157191)
	default:
		panic("unknown borough")
	}
	o.Changes = s.Changes.Filter(b)
	return o
}

type Change struct {
	EffectiveDate time.Time
	Spaces        int
	Borough       string
	Category      string
	Name          string
	Description   string
	ReferenceURL  string
	Source        string
}

func (c Change) Future() bool {
	return time.Now().Before(c.EffectiveDate)
}

type Changes []Change

func (c Changes) Filter(b Borough) Changes {
	var o Changes
	for _, cc := range c {
		if cc.Borough == b.Name {
			o = append(o, cc)
		}
	}
	return o
}

func (c OnStreet) EstimateSpaces() int {
	return c.BaseEst + c.DeltaSpaces()
}

func (c Changes) DeltaSpaces() (spaces int) {
	for _, r := range c {
		spaces += r.Spaces
	}
	return
}

func ParseCurbChangesFromFile(file string) (OnStreet, error) {
	f, err := os.Open(file)
	if err != nil {
		return OnStreet{}, err
	}
	defer f.Close()
	return ParseCurbChanges(f)
}

func ParseCurbChanges(r io.Reader) (OnStreet, error) {
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return OnStreet{}, err
	}
	var data []Change
	for i, row := range records {
		if len(row) != 8 {
			return OnStreet{}, fmt.Errorf("invalid row. expect 8 values got %#v", row)
		}
		if i == 0 {
			continue
		}
		t, err := time.Parse("2006/1/2", row[0])
		if err != nil {
			return OnStreet{}, fmt.Errorf("row %d invalid date %q expected YYYY/mm/dd format. %s", i, row[0], err)
		}
		count, err := strconv.Atoi(row[1])
		if err != nil {
			return OnStreet{}, fmt.Errorf("row %d invalid space count %q %s", i, row[1], err)
		}
		data = append(data, Change{
			EffectiveDate: t,
			Spaces:        count,
			Borough:       row[2],
			Category:      row[3],
			Name:          row[4],
			Description:   row[5],
			ReferenceURL:  row[6],
			Source:        row[7],
		})
	}
	return OnStreet{BaseOnStreetParking, data}, nil
}
