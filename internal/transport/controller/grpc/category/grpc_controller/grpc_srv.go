package grpccontroller

import (
	"context"
	"log"

	"github.com/AlexSamarskii/delivery_service/internal/entity"
	"github.com/AlexSamarskii/delivery_service/internal/transport/controller/grpc/category"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CategoryRepository interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) ([]entity.Category, error)
}

type CategoryGrpcServer struct {
	category.UnimplementedCategoryServiceServer
	repo CategoryRepository
}

func NewCategoryGrpcServer(repo CategoryRepository) *CategoryGrpcServer {
	return &CategoryGrpcServer{repo: repo}
}

func (s *CategoryGrpcServer) GetCategoriesByIds(
	ctx context.Context,
	req *category.GetCategoriesByIdsRequest,
) (*category.GetCategoriesResponse, error) {
	log.Println("[START] GRPC - GetCategoriesByIds")

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if len(req.Ids) == 0 {
		return &category.GetCategoriesResponse{Data: []*category.CategoryDTO{}}, nil
	}

	// Парсим UUID
	uuidIds := make([]uuid.UUID, 0, len(req.Ids))
	for _, idStr := range req.Ids {
		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Printf("Invalid UUID: %s", idStr)
			return nil, status.Errorf(codes.InvalidArgument, "invalid category ID: %s", idStr)
		}
		uuidIds = append(uuidIds, id)
	}

	cats, err := s.repo.FindByIds(ctx, uuidIds)
	if err != nil {
		log.Printf("Failed to fetch categories: %v", err)
		return nil, status.Error(codes.Internal, "failed to fetch categories")
	}

	result := make([]*category.CategoryDTO, len(cats))
	for i, cat := range cats {
		result[i] = &category.CategoryDTO{
			Id:          cat.Id.String(),
			Name:        cat.Name,
			Description: cat.Description,
			Icon:        cat.Icon,
			Status:      cat.Status,
		}
	}

	log.Println("[END] GRPC - GetCategoriesByIds")
	return &category.GetCategoriesResponse{Data: result}, nil
}
