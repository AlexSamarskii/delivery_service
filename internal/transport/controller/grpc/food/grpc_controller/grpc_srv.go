package grpccontroller

import (
	"context"
	"log"

	"github.com/AlexSamarskii/delivery_service/internal/entity"
	"github.com/AlexSamarskii/delivery_service/internal/entity/dto"
	"github.com/AlexSamarskii/delivery_service/internal/transport/controller/grpc/food"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FoodRepository interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) ([]dto.FoodInfoDto, error)
}

type UpdateService interface {
	Execute(ctx context.Context, req entity.FoodUpdateReq) error
}

type FoodGrpcServer struct {
	food.UnimplementedFoodServiceServer
	repo          FoodRepository
	updateService UpdateService
}

func NewFoodGrpcServer(repo FoodRepository, updateService UpdateService) *FoodGrpcServer {
	return &FoodGrpcServer{
		repo:          repo,
		updateService: updateService,
	}
}

// GetFoodsByIds — получает список блюд по ID
func (f *FoodGrpcServer) GetFoodsByIds(
	ctx context.Context,
	req *food.GetFoodsByIdsRequest,
) (*food.GetFoodsResponse, error) {
	log.Println("[START] GRPC - GetFoodsByIds")

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if len(req.FoodIds) == 0 {
		return &food.GetFoodsResponse{Data: []*food.FoodDTO{}}, nil
	}

	// Парсим UUID
	uuidIds := make([]uuid.UUID, 0, len(req.FoodIds))
	for _, idStr := range req.FoodIds {
		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Printf("Invalid UUID: %s", idStr)
			return nil, status.Errorf(codes.InvalidArgument, "invalid food ID: %s", idStr)
		}
		uuidIds = append(uuidIds, id)
	}

	foods, err := f.repo.FindByIds(ctx, uuidIds)
	if err != nil {
		log.Printf("Failed to fetch foods: %v", err)
		return nil, status.Error(codes.Internal, "failed to fetch foods")
	}

	result := make([]*food.FoodDTO, len(foods))
	for i, foodItem := range foods {
		result[i] = &food.FoodDTO{
			Id:           foodItem.Id.String(),
			Name:         foodItem.Name,
			Description:  foodItem.Description,
			Image:        foodItem.Images,
			Price:        float32(foodItem.Price),
			AvgPoint:     float32(foodItem.AvgPoint),
			CommentQty:   int64(foodItem.CommentQty),
			CategoryId:   foodItem.CategoryId.String(),
			RestaurantId: foodItem.RestaurantId.String(),
			Status:       foodItem.Status,
		}
	}

	log.Println("[END] GRPC - GetFoodsByIds")
	return &food.GetFoodsResponse{Data: result}, nil
}

func (f *FoodGrpcServer) UpdateFood(
	ctx context.Context,
	req *food.UpdateFoodRequest,
) (*food.UpdateFoodResponse, error) {
	log.Println("[START] GRPC - UpdateFood")

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	foodId, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Invalid food ID: %s", req.Id)
		return nil, status.Errorf(codes.InvalidArgument, "invalid food ID: %s", req.Id)
	}

	updateReq := entity.FoodUpdateReq{
		Id:          foodId,
		Name:        req.Name,
		Description: req.Description,
	}

	// Опциональные поля — только если заданы
	if req.Status != "" {
		updateReq.Status = &req.Status
	}
	if req.RestaurantId != "" {
		updateReq.RestaurantId = &req.RestaurantId
	}
	if req.CategoryId != "" {
		updateReq.CategoryId = &req.CategoryId
	}
	if req.Image != "" {
		updateReq.Image = &req.Image
	}

	if err := f.updateService.Execute(ctx, updateReq); err != nil {
		log.Printf("UpdateService.Execute failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to update food")
	}

	log.Println("[END] GRPC - UpdateFood")
	return &food.UpdateFoodResponse{FoodId: req.Id}, nil
}
