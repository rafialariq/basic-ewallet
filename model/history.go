package model

type History struct {
	Phone    string  `json:"phone_number"`
	MoreThan float64 `json:"more_than"`
	LessThan float64 `json:"less_than"`
}
