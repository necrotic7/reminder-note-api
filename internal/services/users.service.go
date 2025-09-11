package services

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/repositories"
	"github.com/zivwu/reminder-note-api/internal/types"
)

type UsersService struct {
	repo *repositories.UsersRepository
}

func NewUsersService(repo *repositories.UsersRepository) *UsersService {
	return &UsersService{
		repo: repo,
	}
}

func (s *UsersService) Login(ctx context.Context, req *types.ReqLoginBody) (result gin.H, err error) {
	userId, err := s.repo.UpsertUser(ctx, req.LineID, req.Name)
	if err != nil {
		log.Println("登入失敗：", err)
		return
	}
	result = gin.H{
		"userId": userId,
	}
	return
}
