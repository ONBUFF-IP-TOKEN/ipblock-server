package token

import (
	"time"

	"github.com/ONBUFF-IP-TOKEN/basenet"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
)

const (
	TokenCmd_CreateNft   uint32 = 0
	TokenCmd_DeleteToken uint32 = 1
)

type TokenCmd struct {
	itoken  *IToken
	command chan *basenet.CommandData
}

func NewTokenCmd(itoken *IToken) *TokenCmd {
	tokenCmd := new(TokenCmd)
	tokenCmd.itoken = itoken
	return tokenCmd
}

func (o *TokenCmd) GetTokenCmdChannel() chan *basenet.CommandData {
	return o.command
}

func (o *TokenCmd) StartTokenCommand() {
	context.GetInstance().Put(context.TokenChannel, o.command)

	go func() {
		ticker := time.NewTicker(1 * time.Second)

		defer func() {
			ticker.Stop()
		}()

		for {
			select {
			case ch := <-o.command:
				o.CommandProc(ch)
			case <-ticker.C:
			}
		}
	}()
}

func (o *TokenCmd) CommandProc(data *basenet.CommandData) error {

	if data.Data != nil {
		start := time.Now()
		switch data.CommandType {
		case TokenCmd_CreateNft:
			o.CreateNft(data.Data)
		case TokenCmd_DeleteToken:
			o.DeleteToken(data.Data)
		}
		end := time.Now()

		log.Debug("cmd.kind:", data.CommandType, ",elapsed", end.Sub(start))
	}
	return nil
}

func (o *TokenCmd) CreateNft(data interface{}) {
	productInfo := data.(*context.ProductInfo)

	_ = productInfo
}

func (o *TokenCmd) DeleteToken(data interface{}) {
	productInfo := data.(*context.ProductInfo)

	_ = productInfo
}
