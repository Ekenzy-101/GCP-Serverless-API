package model

import "github.com/Ekenzy-101/GCP-Serverless/types"

type Post struct {
	ID      string  `json:"id,omitempty" firestore:"id,omitempty"`
	Content string  `json:"content,omitempty" firestore:"content,omitempty" validate:"required"`
	Image   string  `json:"image,omitempty" firestore:"image,omitempty"`
	Title   string  `json:"title,omitempty" firestore:"title,omitempty" validate:"required"`
	User    types.M `json:"user,omitempty" firestore:"user,omitempty"`
}
