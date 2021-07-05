package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
)

func (o *DB) GetHistoryTransferItem(req *context.GetHistoryTransferItem) ([]context.ItemTransferHistory, int64, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM ipblock.items_history_transfer WHERE item_id=%v ORDER BY idx DESC LIMIT %v,%v", req.ItemId, req.PageSize*req.PageOffset, req.PageSize)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	defer rows.Close()

	historys := make([]context.ItemTransferHistory, 0)
	for rows.Next() {
		history := context.ItemTransferHistory{}
		if err := rows.Scan(&history.Idx, &history.ItemId, &history.FromAddr, &history.ToAddr, &history.TokenId, &history.State, &history.Hash, &history.Timestamp); err != nil {
			log.Error(err)
		}
		historys = append(historys, history)
	}

	totalCount, err := o.GetTotalHistoryTransferItemSize(req.ItemId)

	return historys, totalCount, err
}

func (o *DB) GetTotalHistoryTransferItemSize(ItemId int64) (int64, error) {
	sqlQuery := fmt.Sprintf("SELECT COUNT(*) as count FROM ipblock.items_history_transfer WHERE item_id=%v", ItemId)
	rows, err := o.Mysql.Query(sqlQuery)

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

func (o *DB) GetHistoryTransferMe(req *context.GetHistoryTransferMe) ([]context.ItemTransferHistory, int64, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM ipblock.items_history_transfer WHERE from_addr='%v' OR to_addr='%v' ORDER BY idx DESC LIMIT %v,%v", req.WalletAddress, req.WalletAddress, req.PageSize*req.PageOffset, req.PageSize)
	rows, err := o.Mysql.Query(sqlQuery)
	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	defer rows.Close()

	historys := make([]context.ItemTransferHistory, 0)
	for rows.Next() {
		history := context.ItemTransferHistory{}
		if err := rows.Scan(&history.Idx, &history.ItemId, &history.FromAddr, &history.ToAddr, &history.TokenId, &history.State, &history.Hash, &history.Timestamp); err != nil {
			log.Error(err)
		}
		historys = append(historys, history)
	}

	totalCount, err := o.GetTotalHistoryTransferMeSize(req.WalletAddress)

	return historys, totalCount, err
}

func (o *DB) GetTotalHistoryTransferMeSize(walletAddr string) (int64, error) {
	sqlQuery := fmt.Sprintf("SELECT COUNT(*) as count FROM ipblock.items_history_transfer WHERE from_addr='%v' OR to_addr='%v'", walletAddr, walletAddr)
	rows, err := o.Mysql.Query(sqlQuery)

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

func (o *DB) InsertHistory(TxHash, FromAddr, ToAddr string, TokenID int64, State string) (int64, error) {
	item, err := o.GetItemByTokenId(TokenID)
	if err != nil || item.ItemId == 0 {
		log.Error(err)
		return -1, err
	}

	sqlQuery := fmt.Sprintf("INSERT INTO ipblock.items_history_transfer(item_id, from_addr, to_addr, token_id, state, hash, timestamp" +
		") VALUES (?,?,?,?,?,?,?)")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, item.ItemId, FromAddr, ToAddr, TokenID, State, TxHash, datetime.GetTS2MilliSec())
	if err != nil {
		log.Error(err)
		return -1, err
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	log.Debug("insert id:", insertId)
	return insertId, nil
}

func (o *DB) UpdateHistory(TxHash, FromAddr, ToAddr string, TokenID int64) {

}
