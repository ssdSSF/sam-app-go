package model

type Student struct {
	Id        int64 `json:"Id,string"`
	FirstName string
	LastName  string
	Email     string
}
