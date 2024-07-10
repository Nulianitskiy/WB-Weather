package model

type Weather struct {
	Id          int     `json:"weather_id"`
	IdCity      int     `json:"city_id"`
	Temp        float64 `json:"temp"`
	Date        string  `json:"date"`
	WeatherData []byte  `json:"weather_data"`
}

type ResponseShortWeatherByCity struct {
	City    string   `db:"name" json:"city"`
	Country string   `db:"country" json:"country"`
	Temp    string   `db:"temp" json:"temp"`
	Date    []string `db:"weather_dates" json:"date"`
}

type ResponseFullWeather struct {
	WeatherData map[string]interface{} `db:"weather_data" json:"weather"`
}

type WeatherAPIResponse struct {
	List []struct {
		Dt   int64 `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			SeaLevel  int     `json:"sea_level"`
			GrndLevel int     `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			TempKf    float64 `json:"temp_kf"`
		} `json:"main"`
		Weather []struct {
			Id          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
			Gust  float64 `json:"gust"`
		} `json:"wind"`
		Visibility int     `json:"visibility"`
		Pop        float64 `json:"pop"`
		Sys        struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
}
