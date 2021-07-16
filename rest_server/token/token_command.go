package token

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	ethCtrl "github.com/ONBUFF-IP-TOKEN/baseEthereum/ethcontroller"
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/basenet"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
)

const (
	TokenCmd_CreateNft    uint32 = 0
	TokenCmd_DeleteToken  uint32 = 1
	TokenCmd_OrderProduct uint32 = 2
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
		case TokenCmd_OrderProduct:
			o.OrderProduct(data)
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

func (o *TokenCmd) OrderProduct(data *basenet.CommandData) {
	order := data.Data.(*context.OrderProduct)
	productInfo, err := model.GetDB().GetProductInfo(order.ProductId)
	if err != nil {
		log.Error("OrderProduct GetProductInfo error ", err)
		return
	}
	token := o.itoken.Tokens[Token_nft]
	// 1. receipt 정상 체크
	go func() {
		errCnt := 0
	POLLING:
		tx, isPanding, err := token.eth.GetTransactionByTxHash(order.PurchaseTxHash)
		if err == nil {
			if isPanding {
				log.Debug("is panding : ", isPanding)
				time.Sleep(time.Second * 1)
				errCnt = 0
				goto POLLING
			}
			receipt, err := token.eth.GetTransactionReceipt(tx)
			if err == nil {
				log.Info("GetTransactionReceipt Type:", receipt.Type)
				log.Info("GetTransactionReceipt PostState:", receipt.PostState)
				log.Info("GetTransactionReceipt status :", receipt.Status)
				log.Info("GetTransactionReceipt CumulativeGasUsed:", receipt.CumulativeGasUsed)
				log.Info("GetTransactionReceipt Bloom :", receipt.Bloom)
				for _, logInfo := range receipt.Logs {
					fmt.Printf("GetTransactionReceipt Logs %+v\n", logInfo)
				}

				log.Info("topics 0 : ", receipt.Logs[0].Topics[0].Hex())
				log.Info("topics 1 : ", receipt.Logs[0].Topics[1].Hex())
				log.Info("topics 2 : ", receipt.Logs[0].Topics[2].Hex())

				log.Info("GetTransactionReceipt TxHash:", receipt.TxHash.Hex())
				log.Info("GetTransactionReceipt contractAddress :", receipt.ContractAddress.Hex())
				log.Info("GetTransactionReceipt GasUsed:", receipt.GasUsed)
				log.Info("GetTransactionReceipt blockhash :", receipt.BlockHash.Hex())
				log.Info("GetTransactionReceipt blocknumber :", receipt.BlockNumber)
				log.Info("GetTransactionReceipt TransactionIndex:", receipt.TransactionIndex)

				//token contract address check
				log.Info("token address : ", receipt.Logs[0].Address.Hex())
				if strings.ToUpper(o.conf.TokenAddrs[Token_onit]) != strings.ToUpper(receipt.Logs[0].Address.Hex()) {
					log.Error("Invalid token address :", receipt.Logs[0].Address.Hex())
					return
				}

				//받는 사람 보내는 사람 check
				fromAddr := strings.Replace(receipt.Logs[0].Topics[1].Hex(), "000000000000000000000000", "", -1)
				toAddr := strings.Replace(receipt.Logs[0].Topics[2].Hex(), "000000000000000000000000", "", -1)
				if strings.ToUpper(order.WalletAddr) != strings.ToUpper(fromAddr) {
					log.Error("Invalid from address :", fromAddr)
					return
				}
				if strings.ToUpper(o.conf.ServerWalletAddr) != strings.ToUpper(toAddr) {
					log.Error("Invalid to address :", toAddr)
					return
				}
				// 구입 액수 check
				value := new(big.Int)
				value.SetString(hex.EncodeToString(receipt.Logs[0].Data), 16)
				log.Info("transfer value :", value)

				transferEther := ethCtrl.Convert(value.String(), ethCtrl.Wei, ethCtrl.Ether)
				price := new(big.Rat).SetInt64(int64(productInfo.Price))
				if transferEther.Cmp(price) != 0 {
					log.Error("Invalid purchase price :", transferEther.String())
					return
				}

				// todo 거래 정보가 다 정확하면 nft를 전송하고 nft 테이블에 소유자를 갱신해준다.
			} else if err.Error() == "not found" {
				log.Debug("not found retry GetTransactionReceipt : ", order.PurchaseTxHash)
				time.Sleep(time.Second * 1)
				if errCnt > 3 {
					log.Error("GetTransactionReceipt max try")
					return
				}
				errCnt++
				goto POLLING
			}
		} else {
			log.Debug("GetTransactionByTxHash error : ", err)
			if errCnt > 3 {
				log.Error("GetTransactionByTxHash max try : ", order.PurchaseTxHash)
				return
			}
			errCnt++
			goto POLLING
		}
	}()

	// 2. 영수증 확인 후 nft 전송

	// 3. nft 콜백 스레드에서 db 업데이트 처리

	_ = order
}
