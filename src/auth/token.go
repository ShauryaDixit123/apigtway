package auth

import (
	"apigtway/src/dtos"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/gofrs/uuid"
)

type Service struct {
	Rcl *redis.Client
}

func (s *Service) CreateToken(userid int64) (*dtos.Token, error) {
	td := &dtos.Token{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	tempUUID, _ := uuid.NewV4()
	td.AccessUUID = tempUUID.String()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	tempUUID, _ = uuid.NewV4()
	td.RefreshUUID = tempUUID.String()

	var err error
	//Creating Access Token
	// os.Setenv("ACCESS_SECRET", os.Getenv("")) //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET_KEY")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	// os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET_KEY")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (s *Service) CreateAuth(userid int64, t *dtos.Token) error {
	at := time.Unix(t.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(t.RtExpires, 0)
	now := time.Now()

	errAccess := s.Rcl.Set(t.AccessUUID, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := s.Rcl.Set(t.RefreshUUID, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (s *Service) FetchAuth(ad dtos.AccessDetails) (uint64, error) {
	userid, er := s.Rcl.Get(ad.AccessUUID).Result()
	if er != nil {
		return 0, er
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

func (s *Service) DeleteAuth(id string) (int64, error) {
	deleted, err := s.Rcl.Del(id).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
