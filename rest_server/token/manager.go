package token

import (
	"github.com/ONBUFF-IP-TOKEN/baseEthereum/ethcontroller"
	"github.com/ONBUFF-IP-TOKEN/basenet"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
)

type IToken struct {
	ethClient   *ethcontroller.EthClient
	conf        *config.TokenInfo
	tokenName   string
	tokenSymbol string

	asyncResponse chan *basenet.CommandData
}

func NewTokenManager(conf *config.TokenInfo) *IToken {
	gToken = new(IToken)
	gToken.conf = conf
	return gToken
}

func GetToken() *IToken {
	return gToken
}

func (o *IToken) Init() error {
	o.CreateChannel()
	o.ethClient = ethcontroller.NewEthClient(o.asyncResponse)

	// mainnet 연결
	if err := o.ConnectMainNet(); err != nil {
		return err
	}

	// subscribe contract
	if err := o.SubcribeContract(); err != nil {
		return err
	}

	// load contract
	if err := o.LoadContract(o.conf.NFTTokenAddr); err != nil {
		return err
	}

	// callback channel 생성

	return nil
	//test code
	//uri := "http://www.naver.com/" + strconv.FormatInt(datetime.GetTS2MilliSec(), 10)
	//uri := "http://www.naver.com/ip1"
	//o.CreateERC721Token(o.conf.ServerWalletAddr, o.conf.ServerWalletAddr, uri, o.conf.ServerPrivateKey)
	//o.Approve(o.conf.ServerWalletAddr, o.conf.ServerPrivateKey, 113)
	txhash, err := o.TransferERC721("0xfc788F6956E98feb367b04f442F7CF8C771c25E9", "0x38f998d033990a315b08AFc0F78059Fb7D11Dc4d", 2)
	if err != nil {
		log.Debug("err ", err)
	}
	log.Debug("txhash ", txhash)

	return nil
}

func (o *IToken) ConnectMainNet() error {
	if err := o.ethClient.GetDial(o.conf.MainnetHost); err != nil {
		log.Error("Mainnet Dial error: ", err)
		return err
	}

	log.Info("Mainnet Dial Success")
	return nil
}

func (o *IToken) SubcribeContract() error {
	if err := o.ethClient.SubcribeContract(o.conf.NFTTokenAddr); err != nil {
		log.Error("StartSubscribeBlock error:", err)
		return err
	}
	log.Info("SubcribeContract Success")
	return nil
}

func (o *IToken) CreateERC721(wallertAddr, uri string) (string, error) {
	return o.CreateERC721Token(o.conf.ServerWalletAddr, wallertAddr, uri, o.conf.ServerPrivateKey)
}

func (o *IToken) TransferERC721(fromAddr, toAddr string, tokenId int64) (string, error) {
	return o.TransferERC721Token(o.conf.ServerWalletAddr, fromAddr, toAddr, o.conf.ServerPrivateKey, tokenId)
}

func (o *IToken) Burn(tokenId int64) (string, error) {
	return o.BurnToken(o.conf.ServerWalletAddr, o.conf.ServerPrivateKey, tokenId)
}

func (o *IToken) CreateChannel() {
	o.asyncResponse = make(chan *basenet.CommandData)

	go func() {
		defer close(o.asyncResponse)
		for {
			cmd := <-o.asyncResponse
			log.Debug("callback type : ", cmd.CommandType)
			log.Debug("callback data : ", cmd.Data)
			o.CallBackCmdProc(cmd)
		}
	}()
}

func (o *IToken) CallBackCmdProc(cmd *basenet.CommandData) {
	cmdType := cmd.CommandType
	switch cmdType {
	case ethcontroller.Ch_type_transfer:
		transInfo := cmd.Data.(ethcontroller.CallBack_Transfer)
		if transInfo.FromAddr == gNullAddress && transInfo.ToAddr != gNullAddress {
			// 최초 생성 처리
			model.GetDB().UpdateTokenID(transInfo.TxHash, transInfo.TokenID)
			model.GetDB().InsertHistory(transInfo.TxHash, transInfo.FromAddr, transInfo.ToAddr, transInfo.TokenID, token_state_mint)
		} else if transInfo.FromAddr != gNullAddress && transInfo.ToAddr != gNullAddress {
			// 코인 전송 처리
			model.GetDB().UpdateTransfer(transInfo.TxHash, transInfo.FromAddr, transInfo.ToAddr, transInfo.TokenID)
			model.GetDB().InsertHistory(transInfo.TxHash, transInfo.FromAddr, transInfo.ToAddr, transInfo.TokenID, token_state_transfer)
		} else if transInfo.FromAddr != gNullAddress && transInfo.ToAddr == gNullAddress {
			// 코인 삭체 처리 : 히스토리에 먼저 남기고 item 테이블 삭제 한다.
			insertId, err := model.GetDB().InsertHistory(transInfo.TxHash, transInfo.FromAddr, transInfo.ToAddr, transInfo.TokenID, token_state_burn)
			if err == nil && insertId != 0 {
				model.GetDB().DeleteItemByTokenId(transInfo.TokenID)
			}

		}
	}
}
