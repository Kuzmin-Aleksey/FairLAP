package entity

type LapParameter struct {
	Class string `json:"class" db:"class"`
	Value int    `json:"value" db:"value"`
}
