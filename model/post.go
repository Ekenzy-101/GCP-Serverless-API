package model

import (
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Ekenzy-101/GCP-Serverless/types"
)

type Post struct {
	ID        string    `json:"id,omitempty" firestore:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty" firestore:"createdAt,omitempty"`
	Content   string    `json:"content,omitempty" firestore:"content,omitempty" validate:"required"`
	Image     string    `json:"image,omitempty" firestore:"image,omitempty"`
	Title     string    `json:"title,omitempty" firestore:"title,omitempty" validate:"required"`
	User      types.M   `json:"user,omitempty" firestore:"user,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" firestore:"updatedAt,omitempty"`
}

func NewPostFromDocument(document *firestore.DocumentSnapshot) (*Post, error) {
	post := &Post{}
	err := document.DataTo(post)
	return post, err
}

func (p *Post) SetContent(value string) *Post {
	p.Content = value
	return p
}

func (p *Post) SetID(value string) *Post {
	p.ID = value
	return p
}

func (p *Post) SetImage(value string) *Post {
	p.Image = value
	return p
}

func (p *Post) SetTimestamps(createTime, updateTime time.Time) *Post {
	p.CreatedAt = createTime
	p.UpdatedAt = updateTime
	return p
}

func (p *Post) SetTitle(value string) *Post {
	p.Title = value
	return p
}

func (p *Post) SetUser(value types.M) *Post {
	p.User = value
	return p
}
