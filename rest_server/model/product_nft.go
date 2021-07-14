package model

import (
	"errors"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
)

func (o *DB) InsertProductNFT(product *context.ProductInfo, quantityIndex, state int64, createHash, ownerWalletAddr, uri string) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO ipblock.product_nft(product_id, create_ts, create_hash, quantity_index, owner_wallet_address, nft_uri, state)" +
		" VALUES (?,?,?,?,?,?,?)")
	result, err := o.Mysql.PrepareAndExec(sqlQuery, product.Id, product.CreateTs, createHash, quantityIndex, ownerWalletAddr, uri, state)
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

func (o *DB) UpdateProductNftTokenID(createHash string, tokenId int64, state int64) (int64, error) {
	sqlQuery := "UPDATE ipblock.product_nft set token_id=?, state=? WHERE create_hash=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, tokenId, state, createHash)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if cnt == 0 {
		err = errors.New("RowsAffected none")
		log.Error(err)
		return 0, err
	}
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return cnt, nil
}
