package auth

import (
	"context"
	"errors"
)

var (
	db = map[string]UserData{}
)

func createUser(ctx context.Context, in UserCredential) (user User, err error) {
	user = User{
		UserID:   in.UserID,
		Nickname: in.UserID,
	}
	db[in.UserID] = UserData{
		User:     user,
		Password: in.Password,
	}
	return
}

func isUserIdExisted(ctx context.Context, userID string) (existed bool, err error) {
	_, existed = db[userID]
	if existed {
		err = errors.New("already same user_id is used")
	}
	return
}

func getUserByID(ctx context.Context, userID string) (user User, err error) {
	userData, existed := db[userID]
	if !existed {
		err = errors.New("No User found")
		return
	}

	user = userData.User
	return
}

func updateUserData(ctx context.Context, userID string, update UserUpdate) (user User, err error) {
	currentUser := db[userID]
	db[userID] = UserData{
		User: User{
			UserID:   currentUser.UserID,
			Nickname: update.Nickname,
			Comment:  update.Comment,
		},
		Password: currentUser.Password,
	}
	return
}

func deleteUserData(ctx context.Context, userID string) (err error) {
	delete(db, userID)
	return
}
