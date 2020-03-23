package store

import "time"

type Address struct {
	ID        string    `json:"id" bson:"_id"`
	UserID    int32     `json:"user_id" bson:"user_id"`
	Name      string    `json:"name" bson:"name,omitempty"`
	Phone     string    `json:"phone" bson:"phone,omitempty"`
	Area      string    `json:"area" bson:"area,omitempty"`
	Address   string    `json:"address" bson:"address,omitempty"`
	CreatAt   time.Time `json:"created" bson:"CreatAt,omitempty"`
	UpdateAt  time.Time `json:"updated" bson:"UpdateAt,omitempty"`
	IsDefault bool      `json:"isDefault" bson:"isDefault,omitempty"`
}
