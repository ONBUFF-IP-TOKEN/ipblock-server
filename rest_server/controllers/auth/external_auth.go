package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
)

func CheckAuthToken(walletAddr, authToken string) (bool, error) {
	conf := config.GetInstance()

	params := context.NewVerifyAuthToken()
	params.AuthToken = authToken
	params.WalletAddr = walletAddr

	callURL := fmt.Sprintf("%s%s", conf.Auth.ApiAuthDomain, conf.Auth.ApiAuthVerify)
	buff := bytes.NewBuffer(nil)
	pbytes, _ := json.Marshal(params)
	buff = bytes.NewBuffer(pbytes)

	req, err := http.NewRequest("POST", callURL, buff)
	if err != nil {
		log.Error(err)
		return false, err
	}

	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		log.Errorf("membership resp: %v, err: %v", resp, err)
		return false, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Errorf("resp: %v, body close err: %v", resp, err)
		}
	}()

	decoder := json.NewDecoder(resp.Body)
	baseResp := new(base.BaseResponse)
	err = decoder.Decode(baseResp)
	if err != nil {
		log.Errorf("resp: %v, docode err: %v", resp, err)
		return false, err
	}

	if baseResp.Message != "success" {
		err := errors.New(baseResp.Message)
		log.Errorf("resp: %v, body close err: %v", resp, err)
		return false, err
	}

	return true, nil
}
