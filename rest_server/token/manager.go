package token

import (
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
)

type IToken struct {
	conf *config.TokenInfo

	Tokens map[int]*Token

	tokenCmd *TokenCmd
}

func NewTokenManager(conf *config.TokenInfo) *IToken {
	gToken = new(IToken)
	gToken.conf = conf
	gToken.tokenCmd = NewTokenCmd(gToken, conf)
	return gToken
}

func GetToken() *IToken {
	return gToken
}

func (o *IToken) Init() error {
	o.Tokens = map[int]*Token{
		Token_nft:  new(Token),
		Token_onit: new(Token),
	}

	for idx, token := range o.Tokens {
		// callback channel 생성
		token.CreateChannel()

		token.Init(idx, o.conf)

		// mainnet connect
		if err := token.ConnectMainNet(o.conf.MainnetHost); err != nil {
			log.Fatal("ConnectMainNet ", tokenTypes[idx], " error ", err)
		} else {
			log.Info("Mainnet Dial Success ", tokenTypes[idx])
		}

		// subscribe contract
		if err := token.SubcribeContract(o.conf.TokenAddrs[idx]); err != nil {
			log.Fatal("SubcribeContract ", tokenTypes[idx], " error ", err)
		} else {
			log.Info("SubcribeContract Success ", tokenTypes[idx])
		}

		//load contract
		if err := token.LoadContract(o.conf.TokenAddrs[idx]); err != nil {
			log.Fatal("LoadContract ", tokenTypes[idx], " error ", err)
		} else {
			log.Info("LoadContract Success ", tokenTypes[idx])
		}

		//load name, symbol
		if name, symbol, err := token.LoadContractInfo(); err != nil {
			log.Fatal("LoadContractInfo ", tokenTypes[idx], " error ", err)
		} else {
			log.Info("LoadContractInfo ", tokenTypes[idx], " ", name, " ", symbol)
		}
	}

	o.tokenCmd.StartTokenCommand()

	return nil
}
