package model

type City struct {
	Id        int     `db:"id" json:"id"`
	City      string  `db:"city" json:"city"`
	Country   string  `db:"country" json:"country"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	Longitude float64 `db:"longitude" json:"longitude"`
}
