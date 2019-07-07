package parkingdata

import (
	"encoding/json"
	"strconv"
)

// StrFloat64 Unmarshal's from a string. It will not round trip json encoding/unmarshaling
type StrFloat64 float64

func (s *StrFloat64) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return err
	}
	*s = StrFloat64(f)
	return nil
}
