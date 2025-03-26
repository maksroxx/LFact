package fixtures

import (
	"context"
	"fmt"
	"log"

	"github.com/maksroxx/LFact/db"
	"github.com/maksroxx/LFact/types"
)

func AddUser(store *db.Store, fn, ln string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		FirstName: fn,
		LastName:  ln,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	insertedUser, err := store.UserStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}
