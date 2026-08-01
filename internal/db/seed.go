package db

import (
	"context"

	"github.com/dneedsleep/Social/internal/store"
)

var usersArr = []struct {
	FirstName string
	LastName  string
	Email     string
}{
	{"John", "Smith", "john.smith@example.com"},
	{"Emma", "Johnson", "emma.johnson@example.com"},
	{"Michael", "Brown", "michael.brown@example.com"},
	{"Olivia", "Davis", "olivia.davis@example.com"},
	{"William", "Wilson", "william.wilson@example.com"},
}

func Seed(store store.Storage) error {
	ctx := context.Background()

	users := generateUsers(100)

	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			return err
		}
	}

	return nil
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		data := usersArr[i%len(usersArr)]

		users[i] = &store.User{
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Email:     data.Email,
			Password:  "password123",
		}
	}

	return users
}
