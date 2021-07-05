package model

import (
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
)

// 로그인 성공시 정보 추가
func (o *DB) SetAuthInfo(authInfo *AuthInfo) error {
	cKey := genCacheKeyByAuth(authInfo.WalletAuth.WalletAddr)
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	conf := config.GetInstance()
	return o.Cache.Set(cKey, authInfo, time.Duration(conf.Auth.TokenExpiryPeriod*int64(time.Minute)))
}

// 로그인 정보 검색
func (o *DB) GetAuthInfo(walletAddr string) (*AuthInfo, error) {
	cKey := genCacheKeyByAuth(walletAddr)
	authInfo := new(AuthInfo)
	err := o.Cache.Get(cKey, authInfo)
	return authInfo, err
}

func genCacheKeyByAuth(id string) string {
	return config.GetInstance().DBPrefix + ":AUTH:" + id
}
