package token

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	ethCtrl "github.com/ONBUFF-IP-TOKEN/baseEthereum/ethcontroller"
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/basenet"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/cdn/azure"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
)

const (
	TokenCmd_CreateNft                uint32 = 0
	TokenCmd_DeleteToken              uint32 = 1
	TokenCmd_OrderProduct             uint32 = 2
	TokenCmd_CreatNftByAut            uint32 = 3
	TokenCmd_Bid_Deposit_CheckReceipt uint32 = 4
	TokenCmd_Bid_Winner_CheckReceipt  uint32 = 5
)

type TokenCmd struct {
	itoken  *IToken
	conf    *config.TokenInfo
	confCdn *config.Cdn
	command chan *basenet.CommandData
}

func NewTokenCmd(itoken *IToken, conf *config.TokenInfo, confCdn *config.Cdn) *TokenCmd {
	tokenCmd := new(TokenCmd)
	tokenCmd.itoken = itoken
	tokenCmd.conf = conf
	tokenCmd.confCdn = confCdn
	tokenCmd.command = make(chan *basenet.CommandData)
	return tokenCmd
}

func (o *TokenCmd) GetTokenCmdChannel() chan *basenet.CommandData {
	return o.command
}

func (o *TokenCmd) StartTokenCommand() {
	context.GetChanInstance().Put(context.TokenChannel, o.command)

	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)

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
		case TokenCmd_CreatNftByAut:
			o.CreateNftbyAut(data.Data, data.Callback)
		case TokenCmd_Bid_Deposit_CheckReceipt:
			o.BidDepositCheckReceipt(data.Data)
		case TokenCmd_Bid_Winner_CheckReceipt:
			o.BidWinnerCheckReceipt(data.Data)
		}

		end := time.Now()

		log.Debug("cmd.kind:", data.CommandType, ",elapsed", end.Sub(start))
	}
	return nil
}

func (o *TokenCmd) CreateNft(data interface{}, cb chan interface{}) {
	product := data.(*context.ProductInfo)

	for i := int64(0); i < product.QuantityTotal; i++ {
		//2-1. nft ?????? ??????

		uri := GetNftUri(o.conf.NftUriDomain, product.Id, i+1)

		if txHash, err := o.itoken.Tokens[Token_nft].Nft_CreateERC721(o.conf.ServerWalletAddr, uri); err != nil {
			log.Error("Nft_CreateERC721 error :", err)
		} else {
			//2-2. db ??????
			if _, err := model.GetDB().InsertProductNFT(product, i+1, context.Nft_state_pending, context.Nft_order_state_sale_ready, txHash, o.conf.ServerWalletAddr, uri); err != nil {
				log.Error("InsertProductNFT :", err)
			}
		}
	}

	cb <- base.MakeBaseResponse(resultcode.Result_Success)
}

func (o *TokenCmd) CreateNftbyAut(data interface{}, cb chan interface{}) {
	product := data.(*context_auc.ProductInfo)

	var err error

	product.NftUri = GetNftUri(o.confCdn.Azure.Domain+o.confCdn.Azure.ContainerNft, product.Id, 0)
	product.NftContract = o.conf.TokenAddrs[Token_nft]
	if product.NftCreateTxHash, err = o.itoken.Tokens[Token_nft].Nft_CreateERC721(o.conf.ServerWalletAddr, product.NftUri); err != nil {
		log.Error("Nft_CreateERC721 error :", err)
		cb <- base.MakeBaseResponse(resultcode.Result_TokenERC721CreateError)
	} else {
		errCnt := 0
	POLLING:
		time.Sleep(time.Second * 5)
		_, isPanding, err := o.itoken.Tokens[Token_onit].eth.GetTransactionByTxHash(product.NftCreateTxHash)
		if err == nil {
			if isPanding {
				log.Debug("is panding : ", isPanding, " tx:", product.NftCreateTxHash)
				time.Sleep(time.Second * 1)
			}
		} else {
			if errCnt > 30 {
				log.Error("GetTransactionByTxHash max try : ", product.NftCreateTxHash)
				cb <- base.MakeBaseResponse(resultcode.Result_TokenERC721CreateError)
				return
			}

			errCnt++
			goto POLLING
		}

		// txhash update
		if _, err := model.GetDB().UpdateAucProductNft(product); err != nil {
			log.Error("UpdateAucProductNft fail : ", err, " product_id:", product.Id, " txhash:", product.NftCreateTxHash)
			cb <- base.MakeBaseResponse(resultcode.Result_DBError)
		} else {
			go func() {
				// nft ????????? json ?????? cdn ?????????
				data := &NftUri_AuctionProduct{
					SNo:      product.SNo,
					Media:    product.Media,
					CardInfo: product.CardInfo,
				}

				nftData, _ := json.Marshal(data)
				productId := strconv.FormatInt(product.Id, 10)
				azure.GetAzure().UploadNftInfoBuffer(nftData, productId+"/"+productId+"_0.json")
			}()

			// product list cache ?????? ??????
			model.GetDB().DeleteProductList()
			model.GetDB().DeleteAuctionList()
			cb <- base.MakeBaseResponse(resultcode.Result_Success)
		}
	}
}

func (o *TokenCmd) DeleteToken(data interface{}) {
	productInfo := data.(*context.ProductInfo)

	_ = productInfo
}

func (o *TokenCmd) OrderProduct(data *basenet.CommandData) {
	go func() {
		order := data.Data.(*context.OrderProduct)
		_, err := model.GetDB().GetProductInfo(order.ProductId)
		if err != nil {
			log.Error("OrderProduct GetProductInfo error ", err, " product_id:")
			return
		}
		token := o.itoken.Tokens[Token_nft]
		errCnt := 0
	POLLING:
		//transaction??? ???????????? ??????
		tx, isPanding, err := token.eth.GetTransactionByTxHash(order.PurchaseTxHash)
		if err == nil {
			if isPanding {
				log.Debug("is panding : ", isPanding, " tx:", order.PurchaseTxHash)
				time.Sleep(time.Second * 1)
				errCnt = 0
				goto POLLING
			}
			// 1. receipt ?????? ??????
			receipt, err := token.eth.GetTransactionReceipt(tx)
			if err == nil {
				log.Info("GetTransactionReceipt Type:", receipt.Type)
				log.Info("GetTransactionReceipt PostState:", receipt.PostState)
				log.Info("GetTransactionReceipt status :", receipt.Status)
				log.Info("GetTransactionReceipt CumulativeGasUsed:", receipt.CumulativeGasUsed)
				//log.Info("GetTransactionReceipt Bloom :", receipt.Bloom)

				log.Info("GetTransactionReceipt topics 0 : ", receipt.Logs[0].Topics[0].Hex())
				log.Info("GetTransactionReceipt topics 1 : ", receipt.Logs[0].Topics[1].Hex())
				log.Info("GetTransactionReceipt topics 2 : ", receipt.Logs[0].Topics[2].Hex())

				log.Info("GetTransactionReceipt TxHash:", receipt.TxHash.Hex())
				log.Info("GetTransactionReceipt contractAddress :", receipt.ContractAddress.Hex())
				log.Info("GetTransactionReceipt GasUsed:", receipt.GasUsed)
				log.Info("GetTransactionReceipt blockhash :", receipt.BlockHash.Hex())
				log.Info("GetTransactionReceipt blocknumber :", receipt.BlockNumber)
				log.Info("GetTransactionReceipt TransactionIndex:", receipt.TransactionIndex)

				for _, logInfo := range receipt.Logs {
					fmt.Printf("GetTransactionReceipt Logs %+v\n", logInfo)
				}

				//token contract address check
				log.Info("token address : ", receipt.Logs[0].Address.Hex())
				if !strings.EqualFold(strings.ToUpper(o.conf.TokenAddrs[Token_onit]), strings.ToUpper(receipt.Logs[0].Address.Hex())) {
					model.GetDB().UpdateProductRemain(true, order.ProductId)
					model.GetDB().UpdateProductNftOrderState(order.TokenId, context.Nft_order_state_sale_ready)
					model.GetDB().UpdateOrderState(order.TokenId, context.Order_state_cancel)
					log.Error("Invalid token address :", receipt.Logs[0].Address.Hex())
					return
				}

				//?????? ?????? ????????? ?????? check
				fromAddr := strings.Replace(receipt.Logs[0].Topics[1].Hex(), "000000000000000000000000", "", -1)
				toAddr := strings.Replace(receipt.Logs[0].Topics[2].Hex(), "000000000000000000000000", "", -1)
				if !strings.EqualFold(strings.ToUpper(order.WalletAddr), strings.ToUpper(fromAddr)) {
					model.GetDB().UpdateProductRemain(true, order.ProductId)
					model.GetDB().UpdateProductNftOrderState(order.TokenId, context.Nft_order_state_sale_ready)
					model.GetDB().UpdateOrderState(order.TokenId, context.Order_state_cancel)
					log.Error("Invalid from address :", fromAddr)
					return
				}
				if !strings.EqualFold(strings.ToUpper(o.conf.ServerWalletAddr), strings.ToUpper(toAddr)) {
					model.GetDB().UpdateProductRemain(true, order.ProductId)
					model.GetDB().UpdateProductNftOrderState(order.TokenId, context.Nft_order_state_sale_ready)
					model.GetDB().UpdateOrderState(order.TokenId, context.Order_state_cancel)
					log.Error("Invalid to address :", toAddr)
					return
				}
				// ?????? ?????? check
				value := new(big.Int)
				value.SetString(hex.EncodeToString(receipt.Logs[0].Data), 16)
				log.Info("transfer value :", value)

				transferEther := ethCtrl.Convert(value.String(), ethCtrl.Wei, ethCtrl.Ether)

				var price big.Rat
				productInfo, err := model.GetDB().GetProductInfo(order.ProductId)
				if err != nil {
					model.GetDB().UpdateProductRemain(true, order.ProductId)
					model.GetDB().UpdateProductNftOrderState(order.TokenId, context.Nft_order_state_sale_ready)
					model.GetDB().UpdateOrderState(order.TokenId, context.Order_state_cancel)
					log.Error("OrderProduct GetProductInfo error ", err, " product_id:")
					return
				}
				for _, pricePos := range productInfo.Prices {
					if strings.EqualFold(pricePos.TokenType, order.TokenType) {
						price = *price.SetFloat64(pricePos.Price)
						break
					}
				}

				temp1, _ := transferEther.Float64()
				temp2, _ := price.Float64()
				if temp1 != temp2 {
					model.GetDB().UpdateProductRemain(true, order.ProductId)
					model.GetDB().UpdateProductNftOrderState(order.TokenId, context.Nft_order_state_sale_ready)
					model.GetDB().UpdateOrderState(order.TokenId, context.Order_state_cancel)
					log.Error("Invalid purchase receipt price :", temp1, " real price :", temp2)
					return
				}
			} else if err.Error() == "not found" {
				log.Debug("not found retry GetTransactionReceipt : ", order.PurchaseTxHash)
				time.Sleep(time.Second * 1)
				if errCnt > 3 {
					model.GetDB().UpdateProductRemain(true, order.ProductId)
					model.GetDB().UpdateProductNftOrderState(order.TokenId, context.Nft_order_state_sale_ready)
					model.GetDB().UpdateOrderState(order.TokenId, context.Order_state_cancel)
					log.Error("GetTransactionReceipt max try from hash : ", order.PurchaseTxHash)
					return
				}
				errCnt++
				goto POLLING
			}
		} else {
			log.Debug("GetTransactionByTxHash error : ", err)
			if errCnt > 3 {
				model.GetDB().UpdateProductRemain(true, order.ProductId)
				model.GetDB().UpdateProductNftOrderState(order.TokenId, context.Nft_order_state_sale_ready)
				model.GetDB().UpdateOrderState(order.TokenId, context.Order_state_cancel)
				log.Error("GetTransactionByTxHash max try : ", order.PurchaseTxHash)
				return
			}
			errCnt++
			goto POLLING
		}

		// 2. ????????? ?????? ?????? ?????? nft ??????
		model.GetDB().UpdateOrderState(order.TokenId, context.Order_state_txhash_complete)
		// ????????? nft token id ??????
		if nfts, err := model.GetDB().GetNftListByProductId(order.ProductId); err != nil {
			model.GetDB().UpdateProductRemain(true, order.ProductId)
			model.GetDB().UpdateProductNftOrderState(order.TokenId, context.Nft_order_state_sale_ready)
			model.GetDB().UpdateOrderState(order.TokenId, context.Order_state_cancel)
			log.Error("GetNftListByProductId error : ", err)
			return
		} else {
			for _, nft := range nfts {
				if nft.TokenId == order.TokenId {
					//nft ??????
					txHash, err := o.itoken.Tokens[Token_nft].Nft_TransferERC721(nft.OwnerWalletAddr, order.WalletAddr, nft.TokenId)
					if err != nil {
						log.Error("Nft_TransferERC721 error : ", err, " token_id:", nft.TokenId)
						model.GetDB().UpdateProductRemain(true, order.ProductId)
						model.GetDB().UpdateProductNftOrderState(order.TokenId, context.Nft_order_state_sale_ready)
						model.GetDB().UpdateOrderState(order.TokenId, context.Order_state_cancel)
						return
					} else {
						model.GetDB().UpdateOrderState(order.TokenId, context.Order_state_nft_transfer_start)
						log.Info("Nft_TransferERC721 txhash : ", txHash)
					}
					break
				}
			}
		}
	}()

	// 3. nft ?????? ??????????????? db ???????????? ??????
}
