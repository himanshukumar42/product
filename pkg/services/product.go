package services

import (
	"context"
	"net/http"

	"github.com/himanshuk42/product/pkg/db"
	"github.com/himanshuk42/product/pkg/models"
	"github.com/himanshuk42/product/pkg/pb"
)

type Server struct {
	pb.ProductServiceServer
	H db.Handler
}

func (s *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {

	product := &models.Product{
		Name:  req.GetName(),
		Stock: req.GetStock(),
		Price: req.GetPrice(),
	}

	if result := s.H.DB.Create(&product); result.Error != nil {
		return &pb.CreateProductResponse{
			Status: http.StatusConflict,
			Error:  result.Error.Error(),
		}, nil
	}

	return &pb.CreateProductResponse{
		Status: http.StatusCreated,
		Id:     product.Id,
	}, nil
}

func (s *Server) FindOne(ctx context.Context, req *pb.FindOneRequest) (*pb.FindOneResponse, error) {
	var product models.Product

	if result := s.H.DB.First(&product, req.GetId()); result.Error != nil {
		return &pb.FindOneResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	data := &pb.FindOneData{
		Id:    product.Id,
		Name:  product.Name,
		Stock: product.Stock,
		Price: product.Price,
	}

	return &pb.FindOneResponse{
		Status: http.StatusOK,
		Data:   data,
	}, nil
}

func (s *Server) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {
	var product models.Product

	if result := s.H.DB.First(&product, req.GetId()); result.Error != nil {
		return &pb.DecreaseStockResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	if product.Stock <= 0 {
		return &pb.DecreaseStockResponse{
			Status: http.StatusConflict,
			Error:  "stock too low",
		}, nil
	}

	var log models.StockDecreaseLog

	if result := s.H.DB.Where(&models.StockDecreaseLog{OrderId: req.GetOrderId()}).First(&log); result.Error == nil {
		return &pb.DecreaseStockResponse{
			Status: http.StatusConflict,
			Error:  "stock already decreased",
		}, nil
	}

	product.Stock = product.Stock - 1

	s.H.DB.Save(&product)

	log.OrderId = req.GetOrderId()
	log.ProductRefer = product.Id

	s.H.DB.Create(&log)

	return &pb.DecreaseStockResponse{
		Status: http.StatusOK,
	}, nil
}
