package model

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

func (o *DB) InsertAucBid(bid *context_auc.BidDeposit) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO auc_bids (auc_id, product_id, " +
		"bid_state, bid_ts, bid_attendee_wallet_address, " +
		"deposit_amount, deposit_txhash, deposit_state, token_type ) VALUES (?,?,?,?,?,?,?,?,?)")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, bid.AucId, bid.ProductId,
		bid.BidState, bid.BidTs, bid.BidAttendeeWalletAddr, bid.DepositAmount, bid.DepositTxHash, bid.DepositState, bid.TokenType)

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
	o.DeleteBidList(bid.AucId)
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
	o.DeleteBidList(bidSubmit.AucId)
	return id, nil
}

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
	// if id == 0 {
	// 	err = errors.New("RowsAffected none")
	// 	log.Error(err)
	// 	return id, err
	// }
	log.Debug("UpdateAucBidWinner id:", id)
	o.DeleteBidList(bid.AucId)
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

	var bidWinnerTxHash sql.NullString
	bid := &context_auc.Bid{}
	for rows.Next() {
		if err := rows.Scan(&bid.Id, &bid.AucId, &bid.ProductId,
			&bid.BidState, &bid.BidTs, &bid.BidAttendeeWalletAddr, &bid.BidAmount, &bidWinnerTxHash, &bid.BidWinnerState,
			&bid.DepositAmount, &bid.DepositTxHash, &bid.DepositState, &bid.TokenType); err != nil {
			log.Error(err)
		} else {
			bid.BidWinnerTxHash = bidWinnerTxHash.String
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

	var bidWinnerTxHash sql.NullString
	bid := &context_auc.Bid{}
	for rows.Next() {
		if err := rows.Scan(&bid.Id, &bid.AucId, &bid.ProductId,
			&bid.BidState, &bid.BidTs, &bid.BidAttendeeWalletAddr, &bid.BidAmount, &bidWinnerTxHash, &bid.BidWinnerState,
			&bid.DepositAmount, &bid.DepositTxHash, &bid.DepositState, &bid.TokenType); err != nil {
			log.Error(err)
		} else {
			bid.BidWinnerTxHash = bidWinnerTxHash.String
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

	var bidWinnerTxHash sql.NullString
	bids := make([]context_auc.Bid, 0)
	for rows.Next() {
		bid := context_auc.Bid{}
		if err := rows.Scan(&bid.Id, &bid.AucId, &bid.ProductId,
			&bid.BidState, &bid.BidTs, &bid.BidAttendeeWalletAddr, &bid.BidAmount, &bidWinnerTxHash, &bid.BidWinnerState,
			&bid.DepositAmount, &bid.DepositTxHash, &bid.DepositState, &bid.TokenType); err != nil {
			log.Error(err)
		} else {
			bid.BidWinnerTxHash = bidWinnerTxHash.String
			bids = append(bids, bid)
		}
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

	var bidWinnerTxHash sql.NullString
	cnt := 0
	for rows.Next() {
		bid := context_auc.Bid{}
		if err := rows.Scan(&bid.Id, &bid.AucId, &bid.ProductId,
			&bid.BidState, &bid.BidTs, &bid.BidAttendeeWalletAddr, &bid.BidAmount, &bidWinnerTxHash, &bid.BidWinnerState,
			&bid.DepositAmount, &bid.DepositTxHash, &bid.DepositState, &bid.TokenType); err != nil {
			log.Error(err)
		} else {
			cnt++
		}
	}

	if cnt == 0 {
		return false, nil
	}
	return true, nil
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
