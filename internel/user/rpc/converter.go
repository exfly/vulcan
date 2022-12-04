package rpc

import (
	userrepo "github.com/exfly/vulcan/internel/user/repo"

	vulcanv1 "github.com/exfly/vulcan/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user userrepo.User) *vulcanv1.User {
	return &vulcanv1.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
