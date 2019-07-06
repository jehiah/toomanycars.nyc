package parkingdata

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type DCALicense struct {
	LicenseNumber     string `json:"license_nbr"`
	LicenseStatus     string `json:"license_status"`
	Creation          time.Time
	Expiration        time.Time
	Industry          string `json:"industry"` // Garage, Parking Lot, Garage and Parking Lot
	BusinessName      string `json:"business_name"`
	BusinessName2     string `json:"business_name_2"`
	AddressBuilding   string `json:"address_building"`
	AddressStreetName string `json:"address_street_name"`
	Borough           string `json:"address_borough"`
	Detail            string `json:"detail_2"` // "Vehicle Spaces: %d, Bicycle Spaces: %d"
}

func (d DCALicense) Change() Change {
	dt := d.Creation
	s := d.Spaces()
	if d.LicenseStatus == "Inactive" {
		dt = d.Expiration
		s = -1 * s
	}
	n := d.BusinessName
	if d.BusinessName2 != "" {
		n = d.BusinessName2
	}
	return Change{
		EffectiveDate: dt,
		Spaces:        s,
		Borough:       d.Borough,
		Category:      d.Industry,
		Name:          fmt.Sprintf("%s %s %s", n, d.AddressBuilding, d.AddressStreetName),
		Source:        fmt.Sprintf("DCA License %s", d.LicenseNumber),
	}
}

func (d *DCALicense) UnmarshalJSON(b []byte) error {
	type localLicense DCALicense
	type tempType struct {
		localLicense
		CreationStr   string `json:"license_creation_date"`
		ExpirationStr string `json:"lic_expir_dd"`
	}
	var data tempType
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	if d == nil {
		d = &DCALicense{}
	}
	*d = DCALicense(data.localLicense)
	d.Creation, err = time.Parse("2006-01-02T03:04:05.000", data.CreationStr)
	if err != nil {
		return err
	}
	if data.ExpirationStr != "" {
		d.Expiration, err = time.Parse("2006-01-02T03:04:05.000", data.ExpirationStr)
		if err != nil {
			return err
		}
	}
	return nil
}

type DCALicenses []DCALicense

// Spaces parses the Detail field for the number of Vehcile Spaces
func (d DCALicense) Spaces() int {
	f := func(c rune) bool {
		return unicode.IsSpace(c) || c == ','
	}
	c := strings.FieldsFunc(d.Detail, f)
	if len(c) < 3 {
		return 0
	}
	n, _ := strconv.Atoi(c[2])
	return n
}

func ParseDCA(r io.Reader) (DCALicenses, error) {
	var o []DCALicense
	err := json.NewDecoder(r).Decode(&o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func ParseDCAFromFile(file string) (DCALicenses, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseDCA(f)
}

func (d DCALicenses) Spaces() (spaces int) {
	for _, dd := range d {
		if dd.LicenseStatus != "Active" {
			continue
		}
		spaces += dd.Spaces()
	}
	return
}
func (d DCALicenses) RecentChanges() Changes {
	var o Changes
	cutoff := time.Now().AddDate(0, -12, 0)
	for _, dd := range d {
		switch dd.LicenseStatus {
		case "Active":
			if dd.Creation.After(cutoff) {
				o = append(o, dd.Change())
			}
		case "Inactive":
			if dd.Expiration.After(cutoff) {
				o = append(o, dd.Change())
			}
		}
	}
	return o
}
