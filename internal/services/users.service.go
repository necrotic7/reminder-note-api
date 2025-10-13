package services

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/repositories"
	"github.com/zivwu/reminder-note-api/internal/types"
	"github.com/zivwu/reminder-note-api/internal/utils"
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
	if len(req.LineID) < 1 {
		err = fmt.Errorf("missing line id")
		return
	}

	user, err := s.repo.UpsertUser(ctx, req.LineID, req.Name)
	if err != nil {
		log.Println("登入失敗：", err)
		return
	}
	// 產生token
	token, err := utils.GenToken(utils.TokenClaims{
		UserID: user.ID.Hex(),
		Name:   user.Name,
	})
	if err != nil {
		log.Println("產生token失敗：", err)
		return
	}
	result = gin.H{
		"token": token,
	}
	return
}
