package model

type User struct {
	Id           int     `json:"id"`
	Username     string  `json:"username"`
	Password     string  `json:"password"`
	Email        string  `json:"email"`
	PhoneNumber  string  `json:"phone_number"`
	PhotoProfile string  `json:"photo_profile"`
	Balance      float64 `json:"balance"`
}
