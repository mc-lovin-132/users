package handlers

import (
	"errors"
	"net/mail"

	"github.com/mc-lovin-132/users/internal/domain"
	"github.com/mc-lovin-132/users/pb"

	passwordvalidator "github.com/ginbun/go-password-validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func errorMapper(err error) error {
	if errors.Is(err, domain.ErrNotUniqueEmail) {
		return status.Error(codes.AlreadyExists, err.Error())
	} else if errors.Is(err, domain.ErrEmptyName) ||
		errors.Is(err, domain.ErrInvalidEmail) ||
		errors.Is(err, domain.ErrInvalidPassword) ||
		errors.Is(err, domain.ErrNotEnoughArgs) ||
		errors.Is(err, domain.ErrToMuchArgs) {
		return status.Error(codes.InvalidArgument, err.Error())
	} else if errors.Is(err, domain.ErrUserNotFound) {
		return status.Error(codes.NotFound, err.Error())
	} else if errors.Is(err, domain.ErrInternal) {
		return status.Error(codes.Internal, err.Error())
	}
	return status.Error(codes.Internal, err.Error())
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidByEntropy(password string) bool {
	const minEntropyBits = 60
	err := passwordvalidator.Validate(password, minEntropyBits)
	return err == nil
}

func createRequestToDomain(in *pb.CreateRequest) (*domain.User, error) {
	if in.Name == "" {
		return nil, domain.ErrEmptyName
	} else if !isEmailValid(in.Email) {
		return nil, domain.ErrInvalidEmail
	} else if !isValidByEntropy(in.Password) {
		return nil, domain.ErrInvalidPassword
	}
	return &domain.User{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	}, nil
}

func domainToProto(user *domain.User) *pb.User {
	return &pb.User{
		Id:       int64(user.ID),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
}

func domainToGetResponse(user *domain.User) *pb.GetResponse {
	return &pb.GetResponse{
		User: domainToProto(user),
	}
}
