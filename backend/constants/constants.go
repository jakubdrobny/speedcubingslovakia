package constants

import "math"

const (
	EPS                        = 1e-9
	DNF                        = math.MaxInt - 5
	DNS                        = math.MaxInt - 4
	VERY_SLOW                  = math.MaxInt - 20
	MBLD_MAX_CUBES_PER_ATTEMPT = 69
	PR_COLOR                   = "#FC4A0A"
	NR_COLOR                   = "#28A745"
	CR_COLOR                   = "#D00404"
	WR_COLOR                   = "#0B6BCB"
)

var US_STATE_NAMES = []string{
	"Alabama",
	"Alaska",
	"Arizona",
	"Arkansas",
	"California",
	"Colorado",
	"Connecticut",
	"Delaware",
	"Florida",
	"Georgia",
	"Hawaii",
	"Idaho",
	"Illinois",
	"Indiana",
	"Iowa",
	"Kansas",
	"Kentucky",
	"Louisiana",
	"Maine",
	"Maryland",
	"Massachusetts",
	"Michigan",
	"Minnesota",
	"Mississippi",
	"Missouri",
	"Montana",
	"Nebraska",
	"Nevada",
	"New Hampshire",
	"New Jersey",
	"New Mexico",
	"New York",
	"North Carolina",
	"North Dakota",
	"Ohio",
	"Oklahoma",
	"Oregon",
	"Pennsylvania",
	"Rhode Island",
	"South Carolina",
	"South Dakota",
	"Tennessee",
	"Texas",
	"Utah",
	"Vermont",
	"Virginia",
	"Washington",
	"West Virginia",
	"Wisconsin",
	"Wyoming",
}

var COUNTRY_GROUPS_ISO2 = []string{"XA", "XE", "XF", "XM", "XN", "XO", "XS", "XW"}

var CONTINENT_ID_TO_COUNTRY_GROUP_ISO2 = map[string]string{
	"_Asia":          "XA",
	"_Europe":        "XE",
	"_Africa":        "XF",
	"_North America": "XN",
	"_Oceania":       "XO",
	"_South America": "XS",
}
