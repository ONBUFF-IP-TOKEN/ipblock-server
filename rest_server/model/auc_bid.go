package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

func (o *DB) InsertAucBid(bid *context_auc.BidDeposit) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO auc_bids (auc_id, product_id, " +
		"bid_state, bid_ts, bid_attendee_wallet_address, " +
		"deposit_amount, deposit_txhash, deposit_state ) VALUES (?,?,?,?,?,?,?,?)")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, bid.AucId, bid.ProductId,
		bid.BidState, bid.BidTs, bid.BidAttendeeWalletAddr, bid.DepositAmount, bid.DepositTxHash, bid.DepositState)

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
	return insertId, nil
}

func (o *DB) UpdateAucBidDepositState(bid *context_auc.BidDeposit, state int64) (int64, error) {
	sqlQuery := fmt.Sprintf("UPDATE auc_bids set deposit_state=? WHERE id=?")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, state, bid.Id)

	if err != nil {
		log.Error(err)
		return -1, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	log.Debug("UpdateAucBidDepositState id:", id, " deposit_state:", state)

	return id, nil
}

func (o *DB) UpdateAucBidSubmit(bidSubmit *context_auc.BidSubmit) (int64, error) {
	sqlQuery := fmt.Sprintf("UPDATE auc_bids set bid_state=?, bid_amount=? WHERE auc_id=? and bid_attendee_wallet_address=?")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, bidSubmit.BidState, bidSubmit.BidAmount, bidSubmit.AucId, bidSubmit.BidAttendeeWalletAddr)

	if err != nil {
		log.Error(err)
		return -1, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	log.Debug("UpdateAucBidSubmit id:", id)

	return id, nil
}

func (o *DB) GetAucBidBestAttendee(aucId int64) (*context_auc.Bid, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM auc_bids ORDER BY bid_amount DESC limit 1")
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	bid := &context_auc.Bid{}
	for rows.Next() {
		if err := rows.Scan(&bid.Id, &bid.AucId, &bid.ProductId,
			&bid.BidState, &bid.BidTs, &bid.BidAttendeeWalletAddr, &bid.BidAmount,
			&bid.DepositAmount, &bid.DepositTxHash, &bid.DepositState); err != nil {
			log.Error(err)
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
		if err := rows.Scan(&bid.Id, &bid.AucId, &bid.ProductId,
			&bid.BidState, &bid.BidTs, &bid.BidAttendeeWalletAddr, &bid.BidAmount,
			&bid.DepositAmount, &bid.DepositTxHash, &bid.DepositState); err != nil {
			log.Error(err)
		}
	}

	if len(bid.BidAttendeeWalletAddr) == 0 {
		return nil, err
	}
	return bid, err
}

func (o *DB) GetAucBidBestAttendeeList(pageInfo *context_auc.BidAttendeeList) ([]context_auc.Bid, int64, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM auc_bids WHERE bid_amount != 0 ORDER BY bid_amount DESC LIMIT %v,%v", pageInfo.PageSize*pageInfo.PageOffset, pageInfo.PageSize)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	defer rows.Close()

	bids := make([]context_auc.Bid, 0)
	for rows.Next() {
		bid := context_auc.Bid{}
		if err := rows.Scan(&bid.Id, &bid.AucId, &bid.ProductId,
			&bid.BidState, &bid.BidTs, &bid.BidAttendeeWalletAddr, &bid.BidAmount,
			&bid.DepositAmount, &bid.DepositTxHash, &bid.DepositState); err != nil {
			log.Error(err)
		} else {
			bids = append(bids, bid)
		}
	}

	totalCount, err := o.GetTotalAucBidSize()

	return bids, totalCount, err
}

func (o *DB) GetTotalAucBidSize() (int64, error) {
	rows, err := o.Mysql.Query("SELECT COUNT(*) as count FROM auc_bids WHERE bid_amount != 0")
	var count int64
	if err != nil {
		log.Error(err)
		return count, err
	}

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&count)
	}

	return count, nil
}
