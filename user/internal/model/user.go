package model

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	Email    string `json:"email"`
}
