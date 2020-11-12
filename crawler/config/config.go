package config

const (
	// Parser names
	ParseCity     = "ParseCity"
	ParseCityList = "ParseCityList"
	ParseProfile  = "ParseProfile"

	ParseCarDetail = "ParseCarDetail"
	ParseCarList   = "ParseCarList"
	ParseCarModel  = "ParseCarModel"

	NilParser = "NilParser"

	// ElasticSearch
	CarElasticIndex    = "car_profile"
	DatingElasticIndex = "dating_profile"

	// Rate limiting
	Qps = 2
)
