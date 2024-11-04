package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	Id         primitive.ObjectID `bson:"_id" json:"order_id"`
	User_id    string             `bson:"user_id" json:"user_id"`
	Book_id    string             `bson:"book_id" json:"book_id"`
	Created_at time.Time          `bson:"created_at" json:"created_at"`
}

type Orders struct {
	ListOfOrders []Order
}

func (O *Orders) AddOrderToList(order Order) {
	O.ListOfOrders = append(O.ListOfOrders, order)
}
