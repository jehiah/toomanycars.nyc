package data

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

func (d DCALicense) addressKey() string {
	return fmt.Sprintf("%s %s %s", d.AddressBuilding, d.AddressStreetName, d.Borough)
}

func (d DCALicense) Change() Change {
	src := fmt.Sprintf("License %s", d.LicenseNumber)
	dt := d.Creation
	s := d.Spaces()
	if d.LicenseStatus == "Inactive" {
		dt = d.Expiration
		s = -1 * s
		src = fmt.Sprintf("Expired License %s (%s)", d.LicenseNumber, d.Creation.Format("2006"))
	} else {
		src = "New " + src
	}
	n := d.BusinessName
	if d.BusinessName2 != "" {
		n = d.BusinessName2
	}
	addr := fmt.Sprintf("%s %s", d.AddressBuilding, d.AddressStreetName)
	return Change{
		EffectiveDate: dt,
		Spaces:        s,
		Borough:       d.Borough,
		Category:      d.Industry,
		Name:          n,
		Description:   addr,
		Source:        src,
		ReferenceURL:  fmt.Sprintf("https://www.google.com/maps/place/%s", strings.ReplaceAll(fmt.Sprintf("%s, %s, NY", addr, d.Borough), " ", "+")),
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
	d.Creation, err = time.Parse("2006-01-02T03:04:05", data.CreationStr)
	if err != nil {
		return err
	}
	if data.ExpirationStr != "" {
		d.Expiration, err = time.Parse("2006-01-02T03:04:05", data.ExpirationStr)
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
	for _, dd := range d.Active() {
		spaces += dd.Spaces()
	}
	return
}
func (d DCALicenses) EstimateLotSpaces() (spaces int) {
	for _, dd := range d.Active() {
		switch dd.Industry {
		case "Garage and Parking Lot":
			spaces += (dd.Spaces() / 2)
		case "Parking Lot":
			spaces += dd.Spaces()
		}
	}
	return
}
func (d DCALicenses) Active() DCALicenses {
	var o DCALicenses
	now := time.Now()
	for _, dd := range d {
		if dd.LicenseStatus == "Inactive" && dd.Expiration.Before(now) {
			continue
		}
		o = append(o, dd)
	}
	return o
}

func (d DCALicenses) RecentChanges() Changes {
	// build skip list
	skip := make(map[string]bool)
	for _, e := range replacedLicenses {
		skip[e.ReplacesLicense] = true
	}

	// build a list of most recent address start / end
	recentExpired := make(map[string]DCALicense)
	recentNew := make(map[string]DCALicense)
	for _, dd := range d {
		addr := dd.addressKey()
		switch dd.LicenseStatus {
		case "Active":
			if c, ok := recentNew[addr]; !ok || c.Creation.Before(dd.Creation) {
				recentNew[addr] = dd
			}
		case "Inactive":
			if c, ok := recentExpired[addr]; !ok || c.Expiration.Before(dd.Expiration) {
				recentExpired[addr] = dd
			}
		}
	}

	var o Changes
	// 3 months
	cutoff := time.Now().AddDate(0, -3, 0)
	for _, dd := range d {
		if skip[dd.LicenseNumber] {
			continue
		}
		switch dd.LicenseStatus {
		case "Inactive":
			if dd.Expiration.After(time.Now()) {
				// consider it as still activ
				continue
			}
			if dd.Expiration.Before(cutoff) {
				continue
			}
			// expired w/ a more recent expired is skipped
			if e, ok := recentExpired[dd.addressKey()]; ok && dd.Expiration.Before(e.Expiration) {
				continue
			}
			// "expired" w/ a later start at that address is skipped
			if _, ok := recentNew[dd.addressKey()]; !ok {
				o = append(o, dd.Change())
			}
		case "Active":
			if dd.Creation.Before(cutoff) {
				continue
			}
			// "new" w/ a previous entry at that address isn't an increase (unless space count changes)
			if _, ok := recentExpired[dd.addressKey()]; ok {
				// TODO: if e.Space() != dd.Space()
			} else {
				o = append(o, dd.Change())
			}
		}
	}
	return o
}
