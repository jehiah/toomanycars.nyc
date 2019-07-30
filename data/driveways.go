package data

type Driveways struct {
	GuessSpaces int
	GuessCount  int
}

var DrivewayGuess = Driveways{700000, 350000}

func (d Driveways) Filter(b Borough) Driveways {
	return d
}
