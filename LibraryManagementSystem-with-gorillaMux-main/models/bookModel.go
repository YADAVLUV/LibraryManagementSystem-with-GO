package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	Book_name  string             `bson:"book_name" json:"book_name"`
	Author     string             `bson:"author" json:"author"`
	Created_at time.Time          `bson:"created_at" json:"created_at"`
}

type Books struct {
	ListOfBooks []Book
}

func (B *Books) AddBookToList(book Book) {
	B.ListOfBooks = append(B.ListOfBooks, book)
}
