package services

import (
	"context"
	"github.com/RTS-1989/go-comment-svc/pkg/db"
	"github.com/RTS-1989/go-comment-svc/pkg/models"
	"github.com/RTS-1989/go-comment-svc/pkg/pb"
	"net/http"
)

type Server struct {
	H db.Handler
}

func (s *Server) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	c := &models.Comment{
		NewsId:   req.NewsId,
		ParentId: req.ParentId,
		Text:     req.Text,
		Censored: req.Censored,
	}

	if result := s.H.DB.Create(c); result.Error != nil {
		return &pb.CreateCommentResponse{
			Status: http.StatusConflict,
			Error:  result.Error.Error(),
		}, nil
	}

	return &pb.CreateCommentResponse{
		Status: http.StatusCreated,
		Id:     uint64(c.ID),
	}, nil
}

func (s *Server) CommentsByNews(ctx context.Context, req *pb.CommentsByNewsRequest) (*pb.CommentsByNewsResponse, error) {
	commentsSlice := make([]*pb.Comment, 0)
	if result := s.H.DB.Where(&models.Comment{NewsId: req.NewsId}).Find(&commentsSlice); result.Error != nil {
		return &pb.CommentsByNewsResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	return &pb.CommentsByNewsResponse{
		Status:   http.StatusOK,
		Comments: commentsSlice,
	}, nil
}
