package model

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
)

func (o *DB) InsertItem(item *context.RegisterItem) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO ipblock.items_base(wallet_address, title, token_type, thumbnail_url, token_price," +
		"expire_date, register_date, creator, description, owner_wallet_address, owner, create_hash) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, item.WalletAddr, item.Title, item.TokenType, item.Thumbnail, item.TokenPrice,
		item.ExpireDate, datetime.GetTS2MilliSec(), item.Creator, item.Description, item.WalletAddr, item.Creator, item.CreateHash)
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

func (o *DB) DeleteItem(itemId int64) (bool, error) {
	sqlQuery := "DELETE FROM ipblock.items_base WHERE idx=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, itemId)
	if err != nil {
		log.Error(err)
		return false, err
	}
	cnt, err := result.RowsAffected()
	if cnt == 0 {
		log.Error(err)
		return false, err
	}

	return true, nil
}

func (o *DB) DeleteItemByTokenId(TokenId int64) (bool, error) {
	sqlQuery := "DELETE FROM ipblock.items_base WHERE token_id=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, TokenId)
	if err != nil {
		log.Error(err)
		return false, err
	}
	cnt, err := result.RowsAffected()
	if cnt == 0 {
		log.Error(err)
		return false, err
	}

	return true, nil
}

func (o *DB) GetItem(itemId int64) (*context.ItemInfo, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM ipblock.items_base WHERE idx=%v", itemId)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	var desc, createHash, content sql.NullString
	var tokenId sql.NullInt64
	item := &context.ItemInfo{}
	if rows.Next() {
		if err := rows.Scan(&item.ItemId, &item.WalletAddr, &item.Title, &item.TokenType, &item.Thumbnail, &item.TokenPrice,
			&item.ExpireDate, &item.RegisterDate, &item.Creator, &desc, &item.OwnerWalletAddr, &item.Owner, &tokenId, &createHash, &content); err != nil {
			log.Error(err)
		}
		item.Description = desc.String
		item.TokenId = tokenId.Int64
		item.CreateHash = createHash.String
		item.Content = content.String
	}
	return item, nil
}

func (o *DB) GetItemByTokenId(TokenId int64) (*context.ItemInfo, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM ipblock.items_base WHERE token_id=%v", TokenId)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	var desc, createHash, content sql.NullString
	var tokenId sql.NullInt64

	item := &context.ItemInfo{}
	if rows.Next() {
		if err := rows.Scan(&item.ItemId, &item.WalletAddr, &item.Title, &item.TokenType, &item.Thumbnail, &item.TokenPrice,
			&item.ExpireDate, &item.RegisterDate, &item.Creator, &desc, &item.OwnerWalletAddr, &item.Owner, &tokenId, &createHash, &content); err != nil {
			log.Error(err)
		}
		item.Description = desc.String
		item.TokenId = tokenId.Int64
		item.CreateHash = createHash.String
		item.Content = content.String

	}
	return item, nil
}

func (o *DB) GetItemList(pageInfo *context.GetItemList) ([]context.ItemInfo, int64, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM ipblock.items_base ORDER BY idx DESC LIMIT %v,%v", pageInfo.PageSize*pageInfo.PageOffset, pageInfo.PageSize)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	defer rows.Close()

	var desc, createHash, content sql.NullString
	var tokenId sql.NullInt64
	items := make([]context.ItemInfo, 0)
	for rows.Next() {
		item := context.ItemInfo{}
		if err := rows.Scan(&item.ItemId, &item.WalletAddr, &item.Title, &item.TokenType, &item.Thumbnail, &item.TokenPrice,
			&item.ExpireDate, &item.RegisterDate, &item.Creator, &desc, &item.OwnerWalletAddr, &item.Owner, &tokenId, &createHash, &content); err != nil {
			log.Error(err)
		}
		item.Description = desc.String
		item.TokenId = tokenId.Int64
		item.CreateHash = createHash.String
		item.Content = content.String
		items = append(items, item)
	}

	totalCount, err := o.GetTotalItemSize()

	return items, totalCount, err
}

func (o *DB) GetTotalItemSize() (int64, error) {
	rows, err := o.Mysql.Query("SELECT COUNT(*) as count FROM ipblock.items_base")
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

func (o *DB) UpdateTokenID(txHash string, tokenID int64) error {
	sqlQuery := "UPDATE ipblock.items_base set token_id=? WHERE create_hash=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, tokenID, txHash)
	if err != nil {
		log.Error(err)
		return err
	}
	cnt, err := result.RowsAffected()
	if cnt == 0 {
		err = errors.New("RowsAffected none")
		log.Error(err)
		return err
	}
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (o *DB) UpdateTransfer(txHash string, fromAddr, toAddr string, tokenId int64) error {
	sqlQuery := "UPDATE ipblock.items_base set owner_wallet_address=? WHERE token_id=? and owner_wallet_address=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, toAddr, tokenId, fromAddr)
	if err != nil {
		log.Error(err)
		return err
	}
	cnt, err := result.RowsAffected()
	if cnt == 0 {
		log.Error(err)
		return err
	}

	return nil
}
