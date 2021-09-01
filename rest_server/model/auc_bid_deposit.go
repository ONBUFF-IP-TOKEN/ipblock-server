package model

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

// 입찰금 입력
func (o *DB) InsertAucBidDepsit(bidDeposit *context_auc.BidDepositSubmit) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO auc_bids_deposit (auc_id, product_id, bid_attendee_wallet_address, " +
		"deposit_amount, deposit_txhash, deposit_state, deposit_ts, " +
		"token_type, terms_of_service ) VALUES (?,?,?,?,?,?,?,?,?)")

	tos, _ := json.Marshal(bidDeposit.TermsOfService)

	result, err := o.Mysql.PrepareAndExec(sqlQuery, bidDeposit.AucId, bidDeposit.ProductId, bidDeposit.BidAttendeeWalletAddr,
		bidDeposit.DepositAmount, bidDeposit.DepositTxHash, bidDeposit.DepositState, bidDeposit.DepositTs,
		bidDeposit.TokenType, string(tos))

	if err != nil {
		log.Error(err)
		return -1, err
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	log.Debug("InsertAucBidDepsit id:", insertId)
	o.DeleteBidList(bidDeposit.AucId)
	return insertId, nil
}

// 입찰금 입금 확인
func (o *DB) GetAucBidDeposit(aucId int64, walletAddr string) (*context_auc.BidDeposit, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM auc_bids_deposit WHERE auc_id=%v and bid_attendee_wallet_address='%v'", aucId, walletAddr)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	bidDeposit := &context_auc.BidDeposit{}
	for rows.Next() {
		bidDeposit, err = o.ScanBidDeposit(rows)
		if err != nil {
			log.Error("GetAucBidDeposit::ScanBidDeposit error : ", err)
			continue
		}
	}

	if len(bidDeposit.BidAttendeeWalletAddr) == 0 {
		return nil, err
	}
	return bidDeposit, err
}

// 입찰 보증금 전송 상태 업데이트
func (o *DB) UpdateAucBidDepositState(bidDeposit *context_auc.BidDepositSubmit, state int64) (int64, error) {
	sqlQuery := fmt.Sprintf("UPDATE auc_bids_deposit set deposit_state=? WHERE id=?")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, state, bidDeposit.Id)

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
	o.DeleteBidList(bidDeposit.AucId)
	return id, nil
}

// 입찰 보증금 정보 삭제
func (o *DB) DeleteAucBidDeposit(bid *context_auc.BidRemove) (bool, error) {
	sqlQuery := "DELETE FROM auc_bids_deposit WHERE auc_id=? AND bid_attendee_wallet_address=?"

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

func (o *DB) ScanBidDeposit(rows *sql.Rows) (*context_auc.BidDeposit, error) {
	var tos sql.NullString

	bidDeposit := &context_auc.BidDeposit{}
	if err := rows.Scan(&bidDeposit.Id, &bidDeposit.AucId, &bidDeposit.ProductId, &bidDeposit.BidAttendeeWalletAddr,
		&bidDeposit.DepositAmount, &bidDeposit.DepositTxHash, &bidDeposit.DepositState, &bidDeposit.DepositTs, &bidDeposit.TokenType, &tos); err != nil {
		return nil, err
	}

	getTos := context_auc.TermsOfService{}
	json.Unmarshal([]byte(tos.String), &getTos)
	bidDeposit.TermsOfService = getTos

	return bidDeposit, nil
}
