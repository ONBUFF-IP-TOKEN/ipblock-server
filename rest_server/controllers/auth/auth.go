package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
	"github.com/dgrijalva/jwt-go"
)

type IAuth struct {
	conf  *config.ApiAuth
	block cipher.Block
}

func NewIAuth(conf *config.ApiAuth) (*IAuth, error) {
	if gAuth == nil {
		gAuth = new(IAuth)
		gAuth.conf = conf
		block, err := aes.NewCipher([]byte(conf.AesKey))
		if err != nil {
			return nil, err
		}
		gAuth.block = block
	}
	return gAuth, nil
}

func GetIAuth() *IAuth {
	return gAuth
}

// func (o *IAuth) IsValidCrypto(origin, vd string) bool {
// 	lvd, err := o.DecryptLoginVD(vd)
// 	if err != nil {
// 		log.Error(err)
// 		return false
// 	}

// 	if origin != lvd.WalletAddr {
// 		return false
// 	}
// 	if lvd.Date >= time.Now().Unix() {
// 		return false
// 	}

// 	if lvd.Date > time.Now().Add(time.Minute*time.Duration(o.conf.SignExpiryPeriod)).Unix() ||
// 		time.Unix(lvd.Date, 0).Add(time.Minute*time.Duration(o.conf.SignExpiryPeriod)).Unix() < time.Now().Unix() {
// 		log.Error("out of date vd : ", lvd.WalletAddr)
// 		return false
// 	}

// 	//기존에 동일한 vd 인지 redis 확인
// 	authInfo, err := model.GetDB().GetAuthInfo(origin)
// 	if err != nil {
// 		// redis에 기존 정보가 없다면 성공
// 		return true
// 	}

// 	if authInfo.Vd == vd { //redis에 기존 vd와 동일한 데이터가 있다면 실패 처리
// 		// todo : web 이 vd처리가 완료되면 활성화 시켜야함
// 		//return false
// 	}

// 	return true
// }

// vd encrypt
// func (o *IAuth) EncryptLoginVD(walletAddr string) (string, error) {
// 	lvd := LoginVD{
// 		WalletAddr: walletAddr,
// 		Date:       time.Now().Unix(),
// 	}
// 	bytrArry, err := json.Marshal(lvd)
// 	if err != nil {
// 		return "", err
// 	}
// 	return o.AesEncrypt(bytrArry), nil
// }

// vd decrypt
// func (o *IAuth) DecryptLoginVD(encryptData string) (*LoginVD, error) {
// 	plainText, err := o.AesDecrypt(encryptData)
// 	if err != nil {
// 		return nil, err
// 	}
// 	lvd := new(LoginVD)
// 	err = json.Unmarshal(plainText, lvd)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return lvd, nil
// }

// auth jwt encrypt
func (o *IAuth) EncryptJwt(walletAddr string) (string, int64, error) {
	var authToken string
	expireDate := time.Now().Add(time.Minute * time.Duration(o.conf.TokenExpiryPeriod)).Unix()

	atClaims := jwt.MapClaims{}
	atClaims["wallet_address"] = walletAddr
	atClaims["exp"] = expireDate

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	authToken, err := at.SignedString([]byte(o.conf.JwtSecretKey))
	if err != nil {
		return "", 0, err
	}
	return authToken, expireDate, nil
}

// auth jwt decrypt
func (o *IAuth) DecryptJwt(jwtStr string) (string, int64, error) {
	atClaims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(jwtStr, atClaims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(o.conf.JwtSecretKey), nil
		})
	if err != nil {
		//exp가 만료되면 여기로 에러 리턴됨
		return "", 0, err
	}
	expireDate, ok := atClaims["exp"].(float64)
	if !ok {
		return "", 0, err
	}
	return fmt.Sprintf("%v", atClaims["wallet_address"]), int64(expireDate), nil
}

// auto token 유효한지 검증
func (o *IAuth) IsValidAuthToken(authToken string) (*string, bool) {
	// todo 기능 구현
	walletAddr, expireDate, err := o.DecryptJwt(authToken)
	if err != nil || len(walletAddr) == 0 {
		return nil, false
	}
	//log.Debug("auth check wallet address:", walletAddr)
	//log.Debug("auth check expiredate:", expireDate)
	if time.Now().Unix() > expireDate {
		log.Info("out of auth token exipre date :", walletAddr)
		return nil, false
	}

	authInfo, err := model.GetDB().GetAuthInfo(walletAddr)
	if err != nil {
		return nil, false
	}
	if authInfo.AuthToken != authToken ||
		authInfo.WalletAuth.WalletAddr != walletAddr ||
		authInfo.ExpireDate != expireDate {
		return nil, false
	}

	return &walletAddr, true
}

func (o *IAuth) AesEncrypt(plaintext []byte) string {
	aesGCM, err := cipher.NewGCM(o.block)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext)
}

func (o *IAuth) AesDecrypt(str string) (plaintext []byte, err error) {
	enc, _ := hex.DecodeString(str)

	aesGCM, err := cipher.NewGCM(o.block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err = aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext, err
}
