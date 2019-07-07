package parkingdata

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// Based on @Pollytrott testimony at City Council
const BaseCurbParking int = 3000000

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

func (c Changes) EstimateSpaces() int {
	return BaseCurbParking + c.DeltaSpaces()
}

func (c Changes) DeltaSpaces() (spaces int) {
	for _, r := range c {
		spaces += r.Spaces
	}
	return
}

func ParseCurbChangesFromFile(file string) (Changes, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseCurbChanges(f)
}

func ParseCurbChanges(r io.Reader) (Changes, error) {
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}
	var data []Change
	for i, row := range records {
		if len(row) != 8 {
			return nil, fmt.Errorf("invalid row. expect 8 values got %#v", row)
		}
		if i == 0 {
			continue
		}
		t, err := time.Parse("2006/1/2", row[0])
		if err != nil {
			return nil, fmt.Errorf("row %d invalid date %q expected YYYY/mm/dd format. %s", i, row[0], err)
		}
		count, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, fmt.Errorf("row %d invalid space count %q %s", i, row[1], err)
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
	return data, nil
}
