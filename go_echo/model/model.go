package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Book struct {
	ID          bson.ObjectId `bson:"_id,omitempty"` // omitempty คือเป็นค่าว่างได้
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	Author      string        `bson:"author"`
	Publisher   string        `bson:"publisher"`
	PublishDate time.Time
}
