package token

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	ethCtrl "github.com/ONBUFF-IP-TOKEN/baseEthereum/ethcontroller"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
)

func (o *TokenCmd) BidDepositCheckReceipt(data interface{}) {
	go func() {
		bid := data.(*context_auc.BidDeposit)
		token := o.itoken.Tokens[Token_onit]
		errCnt := 0
	POLLING:
		//transaction이 정상인지 체크
		tx, isPanding, err := token.eth.GetTransactionByTxHash(bid.DepositTxHash)
		if err == nil {
			if isPanding {
				log.Debug("is panding : ", isPanding, " tx:", bid.DepositTxHash)
				time.Sleep(time.Second * 1)
				errCnt = 0
				goto POLLING
			}

			// 1. receipt 정상 체크
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

				// 영수증 유효성 체크
				if receipt.Status == 0 {
					model.GetDB().UpdateAucBidDepositState(bid, context_auc.Deposit_state_fail)
					log.Error("receipt.Status :", receipt.Status)
					return
				}

				//token contract address check
				log.Info("token address : ", receipt.Logs[0].Address.Hex())
				if !strings.EqualFold(strings.ToUpper(o.conf.TokenAddrs[Token_onit]), strings.ToUpper(receipt.Logs[0].Address.Hex())) {
					model.GetDB().UpdateAucBidDepositState(bid, context_auc.Deposit_state_fail)
					log.Error("Invalid token address :", receipt.Logs[0].Address.Hex())
					return
				}

				//받는 사람 보내는 사람 check
				fromAddr := strings.Replace(receipt.Logs[0].Topics[1].Hex(), "000000000000000000000000", "", -1)
				toAddr := strings.Replace(receipt.Logs[0].Topics[2].Hex(), "000000000000000000000000", "", -1)
				if !strings.EqualFold(strings.ToUpper(bid.BidAttendeeWalletAddr), strings.ToUpper(fromAddr)) {
					model.GetDB().UpdateAucBidDepositState(bid, context_auc.Deposit_state_fail)
					log.Error("Invalid from address :", fromAddr)
					return
				}
				if !strings.EqualFold(strings.ToUpper(o.conf.ServerWalletAddr), strings.ToUpper(toAddr)) {
					model.GetDB().UpdateAucBidDepositState(bid, context_auc.Deposit_state_fail)
					log.Error("Invalid to address :", toAddr)
					return
				}
				// 구입 액수 check
				value := new(big.Int)
				value.SetString(hex.EncodeToString(receipt.Logs[0].Data), 16)
				log.Info("transfer value :", value)

				transferEther := ethCtrl.Convert(value.String(), ethCtrl.Wei, ethCtrl.Ether)

				var price big.Rat
				price = *price.SetFloat64(bid.BidAmount)

				temp1, _ := transferEther.Float64()
				temp2, _ := price.Float64()
				if temp1 != temp2 {
					model.GetDB().UpdateAucBidDepositState(bid, context_auc.Deposit_state_fail)
					log.Error("Invalid purchase receipt price :", temp1, " real price :", temp2)
					return
				}
			} else if err.Error() == "not found" {
				log.Debug("not found retry GetTransactionReceipt : ", bid.DepositTxHash, " bid id:", bid.Id)
				time.Sleep(time.Second * 1)
				if errCnt > 3 {
					model.GetDB().UpdateAucBidDepositState(bid, context_auc.Deposit_state_fail)
					log.Error("GetTransactionReceipt max try from hash : ", bid.DepositTxHash, " bid id:", bid.Id)
					return
				}
				errCnt++
				goto POLLING
			}
		} else {
			log.Debug("GetTransactionByTxHash error : ", err)
			if errCnt > 3 {
				model.GetDB().UpdateAucBidDepositState(bid, context_auc.Deposit_state_fail)
				log.Error("GetTransactionByTxHash max try : ", bid.DepositTxHash, " bid id:", bid.Id)
				return
			}
			errCnt++
			goto POLLING
		}

		// 정상 처리 되었으면 입찰자 정보 업데이트
		_, _ = model.GetDB().UpdateAucBidDepositState(bid, context_auc.Deposit_state_complete)
	}()
}

func (o *TokenCmd) BidSuccessCheckReceipt(data interface{}) {
	go func() {
		bid := data.(*context_auc.BidWinner)
		token := o.itoken.Tokens[Token_onit]
		errCnt := 0
	POLLING:
		//transaction이 정상인지 체크
		tx, isPanding, err := token.eth.GetTransactionByTxHash(bid.BidWinnerTxHash)
		if err == nil {
			if isPanding {
				log.Debug("is panding : ", isPanding, " tx:", bid.DepositTxHash)
				time.Sleep(time.Second * 1)
				errCnt = 0
				goto POLLING
			}

			// 1. receipt 정상 체크
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

				// 영수증 유효성 체크
				if receipt.Status == 0 {
					model.GetDB().UpdateAucBidWinnerState(&bid.Bid, context_auc.Bid_winner_state_none)
					log.Error("receipt.Status :", receipt.Status)
					return
				}

				//token contract address check
				log.Info("token address : ", receipt.Logs[0].Address.Hex())
				if !strings.EqualFold(strings.ToUpper(o.conf.TokenAddrs[Token_onit]), strings.ToUpper(receipt.Logs[0].Address.Hex())) {
					model.GetDB().UpdateAucBidWinnerState(&bid.Bid, context_auc.Bid_winner_state_none)
					log.Error("Invalid token address :", receipt.Logs[0].Address.Hex())
					return
				}

				//받는 사람 보내는 사람 check
				fromAddr := strings.Replace(receipt.Logs[0].Topics[1].Hex(), "000000000000000000000000", "", -1)
				toAddr := strings.Replace(receipt.Logs[0].Topics[2].Hex(), "000000000000000000000000", "", -1)
				if !strings.EqualFold(strings.ToUpper(bid.BidAttendeeWalletAddr), strings.ToUpper(fromAddr)) {
					model.GetDB().UpdateAucBidWinnerState(&bid.Bid, context_auc.Bid_winner_state_none)
					log.Error("Invalid from address :", fromAddr)
					return
				}
				if !strings.EqualFold(strings.ToUpper(o.conf.ServerWalletAddr), strings.ToUpper(toAddr)) {
					model.GetDB().UpdateAucBidWinnerState(&bid.Bid, context_auc.Bid_winner_state_none)
					log.Error("Invalid to address :", toAddr)
					return
				}
				// 구입 액수 check
				value := new(big.Int)
				value.SetString(hex.EncodeToString(receipt.Logs[0].Data), 16)
				log.Info("transfer value :", value)

				transferEther := ethCtrl.Convert(value.String(), ethCtrl.Wei, ethCtrl.Ether)

				var price big.Rat
				price = *price.SetFloat64(bid.DepositAmount)

				temp1, _ := transferEther.Float64()
				temp2, _ := price.Float64()
				if temp1 != temp2 {
					model.GetDB().UpdateAucBidWinnerState(&bid.Bid, context_auc.Bid_winner_state_none)
					log.Error("Invalid purchase receipt price :", temp1, " real price :", temp2)
					return
				}
			} else if err.Error() == "not found" {
				log.Debug("not found retry GetTransactionReceipt : ", bid.DepositTxHash, " bid id:", bid.Id)
				time.Sleep(time.Second * 1)
				if errCnt > 3 {
					model.GetDB().UpdateAucBidWinnerState(&bid.Bid, context_auc.Bid_winner_state_none)
					log.Error("GetTransactionReceipt max try from hash : ", bid.DepositTxHash, " bid id:", bid.Id)
					return
				}
				errCnt++
				goto POLLING
			}
		} else {
			log.Debug("GetTransactionByTxHash error : ", err)
			if errCnt > 3 {
				model.GetDB().UpdateAucBidWinnerState(&bid.Bid, context_auc.Bid_winner_state_none)
				log.Error("GetTransactionByTxHash max try : ", bid.DepositTxHash, " bid id:", bid.Id)
				return
			}
			errCnt++
			goto POLLING
		}

		// 정상 처리 되었으면 입찰자 정보 업데이트
		model.GetDB().UpdateAucBidWinnerState(&bid.Bid, context_auc.Bid_winner_state_submit_complete)
	}()
}
