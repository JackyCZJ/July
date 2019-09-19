package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackyczj/NoGhost/pkg/auth"
)

type UserInformation struct {
	Username string `json:"username" validate:"min=1,max=32"`
	Password string `json:"password,omitempty" validate:"min=1,max=32"`
	UserId   string `json:"user_id"`
	Email    string `json:"email"`
	Role     int    `json:"role"`
	Gander   int    `json:"gander"`
	Phone    string `json:"phone"`
	sync.Mutex
}

func (u *UserInformation) Create() error {
	u.Lock()
	defer u.Unlock()
	var err error
	u.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		return err
	}

	id, err := Client.db.Collection("user").InsertOne(context.TODO(), u)
	if err != nil {
		return err
	}
	fmt.Println("New member register , id:", id.InsertedID)
	return nil
}

func (u *UserInformation) GetUser() (*UserInformation, error) {
	u.Lock()
	defer u.Unlock()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	result := Client.db.Collection("user").FindOne(ctx, u)
	if result.Err() != nil {
		return nil, result.Err()
	}
	err := result.Decode(&u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
