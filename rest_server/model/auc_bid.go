package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

// 입찰 하기
func (o *DB) InsertAucBidSubmit(bidSubmit *context_auc.BidSubmit) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO auc_bids (auc_id, product_id, bid_state, bid_ts, bid_attendee_wallet_address, " +
		"bid_amount, token_type ) VALUES (?,?,?,?,?,?,?)")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, bidSubmit.AucId, bidSubmit.ProductId, bidSubmit.BidState, bidSubmit.BidTs, bidSubmit.BidAttendeeWalletAddr,
		bidSubmit.BidAmount, bidSubmit.TokenType)

	if err != nil {
		log.Error(err)
		return -1, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	log.Debug("InsertAucBidSubmit id:", id)
	o.DeleteBidList(bidSubmit.AucId)
	return id, nil
}

// 입찰 정보 삭제
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
