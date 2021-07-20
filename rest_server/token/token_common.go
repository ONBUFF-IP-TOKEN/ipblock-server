package token

import (
	"github.com/ONBUFF-IP-TOKEN/baseEthereum/ethcontroller"
	"github.com/ONBUFF-IP-TOKEN/basenet"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
)

type Token struct {
	TokenType int

	eth         *ethcontroller.EthClient
	tokenName   string
	tokenSymbol string
	resCh       chan *basenet.CommandData

	conf *config.TokenInfo
}

func (o *Token) Init(tokenType int, conf *config.TokenInfo) {
	o.TokenType = tokenType
	o.conf = conf
	o.eth = ethcontroller.NewEthClient(o.resCh)
}

func (o *Token) ConnectMainNet(host string) error {
	if err := o.eth.GetDial(host); err != nil {
		return err
	}
	return nil
}

func (o *Token) SubcribeContract(tokenAddr string) error {
	if o.TokenType == Token_nft {
		if err := o.eth.SubcribeContract(tokenAddr); err != nil {
			return err
		}
	} else if o.TokenType == Token_onit {
		if err := o.eth.Onit_SubcribeContract(tokenAddr); err != nil {
			return err
		}
	}

	return nil
}

func (o *Token) LoadContract(tokenAddr string) error {
	if o.TokenType == Token_nft {
		if err := o.Nft_LoadContract(tokenAddr); err != nil {
			return err
		}
	} else if o.TokenType == Token_onit {
		if err := o.Onit_LoadContract(tokenAddr); err != nil {
			return err
		}
	}

	return nil
}

func (o *Token) LoadContractInfo() (string, string, error) {
	var err error
	if o.TokenType == Token_nft {
		o.Nft_LoadContractInfo()
	} else if o.TokenType == Token_onit {
		o.Onit_LoadContractInfo()
	}

	return o.tokenName, o.tokenSymbol, err
}

func (o *Token) CreateChannel() {
	o.resCh = make(chan *basenet.CommandData)

	go func() {
		defer close(o.resCh)
		for {
			cmd := <-o.resCh
			log.Debug("callback type : ", cmd.CommandType)
			log.Debug("callback data : ", cmd.Data)
			o.CallBackCmdProc(cmd)
		}
	}()
}

func (o *Token) CallBackCmdProc(cmd *basenet.CommandData) {
	cmdType := cmd.CommandType
	switch cmdType {
	case ethcontroller.Ch_type_transfer:
		transInfo := cmd.Data.(ethcontroller.CallBack_Transfer)
		if transInfo.FromAddr == gNullAddress && transInfo.ToAddr != gNullAddress {
			// 최초 생성 처리
			// devlab
			if err := model.GetDB().UpdateTokenID(transInfo.TxHash, transInfo.TokenID); err == nil {
				model.GetDB().InsertHistory(transInfo.TxHash, transInfo.FromAddr, transInfo.ToAddr, transInfo.TokenID, token_state_mint)
			}
			//product
			if _, err := model.GetDB().UpdateProductNftTokenID(transInfo.TxHash, transInfo.TokenID, context.Nft_state_mint); err == nil {
				model.GetDB().InsertHistory(transInfo.TxHash, transInfo.FromAddr, transInfo.ToAddr, transInfo.TokenID, token_state_mint)
			}
		} else if transInfo.FromAddr != gNullAddress && transInfo.ToAddr != gNullAddress {
			// 코인 전송 처리
			// devlab
			if err := model.GetDB().UpdateTransfer(transInfo.TxHash, transInfo.FromAddr, transInfo.ToAddr, transInfo.TokenID); err == nil {
				model.GetDB().InsertHistory(transInfo.TxHash, transInfo.FromAddr, transInfo.ToAddr, transInfo.TokenID, token_state_transfer)
			}
			// product
			if _, err := model.GetDB().UpdateProductNftOwner(transInfo.FromAddr, transInfo.ToAddr, transInfo.TokenID); err == nil {
				model.GetDB().UpdateOrderState(transInfo.TokenID, context.Order_state_nft_transfer_complete)
				model.GetDB().InsertHistory(transInfo.TxHash, transInfo.FromAddr, transInfo.ToAddr, transInfo.TokenID, token_state_transfer)
			}
		} else if transInfo.FromAddr != gNullAddress && transInfo.ToAddr == gNullAddress {
			// 코인 삭체 처리 : 히스토리에 먼저 남기고 item 테이블 삭제 한다.
			insertId, err := model.GetDB().InsertHistory(transInfo.TxHash, transInfo.FromAddr, transInfo.ToAddr, transInfo.TokenID, token_state_burn)
			if err == nil && insertId != 0 {
				model.GetDB().DeleteItemByTokenId(transInfo.TokenID)
			}
		}
	}
}

// func (o *Token) CallBackCmdProc(cmd *basenet.CommandData) {
// 	cmdType := cmd.CommandType
// 	switch cmdType {
// 	case ethcontroller.Ch_type_transfer:
// 		transInfo := cmd.Data.(ethcontroller.CallBack_Transfer)
// 		if transInfo.FromAddr == gNullAddress && transInfo.ToAddr != gNullAddress {
// 			// 최초 생성 처리
// 		} else if transInfo.FromAddr != gNullAddress && transInfo.ToAddr != gNullAddress {
// 			// 코인 전송 처리
// 		} else if transInfo.FromAddr != gNullAddress && transInfo.ToAddr == gNullAddress {
// 			// 코인 삭체 처리 : 히스토리에 먼저 남기고 item 테이블 삭제 한다.
// 		}
// 	}
// }
