package token

import (
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/basenet"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
)

const (
	TokenCmd_CreateNft   uint32 = 0
	TokenCmd_DeleteToken uint32 = 1
)

type TokenCmd struct {
	itoken  *IToken
	conf    *config.TokenInfo
	command chan *basenet.CommandData
}

func NewTokenCmd(itoken *IToken, conf *config.TokenInfo) *TokenCmd {
	tokenCmd := new(TokenCmd)
	tokenCmd.itoken = itoken
	tokenCmd.conf = conf
	tokenCmd.command = make(chan *basenet.CommandData)
	return tokenCmd
}

func (o *TokenCmd) GetTokenCmdChannel() chan *basenet.CommandData {
	return o.command
}

func (o *TokenCmd) StartTokenCommand() {
	context.GetChanInstance().Put(context.TokenChannel, o.command)

	go func() {
		ticker := time.NewTicker(1 * time.Millisecond)

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
			o.CreateNft(data.Data, data.Callback)
		case TokenCmd_DeleteToken:
			o.DeleteToken(data.Data)
		}
		end := time.Now()

		log.Debug("cmd.kind:", data.CommandType, ",elapsed", end.Sub(start))
	}
	return nil
}

func (o *TokenCmd) CreateNft(data interface{}, cb chan interface{}) {
	product := data.(*context.ProductInfo)

	for i := int64(0); i < product.QuantityTotal; i++ {
		//2-1. nft 생성 요청

		uri := GetNftUri(o.conf.NftUriDomain, product.Id, i+1)

		if txHash, err := o.itoken.Tokens[Token_nft].Nft_CreateERC721(o.conf.ServerWalletAddr, uri); err != nil {
			//resp.SetReturn(resultcode.Result_TokenERC721CreateError)
			log.Error("Nft_CreateERC721 error :", err)
		} else {
			//2-2. db 저장
			if _, err := model.GetDB().InsertProductNFT(product, i+1, context.Product_nft_state_pending, txHash, o.conf.ServerWalletAddr, uri); err != nil {
				//resp.SetReturn(resultcode.Result_DBError)
				log.Error("InsertProductNFT :", err)
			}
		}
	}

	cb <- base.MakeBaseResponse(resultcode.Result_Success)
}

func (o *TokenCmd) DeleteToken(data interface{}) {
	productInfo := data.(*context.ProductInfo)

	_ = productInfo
}
