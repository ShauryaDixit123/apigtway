package users

import (
	"apigtway/src/auth"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GenAuth(ctx *gin.Context) {
	id := ctx.Query("id")
	fmt.Println(id)
	userid, er := strconv.Atoi(id)
	if er != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid user id",
		})
		return
	}
	s := auth.Service{
		Rcl: h.rcl,
	}
	t, er := s.CreateToken(int64(userid))
	if er != nil {
		ctx.JSON(400, gin.H{
			"message": "Unable to generate token",
		})
		return
	}
	// creating auth
	if er = s.CreateAuth(int64(userid), t); er != nil {
		ctx.JSON(400, gin.H{
			"message": "Unable to create auth",
		})
		return
	}
	ctx.JSON(200, (gin.H{
		"access_token":  t.AccessToken,
		"refresh_token": t.RefreshToken,
	}))
}

func (h *Handler) SetToken(ctx *gin.Context) {

}
