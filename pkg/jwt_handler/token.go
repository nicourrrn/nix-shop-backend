package jwt_handler

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	accessKey, refreshKey           string
	accessLifeTime, refreshLifeTime time.Duration
)

func init() {
	accessKey = os.Getenv("accessKey")
	refreshKey = os.Getenv("refreshKey")
	lifeTime, err := strconv.Atoi(os.Getenv("accessLifeTime"))
	if err != nil {
		log.Fatalln(err)
	}
	accessLifeTime = time.Duration(lifeTime) * time.Second
	lifeTime, err = strconv.Atoi(os.Getenv("refreshLifeTime"))
	if err != nil {
		log.Fatalln(err)
	}
	refreshLifeTime = time.Duration(lifeTime) * time.Second

}

func GetAccess() string {
	return accessKey
}

func GetRefresh() string {
	return refreshKey
}

type UserClaim struct {
	jwt.StandardClaims
	UserId   int64
	UserType string
}

func NewUserClaim(userId int64, userType string, lifeTime time.Time) *UserClaim {
	return &UserClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: lifeTime.Unix(),
		},
		UserId:   userId,
		UserType: userType,
	}
}

func (c *UserClaim) SetKey(key string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(key))
}

func GetClaim(token, key string) (*UserClaim, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &UserClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if jwtToken == nil {
		return nil, err
	}
	claims, ok := jwtToken.Claims.(*UserClaim)
	if !ok {
		return nil, errors.New("failed to parse")
	}
	return claims, err
}
