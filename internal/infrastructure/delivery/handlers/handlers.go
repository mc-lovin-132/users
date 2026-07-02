package handlers

import (
	"context"

	"github.com/mc-lovin-132/users/internal/domain"
	"github.com/mc-lovin-132/users/pb"
)

type service interface {
	Create(ctx context.Context, data *domain.User) (int, error)
	Delete(ctx context.Context, id int) error
	Get(ctx context.Context, id *int, email *string) (*domain.User, error)
	Update(ctx context.Context, id int, name, email, password *string) (int, error)
}

type Handler struct {
	pb.UnimplementedUserServiceServer
	service service
}

func New(service service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	user, err := createRequestToDomain(in)
	if err != nil {
		return nil, errorMapper(err)
	}
	id, err := h.service.Create(ctx, user)
	if err != nil {
		return nil, errorMapper(err)
	}
	return &pb.CreateResponse{Id: int64(id)}, nil
}
func (h *Handler) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	if in.Email != nil {
		if !isEmailValid(*in.Email) {
			return nil, errorMapper(domain.ErrInvalidEmail)
		}
	}
	if in.Password != nil {
		if !isValidByEntropy(*in.Password) {
			return nil, errorMapper(domain.ErrInvalidPassword)
		}
	}
	id, err := h.service.Update(ctx, int(in.Id), in.Name, in.Email, in.Password)
	if err != nil {
		return nil, errorMapper(err)
	}
	return &pb.UpdateResponse{Id: int64(id)}, nil
}

func (h *Handler) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	var id *int
	var email *string
	switch selector := in.Selector.(type) {
	case *pb.GetRequest_Email:
		email = &selector.Email
	case *pb.GetRequest_Id:
		idValue := int(selector.Id)
		id = &idValue
	default:
		return nil, errorMapper(domain.ErrNotEnoughArgs)
	}
	user, err := h.service.Get(ctx, id, email)
	if err != nil {
		return nil, errorMapper(err)
	}
	return domainToGetResponse(user), nil
}
func (h *Handler) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := h.service.Delete(ctx, int(in.Id))
	if err != nil {
		return nil, errorMapper(err)
	}
	return &pb.DeleteResponse{}, nil
}
