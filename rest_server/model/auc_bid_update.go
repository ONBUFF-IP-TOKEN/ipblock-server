package model

import (
	"errors"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

//낙찰자 정보 업데이트
func (o *DB) UpdateAucBidWinner(bid *context_auc.BidWinner) (int64, error) {
	sqlQuery := fmt.Sprintf("UPDATE auc_bids set bid_winner_txhash=?, bid_winner_state=? WHERE auc_id=? and bid_attendee_wallet_address=?")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, bid.BidWinnerTxHash, bid.BidWinnerState, bid.AucId, bid.BidAttendeeWalletAddr)

	if err != nil {
		log.Error(err)
		return -1, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	if id == 0 {
		err = errors.New("RowsAffected none")
		log.Error(err)
		return id, err
	}
	log.Debug("UpdateAucBidWinner id:", id)
	o.DeleteBidList(bid.AucId)
	return id, nil
}

// 낙찰자 상태 정보 업데이트
func (o *DB) UpdateAucBidWinnerState(bid *context_auc.Bid, state int) (int64, error) {
	sqlQuery := fmt.Sprintf("UPDATE auc_bids set bid_winner_state=? WHERE auc_id=? and bid_attendee_wallet_address=?")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, state, bid.AucId, bid.BidAttendeeWalletAddr)

	if err != nil {
		log.Error(err)
		return -1, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return -1, err
	}

	log.Debug("UpdateAucBidWinner id:", id)
	o.DeleteBidList(bid.AucId)
	return id, nil
}

// 낙찰자 지정하기
func (o *DB) UpdateAucBidFinish(bid *context_auc.Bid, state context_auc.Bid_state) (int64, error) {
	sqlQuery := fmt.Sprintf("UPDATE auc_bids set bid_state=? WHERE id=? AND auc_id=? AND bid_attendee_wallet_address=?")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, state, bid.Id, bid.AucId, bid.BidAttendeeWalletAddr)

	if err != nil {
		log.Error(err)
		return -1, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return -1, err
	}

	log.Debug("UpdateAucBidWinner id:", id)
	o.DeleteBidList(bid.AucId)
	o.DeleteAuctionCache(bid.AucId)
	return id, nil
}
