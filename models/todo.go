package models

type Todo struct {
	Id       int    `json:"id"`
	Text     string `json:"text"`
	Priority int    `json:"priority"`
}

type User struct {
	Email string
	Mdp   string
}
