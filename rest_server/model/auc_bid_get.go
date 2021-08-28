package model

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

// 최고가 입찰 참여자
func (o *DB) GetAucBidBestAttendee(aucId int64) (*context_auc.Bid, error) {
	query := fmt.Sprintf("SELECT * FROM auc_bids WHERE auc_id=%v ORDER BY bid_amount DESC limit 1", aucId)
	rows, err := o.Mysql.Query(query)

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

// 특정 회원이 입찰에 참여한 리스트
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

// 입찰자 리스트
func (o *DB) GetAucBidAttendeeList(pageInfo *context_auc.BidAttendeeList) ([]context_auc.Bid, int64, error) {
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
			log.Error("GetAucBidAttendeeList::ScanBid error : ", err)
			continue
		}
		bids = append(bids, *bid)
	}

	totalCount, err := o.GetTotalAucBidSize(pageInfo.AucId)

	return bids, totalCount, err
}

// tx hash 사용 여부 확인용
func (o *DB) GetAucBidByTxhash(txHash string) (bool, error) {
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
			log.Error("GetAucBidByTxhash::ScanBid error : ", err)
			continue
		}
		cnt++
	}

	if cnt == 0 {
		return false, nil
	}
	return true, nil
}

// 입찰 보증금 반환 리스트
func (o *DB) GetAucBidDepositRefund(req *context_auc.BidDepositRefundList) ([]context_auc.Bid, int64, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM auc_bids WHERE bid_amount != 0 and auc_id=%v and bid_state = 1 and "+
		"bid_winner_state = 0 and deposit_state = 2 "+
		"ORDER BY bid_amount DESC LIMIT %v,%v", req.AucId, req.PageSize*req.PageOffset, req.PageSize)
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
			log.Error("GetAucBidDepositRefund::ScanBid error : ", err)
			continue
		}
		bids = append(bids, *bid)
	}

	totalCount, err := o.GetTotalAucBidDepositRefund(req)

	return bids, totalCount, err
}

func (o *DB) GetTotalAucBidSize(aucId int64) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) as count FROM auc_bids WHERE bid_amount != 0 AND auc_id=%v", aucId)

	err := o.Mysql.QueryRow(query, &count)

	if err != nil {
		log.Error(err)
		return count, err
	}

	return count, nil
}

func (o *DB) GetTotalAucBidDepositRefund(req *context_auc.BidDepositRefundList) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) as count FROM auc_bids WHERE bid_amount != 0 and auc_id=%v and bid_state = 1 and "+
		"bid_winner_state = 0 and deposit_state = 2 "+
		"ORDER BY bid_amount DESC", req.AucId)
	err := o.Mysql.QueryRow(query, &count)

	if err != nil {
		log.Error(err)
		return count, err
	}

	return count, nil
}

func (o *DB) ScanBid(rows *sql.Rows) (*context_auc.Bid, error) {
	var bidWinnerTxHash, tos sql.NullString

	bid := &context_auc.Bid{}
	if err := rows.Scan(&bid.Id, &bid.AucId, &bid.ProductId,
		&bid.BidState, &bid.BidTs, &bid.BidAttendeeWalletAddr, &bid.BidAmount, &bidWinnerTxHash, &bid.BidWinnerState,
		&bid.DepositAmount, &bid.DepositTxHash, &bid.DepositState, &bid.TokenType, &tos); err != nil {
		//log.Error("ScanBid error :", err)
		return nil, err
	}

	bid.BidWinnerTxHash = bidWinnerTxHash.String

	getTos := context_auc.TermsOfService{}
	json.Unmarshal([]byte(tos.String), &getTos)
	bid.TermsOfService = getTos

	return bid, nil
}
