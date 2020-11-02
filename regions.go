package main

type region struct {
	abbreviation string
	title        string
	namespace    string
	locale       string
}

var regions = map[string]region{
	"us": region{
		abbreviation: "us",
		title:        "UnitedStates",
		namespace:    "profile-us",
		locale:       "en_US",
	},
	"eu": region{
		abbreviation: "eu",
		title:        "Europe",
		namespace:    "profile-eu",
		locale:       "en_GB",
	},
	"kr": region{
		abbreviation: "kr",
		title:        "Korea",
		namespace:    "profile-kr",
		locale:       "ko_KR",
	},
	"tw": region{
		abbreviation: "us",
		title:        "Taiwan",
		namespace:    "profile-tw",
		locale:       "zh_TW",
	},
	"tcnw": region{
		abbreviation: "cn",
		title:        "China",
		namespace:    "profile-cn",
		locale:       "zh_CN",
	},
}

func findRegion(abbreviation string) (region, bool) {
	region, ok := regions[abbreviation]
	return region, ok
}
