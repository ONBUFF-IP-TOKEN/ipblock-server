package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

func (o *DB) GetAucBidBestAttendee(aucId int64) (*context_auc.Bid, error) {
	rows, err := o.Mysql.Query("SELECT * FROM auc_bids ORDER BY bid_amount DESC limit 1")

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	bid := &context_auc.Bid{}
	for rows.Next() {
		bid, err = o.ScanBid(rows)
		if err != nil {
			log.Error("GetAucBidBestAttendee::ScanBid error : ", err)
			continue
		}
	}

	if len(bid.BidAttendeeWalletAddr) == 0 {
		return nil, err
	}
	return bid, err
}

func (o *DB) GetAucBidAttendee(aucId int64, walletAddr string) (*context_auc.Bid, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM auc_bids WHERE auc_id=%v and bid_attendee_wallet_address='%v'", aucId, walletAddr)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	bid := &context_auc.Bid{}
	for rows.Next() {
		bid, err = o.ScanBid(rows)
		if err != nil {
			log.Error("GetAucBidAttendee::ScanBid error : ", err)
			continue
		}
	}

	if len(bid.BidAttendeeWalletAddr) == 0 {
		return nil, err
	}
	return bid, err
}

func (o *DB) GetAucBidBestAttendeeList(pageInfo *context_auc.BidAttendeeList) ([]context_auc.Bid, int64, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM auc_bids WHERE bid_amount != 0 and auc_id=%v ORDER BY bid_amount DESC LIMIT %v,%v", pageInfo.AucId, pageInfo.PageSize*pageInfo.PageOffset, pageInfo.PageSize)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	defer rows.Close()

	bids := make([]context_auc.Bid, 0)
	for rows.Next() {
		bid, err := o.ScanBid(rows)
		if err != nil {
			log.Error("GetAucBidBestAttendeeList::ScanBid error : ", err)
			continue
		}
		bids = append(bids, *bid)
	}

	totalCount, err := o.GetTotalAucBidSize()

	return bids, totalCount, err
}

func (o *DB) GetAucBidBestAttendeeByTxhash(txHash string) (bool, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM auc_bids WHERE deposit_txhash = '%v' OR bid_winner_txhash = '%v'", txHash, txHash)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return false, err
	}

	defer rows.Close()

	cnt := 0
	for rows.Next() {
		_, err := o.ScanBid(rows)
		if err != nil {
			log.Error("GetAucBidBestAttendeeByTxhash::ScanBid error : ", err)
			continue
		}
		cnt++
	}

	if cnt == 0 {
		return false, nil
	}
	return true, nil
}
