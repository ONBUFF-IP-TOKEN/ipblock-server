package token

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	ethCtrl "github.com/ONBUFF-IP-TOKEN/baseEthereum/ethcontroller"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
)

func (o *Token) Onit_LoadContract(tokenAddr string) error {
	err := o.eth.Onit_LoadContract(tokenAddr)
	return err
}

func (o *Token) Onit_LoadContractInfo() error {
	var err error
	o.tokenName, err = o.eth.Onit_GetName()
	if err != nil {
		return err
	}
	o.tokenSymbol, err = o.eth.Onit_GetSymbol()
	if err != nil {
		return err
	}
	return err
}

func (o *Token) Onit_GetBalanceOf(walletAddr string) (int64, error) {
	bal, err := o.eth.Onit_GetBalanceOf(walletAddr)
	ne := big.NewInt(1000000000000000000)

	baseBal := big.NewInt(0)
	baseBal = baseBal.Div(bal, ne)

	return baseBal.Int64(), err
}

func (o *Token) CheckTransferReceipt(purchaseInfo *context.PostPurchaseItem, itemInfo *context.ItemInfo) {
	go func() {
		errCnt := 0
	POLLING:
		tx, isPanding, err := o.eth.GetTransactionByTxHash(purchaseInfo.PurchaseTxHash)
		if err == nil {
			if isPanding {
				log.Debug("is panding : ", isPanding)
				time.Sleep(time.Second * 1)
				errCnt = 0
				goto POLLING
			}
			receipt, err := o.eth.GetTransactionReceipt(tx)
			if err == nil {
				log.Info("GetTransactionReceipt Type:", receipt.Type)
				log.Info("GetTransactionReceipt PostState:", receipt.PostState)
				log.Info("GetTransactionReceipt status :", receipt.Status)
				log.Info("GetTransactionReceipt CumulativeGasUsed:", receipt.CumulativeGasUsed)
				log.Info("GetTransactionReceipt Bloom :", receipt.Bloom)
				for _, logInfo := range receipt.Logs {
					fmt.Printf("GetTransactionReceipt Logs %+v\n", logInfo)
				}
				log.Info("GetTransactionReceipt TxHash:", receipt.TxHash.Hex())
				log.Info("GetTransactionReceipt contractAddress :", receipt.ContractAddress.Hex())
				log.Info("GetTransactionReceipt GasUsed:", receipt.GasUsed)
				log.Info("GetTransactionReceipt blockhash :", receipt.BlockHash.Hex())
				log.Info("GetTransactionReceipt blocknumber :", receipt.BlockNumber)
				log.Info("GetTransactionReceipt TransactionIndex:", receipt.TransactionIndex)

				//token contract address check
				log.Info("token address : ", receipt.Logs[0].Address.Hex())
				if strings.EqualFold(strings.ToUpper(o.conf.TokenAddrs[o.TokenType]), strings.ToUpper(receipt.Logs[0].Address.Hex())) {
					log.Error("Invalid token address :", receipt.Logs[0].Address.Hex())
					return
				}

				//받는 사람 보내는 사람 check
				fromAddr := strings.Replace(receipt.Logs[0].Topics[1].Hex(), "000000000000000000000000", "", -1)
				toAddr := strings.Replace(receipt.Logs[0].Topics[2].Hex(), "000000000000000000000000", "", -1)
				if strings.EqualFold(strings.ToUpper(purchaseInfo.WalletAddr), strings.ToUpper(fromAddr)) {
					log.Error("Invalid from address :", fromAddr)
					return
				}
				if strings.EqualFold(strings.ToUpper(o.conf.ServerWalletAddr), strings.ToUpper(toAddr)) {
					log.Error("Invalid to address :", toAddr)
					return
				}

				// 구입 액수 check
				value := new(big.Int)
				value.SetString(hex.EncodeToString(receipt.Logs[0].Data), 16)
				log.Info("transfer value :", value)

				transferEther := ethCtrl.Convert(value.String(), ethCtrl.Wei, ethCtrl.Ether)
				price := new(big.Rat).SetFloat64(itemInfo.TokenPrice)
				if transferEther.Cmp(price) != 0 {
					log.Error("Invalid purchase price :", transferEther.String())
					return
				}

				// 인증 토큰에서 지갑주소 추출해서 item_id와 owner_wallet_address 추출
				txHash, err := GetToken().Tokens[Token_nft].Nft_TransferERC721(itemInfo.OwnerWalletAddr, purchaseInfo.WalletAddr, itemInfo.TokenId)
				if err != nil {
					log.Error("Nft_TransferERC721 error : ", err)
				} else {
					log.Error("Nft_TransferERC721 txhash : ", txHash)
				}

				// nft 정보 생성
				// itemInfo, err := model.GetDB().GetEventItem(purchaseInfo.ItemNum)
				// if err != nil {
				// 	log.Error("UpdateEventVerifyPurchase error : ", err)
				// }
				// nftInfo := &NftUriInfo{
				// 	UriType:  "shose",
				// 	CreateTs: datetime.GetTS2MilliSec(),
				// }
				// nftShose := &NftUriShose{
				// 	Idx:           itemInfo.Idx,
				// 	Name:          itemInfo.Name,
				// 	SerialNo:      itemInfo.Serial,
				// 	Info:          itemInfo.Info,
				// 	Certification: "http://file.onbuff.com/onif/shose/certifi/1.pdf",
				// }
				// nftInfo.Data = nftShose
				// nftFile, err := json.MarshalIndent(nftInfo, "", " ")
				// if err != nil {
				// 	log.Error("MarshalIndent error ", err)
				// }
				// _ = ioutil.WriteFile(strconv.FormatInt(purchaseInfo.ItemNum, 10)+".json", nftFile, 0644)
			} else if err.Error() == "not found" {
				log.Debug("not found retry GetTransactionReceipt : ", purchaseInfo.PurchaseTxHash)
				time.Sleep(time.Second * 1)
				if errCnt > 3 {
					log.Error("GetTransactionReceipt max try")
					return
				}
				errCnt++
				goto POLLING
			}
		} else {

			if errCnt > 3 {
				log.Error("GetTransactionByTxHash max try : ", purchaseInfo.PurchaseTxHash)
				return
			}
			time.Sleep(time.Second * 2)
			errCnt++
			log.Debug("GetTransactionByTxHash error : ", err, " 재시도:", errCnt)
			goto POLLING
		}
	}()
}
