package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

func (o *DB) InsertAucAuction(auction *context_auc.AucAuctionRegister) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO auc_auctions (bid_start_amount, bid_cur_amount, bid_unit, bid_deposit, " +
		"auc_start_ts, auc_end_ts, auc_state, auc_round," +
		"create_ts, active_state, product_id, recommand" +
		") VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, auction.BidStartAmount, auction.BidCurAmount, auction.BidUnit, auction.BidDeposit,
		auction.AucStartTs, auction.AucEndTs, auction.AucState, auction.AucRound,
		auction.CreateTs, auction.ActiveState, auction.ProductId, auction.Recommand)

	if err != nil {
		log.Error(err)
		return -1, err
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	log.Debug("InsertAucAuction id:", insertId)
	// auction list cache 전체 삭제
	o.DeleteAuctionList()
	return insertId, nil
}

func (o *DB) DeleteAucAuction(auctionId int64) (bool, error) {
	sqlQuery := "DELETE FROM auc_auctions WHERE auc_id=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, auctionId)
	if err != nil {
		log.Error(err)
		return false, err
	}
	cnt, err := result.RowsAffected()
	if cnt == 0 {
		log.Error(err)
		return false, err
	}

	// auction list cache 전체 삭제
	o.DeleteAuctionList()
	o.DeleteAuctionCache(auctionId)

	return true, nil
}
