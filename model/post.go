package model

type Post struct {
	ID    string `json:"id,omitempty"`
	Image string `json:"image,omitempty"`
	Title string `json:"title,omitempty"`
	User  *User  `json:"user,omitempty"`
}
