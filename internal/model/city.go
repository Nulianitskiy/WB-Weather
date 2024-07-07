package model

type City struct {
	Id        int    `db:"id" json:"id"`
	City      string `db:"city" json:"city"`
	Country   string `db:"country" json:"country"`
	Latitude  string `db:"latitude" json:"latitude"`
	Longitude string `db:"longitude" json:"longitude"`
}
