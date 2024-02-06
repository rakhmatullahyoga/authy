package auth

import (
	"context"
	"errors"
	"net/http"
)

var (
	createUserRepo    = createUser
	deleteUserRepo    = deleteUserData
	getUserByIdRepo   = getUserByID
	userIdExistedRepo = isUserIdExisted
	updateUserRepo    = updateUserData
)

func registerUser(ctx context.Context, in UserCredential) (user User, err error) {
	existed, err := userIdExistedRepo(ctx, in.UserID)
	if err != nil {
		return
	}

	if existed {
		err = ErrUserAlreadyRegistered
		return
	}

	user, err = createUserRepo(ctx, in)
	return
}

func getUser(ctx context.Context, userID string) (user User, err error) {
	user, err = getUserByIdRepo(ctx, userID)
	return
}

func updateUser(ctx context.Context, userID string, params UserUpdate) (user User, statusCode int, err error) {
	selfID := ctx.Value(ClaimsKeyUserID).(string)
	currentUser, err := getUserByIdRepo(ctx, userID)
	if err != nil {
		statusCode = http.StatusNotFound
		return
	}

	if currentUser.UserID != selfID {
		statusCode = http.StatusForbidden
		err = errors.New("No Permission for Update")
		return
	}

	if params.Nickname == "" {
		params.Nickname = currentUser.Nickname
	}

	if params.Comment == "" {
		params.Comment = currentUser.Comment
	}

	user, err = updateUserRepo(ctx, userID, params)
	statusCode = http.StatusOK
	return
}

func deleteUser(ctx context.Context) (err error) {
	selfID := ctx.Value(ClaimsKeyUserID).(string)
	err = deleteUserRepo(ctx, selfID)
	return
}
