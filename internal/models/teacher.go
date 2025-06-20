package models

type Teacher struct {
	ID        int    `json:"id,omitempty"`
	Firstname string `json:"first_name,omitempty"`
	Lastname  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Class     string `json:"class,omitempty"`
	Subject   string `json:"subject,omitempty"`
}
