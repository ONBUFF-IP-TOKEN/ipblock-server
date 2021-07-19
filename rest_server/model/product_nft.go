package model

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
)

func (o *DB) InsertProductNFT(product *context.ProductInfo, quantityIndex, state int64, createHash, ownerWalletAddr, uri string) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO  product_nft(product_id, create_ts, create_hash, quantity_index, owner_wallet_address, nft_uri, state)" +
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
	sqlQuery := "UPDATE  product_nft set token_id=?, state=? WHERE create_hash=?"

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

func (o *DB) GetNftList(pageInfo *context.NftList) ([]context.NftInfo, int64, error) {
	var sqlQuery string
	if pageInfo.ProductId == 0 {
		sqlQuery = fmt.Sprintf("SELECT * FROM  product_nft ORDER BY product_id DESC LIMIT %v,%v", pageInfo.PageSize*pageInfo.PageOffset, pageInfo.PageSize)
	} else {
		sqlQuery = fmt.Sprintf("SELECT * FROM  product_nft WHERE product_id=%v ORDER BY product_id DESC LIMIT %v,%v", pageInfo.ProductId, pageInfo.PageSize*pageInfo.PageOffset, pageInfo.PageSize)
	}
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	defer rows.Close()

	var tokenId, state sql.NullInt64
	var ownerWalletAddr sql.NullString
	nfts := make([]context.NftInfo, 0)
	for rows.Next() {
		nft := context.NftInfo{}
		if err := rows.Scan(&nft.ProductId, &nft.CreateTs, &nft.CreateHash, &tokenId, &nft.QuantityIndex, &ownerWalletAddr, &nft.NftUri, &state); err != nil {
			log.Error(err)
		}
		nft.TokenId = tokenId.Int64
		nft.OwnerWalletAddr = ownerWalletAddr.String
		nft.State = state.Int64
		nfts = append(nfts, nft)
	}

	totalCount, err := o.GetTotalProductSize()

	return nfts, totalCount, err
}

func (o *DB) GetNftListByProductId(productId int64) ([]context.NftInfo, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM  product_nft WHERE product_id=%v ORDER BY token_id DESC", productId)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	var tokenId, state sql.NullInt64
	var ownerWalletAddr sql.NullString
	nfts := make([]context.NftInfo, 0)
	for rows.Next() {
		nft := context.NftInfo{}
		if err := rows.Scan(&nft.ProductId, &nft.CreateTs, &nft.CreateHash, &tokenId, &nft.QuantityIndex, &ownerWalletAddr, &nft.NftUri, &state); err != nil {
			log.Error(err)
		}
		nft.TokenId = tokenId.Int64
		nft.OwnerWalletAddr = ownerWalletAddr.String
		nft.State = state.Int64
		nfts = append(nfts, nft)
	}

	return nfts, err
}

func (o *DB) UpdateProductNftOwner(createHash, newOwner, oldOwner string, tokenId int64) (int64, error) {
	sqlQuery := "UPDATE  product_nft set owner_wallet_address=? WHERE owner_wallet_address=? and token_id=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, newOwner, createHash, oldOwner, tokenId)
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

func (o *DB) GetTotalNftSize() (int64, error) {
	rows, err := o.Mysql.Query("SELECT COUNT(*) as count FROM  product_nft")
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
