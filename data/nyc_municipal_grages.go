package data

// https://data.cityofnewyork.us/Transportation/Municipal-Parking-Facilities-Manhattan-/i8d5-5ciu
// https://data.cityofnewyork.us/Transportation/Municipal-Parking-Facilities-Staten-Island-/udt3-taj7
// https://data.cityofnewyork.us/Transportation/Municipal-Parking-Facilities-Bronx-/7xqy-uv7r
// https://data.cityofnewyork.us/Transportation/Municipal-Parking-Facilities-Queens-/qfxy-c6k3
// https://data.cityofnewyork.us/Transportation/Municipal-Parking-Facilities-Brooklyn-/evdj-a5z2

type MunicipalGarage struct {
	Borough string
	Name    string
	Spaces  int
}
type MunicipalGarages []MunicipalGarage

func (m MunicipalGarages) Spaces() (total int) {
	for _, mm := range m {
		total += mm.Spaces
	}
	return
}

var AllMunicipalGarages MunicipalGarages = []MunicipalGarage{
	{"Manhattan", "Essex", 356},
	{"Brooklyn", "Avenue M Municipal Parking Field", 51},
	{"Brooklyn", "Bay Ridge Municipal Parking Garage", 205},
	{"Brooklyn", "Bensonhurst #1 Municipal Parking Field", 96},
	{"Brooklyn", "Bensonhurst #2 Municipal Parking Field", 24},
	{"Brooklyn", "Brighton Beach Municipal Parking Field", 312},
	{"Brooklyn", "Canarsie Municipal Parking Field", 281},
	{"Brooklyn", "Flatbush-Caton Municipal Parking Field", 52},
	{"Brooklyn", "Gowanus Municipal Parking Field", 0},
	{"Brooklyn", "Third Avenue between 30th and 41st Streets", 324},
	{"Brooklyn", "Grant Avenue Municipal Parking Field", 203},
	{"Brooklyn", "Sheepshead Bay #1 Municipal Parking Field", 60},
	{"Brooklyn", "Sheepshead Bay #2 Municipal Parking Field", 77},
	{"Queens", "Bayside Municipal Parking Field", 92},
	{"Queens", "Broadway-31st Street Municipal Parking Field", 61},
	{"Queens", "College Point Municipal Parking Field", 35},
	{"Queens", "Court Square Municipal Parking Garage", 703},
	{"Queens", "Ditmars #1 Municipal Parking Field", 57},
	{"Queens", "Ditmars #2 Municipal Parking Field", 67},
	{"Queens", "Far Rockaway #2 Municipal Parking Field", 70},
	{"Queens", "Flushing #2 Municipal Parking Field", 87},
	{"Queens", "Flushing #3 Municipal Parking Field", 156},
	{"Queens", "Flushing #4 Municipal Parking Field", 93},
	{"Queens", "Queens Borough Hall Municipal Parking Garage and Field", 948},
	{"Queens", "Queens Village Municipal Parking Field", 52},
	{"Queens", "Queens Civil Court Garage", 177},
	{"Queens", "Queens Family Court Garage", 207},
	{"Queens", "Rockaway Park Municipal Parking Field", 148},
	{"Queens", "Rosedale Municipal Parking Field", 164},
	{"Queens", "Steinway #1 Municipal Parking Field", 88},
	{"Queens", "Steinway #2 Municipal Parking Field", 46},
	{"Queens", "Sunnyside Municipal Parking Field", 494},
	{"Bronx", "Belmont Municipal Parking Field", 57},
	{"Bronx", "East 149th Street Municipal Parking Garage", 311},
	{"Bronx", "Jerome-190th Street Municipal Garage", 416},
	{"Bronx", "Jerome-Gun Hill Road Municipal Parking Garage", 240},
	{"Bronx", "White Plains Road Municipal Parking Field", 93},
	{"Staten Island", "Ferry Terminal South Municipal Parking Field", 222},
	{"Staten Island", "Great Kills Municipal Parking Field", 62},
	{"Staten Island", "New Dorp Municipal Parking Field", 75},
	{"Staten Island", "Staten Island Courthouse Garage and Parking Lot", 721},
}
