package model

import (
	"encoding/json"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

func (o *DB) InsertAucBid(bid *context_auc.BidDeposit) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO auc_bids (auc_id, product_id, " +
		"bid_state, bid_ts, bid_attendee_wallet_address, " +
		"deposit_amount, deposit_txhash, deposit_state, token_type, terms_of_service ) VALUES (?,?,?,?,?,?,?,?,?,?)")

	tos, _ := json.Marshal(bid.TermsOfService)

	result, err := o.Mysql.PrepareAndExec(sqlQuery, bid.AucId, bid.ProductId,
		bid.BidState, bid.BidTs, bid.BidAttendeeWalletAddr, bid.DepositAmount, bid.DepositTxHash, bid.DepositState, bid.TokenType, string(tos))

	if err != nil {
		log.Error(err)
		return -1, err
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	log.Debug("InsertAucBid id:", insertId)
	o.DeleteBidList(bid.AucId)
	return insertId, nil
}

func (o *DB) DeleteAucBid(bid *context_auc.BidRemove) (bool, error) {
	sqlQuery := "DELETE FROM auc_bids WHERE auc_id=? AND bid_attendee_wallet_address=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, bid.AucId, bid.BidAttendeeWalletAddr)
	if err != nil {
		log.Error(err)
		return false, err
	}
	cnt, err := result.RowsAffected()
	if cnt == 0 {
		log.Error(err)
		return false, err
	}

	// bid list cache 전체 삭제
	o.DeleteBidList(bid.AucId)

	return true, nil
}
