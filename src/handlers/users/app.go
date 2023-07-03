package users

import (
	"apigtway/src/auth"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
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

func (h *Handler) CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		typetkn := "access"
		s := auth.Service{
			Rcl: h.rcl,
		}
		if er := auth.TokenValid(ctx.Request, typetkn); er != nil {
			ctx.JSON(401, gin.H{"message": "INVALID TOKEN"})
			return
		}
		ad, er := auth.ExtractMetaData(ctx.Request, typetkn)
		if er != nil {
			ctx.JSON(401, gin.H{"message": "COULDNT EXTRACT METADATA"})
			return
		}
		if ad != nil {
			id, er := s.FetchAuth(*ad)
			if er != nil {
				ctx.JSON(401, gin.H{"message": "COULDNT FIND TOKEN"})
				return
			}
			ctx.Request.Header.Set("user_id", fmt.Sprint(id, "-OK"))
			ctx.Next()
		}
	}
}

func (h *Handler) RefreshAuth(ctx *gin.Context) {
	tkn := map[string]string{}
	if er := ctx.ShouldBindJSON(&tkn); er != nil {
		ctx.JSON(400, er)
		return
	}
	rft := tkn["refresh_token"]
	rfStr := "REFRESH_SECRET_KEY"
	s := auth.Service{
		Rcl: h.rcl,
	}
	ptkn, er := auth.ParseToken(rft, rfStr)
	if er != nil {
		ctx.JSON(401, er)
		return
	}
	if _, ok := ptkn.Claims.(jwt.Claims); !ok && !ptkn.Valid {
		ctx.JSON(401, er)
		return
	}
	claims, ok := ptkn.Claims.(jwt.MapClaims)
	if ok && ptkn.Valid {
		rfuuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			ctx.JSON(401, er)
			return
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, "Error occurred")
			return
		}
		del, er := s.DeleteAuth(rfuuid)
		if er != nil || del == 0 {
			ctx.JSON(401, er)
			return
		}
		newTkn, er := s.CreateToken(int64(userId))
		if er != nil {
			ctx.JSON(401, er)
			return
		}
		if er = s.CreateAuth(int64(userId), newTkn); er != nil {
			ctx.JSON(401, er)
			return
		}
		ctx.JSON(200, (gin.H{
			"access_token":  newTkn.AccessToken,
			"refresh_token": newTkn.RefreshToken,
		}))
		return
	}
}

func (h *Handler) PingWithAuth(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Pong!",
	})
}
