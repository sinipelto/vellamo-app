package main

type AreaGeoJsonData struct {
	Type     string   `json:"type"`
	Bbox     []string `json:"bbox"`
	Geometry struct {
		Type        string          `json:"type"`
		Coordinates [][][][]float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Slug       string  `json:"slug"`
		Name       string  `json:"name"`
		Rusko      float64 `json:"rusko"`
		Kaupinoja  float64 `json:"kaupinoja"`
		Messukyla  float64 `json:"messukyla"`
		Pinsio     float64 `json:"pinsio"`
		Julkujarvi float64 `json:"julkujarvi"`
		Mustalampi float64 `json:"mustalampi"`
		Hyhky      float64 `json:"hyhky"`
		Ph         float64 `json:"ph"`
		T          float64 `json:"t"`
		Cl         float64 `json:"cl"`
		Retention  float64 `json:"retention"`
		Pohjavesi  float64 `json:"pohjavesi"`
		Pintavesi  float64 `json:"pintavesi"`
		Hardness   float64 `json:"hardness"`
	} `json:"properties"`
}

type AreaSensorTemperatureHourlyData struct {
	Area              string  `json:"area"`
	Measurement       string  `json:"measurement"`
	Interval          string  `json:"interval"`
	Minimum           float64 `json:"minimum"`
	Maximum           float64 `json:"maximum"`
	LowerLimit        float64 `json:"lower_limit"`
	UpperLimit        float64 `json:"upper_limit"`
	RecommendedDigits int     `json:"recommended_digits"`
	Unit              string  `json:"unit"`
	Groups            []struct {
		Title  string `json:"title"`
		Values []struct {
			Datetime string  `json:"datetime"`
			Average  float64 `json:"average"`
			Minimum  float64 `json:"minimum"`
			Maximum  float64 `json:"maximum"`
		} `json:"values"`
	} `json:"groups"`
}
