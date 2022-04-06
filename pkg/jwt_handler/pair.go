package jwt_handler

import (
	"strings"
	"time"
)

type TokenPair struct {
	RefreshToken, AccessToken *UserClaim
}

func NewTokenPair(userId int64, userType string) TokenPair {
	return TokenPair{
		RefreshToken: NewUserClaim(userId, userType, time.Now().Add(refreshLifeTime)),
		AccessToken:  NewUserClaim(userId, userType, time.Now().Add(accessLifeTime)),
	}
}
func NewTokenPairFromStrings(refresh, access string) (pair TokenPair, err error) {
	pair.RefreshToken, err = GetClaim(refresh, refreshKey)
	if err != nil {
		return
	}
	pair.AccessToken, err = GetClaim(access, accessKey)
	if err != nil && !strings.HasPrefix(err.Error(), "token is") {
		return
	}
	err = nil
	return
}

func (p TokenPair) GetStrings() (refresh, access string, err error) {
	refresh, err = p.RefreshToken.SetKey(refreshKey)
	if err != nil {
		return "", "", err
	}
	access, err = p.AccessToken.SetKey(accessKey)
	return
}
