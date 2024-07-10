package model

type City struct {
	Id        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Country   string `db:"country" json:"country"`
	Latitude  string `db:"latitude" json:"latitude"`
	Longitude string `db:"longitude" json:"longitude"`
}

type CityList struct {
	Cities []string `json:"cities"`
}
