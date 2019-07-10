package data

type DCALink struct {
	License         string
	ReplacesLicense string
}

var replacedLicenses = []DCALink{
	{"2082736-DCA", "1288781-DCA"},
	{"2083664-DCA", "2069760-DCA"},
	{"2076276-DCA", "1383664-DCA"},
	{"2069015-DCA", "1117864-DCA"},
	{"2060732-DCA", "2024225-DCA"},
}
