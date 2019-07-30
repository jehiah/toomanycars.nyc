package data

import (
	"fmt"
)

type Driveways struct {
	GuessSpaces int
	GuessCount  int
}

var DrivewayGuess = Driveways{700000, 350000}

func (d Driveways) Scale(f float64) Driveways {
	d.GuessSpaces = int(float64(d.GuessSpaces) * f)
	d.GuessCount = int(float64(d.GuessCount) * f)
	return d
}

func (d Driveways) Filter(b Borough) Driveways {
	switch b.Name {
	case "Manhattan":
		return d.Scale(884 / 375840.0)
	case "Brooklyn":
		return d.Scale(95763 / 375840.0)
	case "Bronx":
		return d.Scale(30501 / 375840.0)
	case "Queens":
		return d.Scale(215338 / 375840.0)
	case "Staten Island":
		return d.Scale(33291 / 375840.0)
	default:
		panic(fmt.Sprintf("unknown borough %s", b.Name))
	}
}
