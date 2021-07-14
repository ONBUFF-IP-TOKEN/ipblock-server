package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
)

func (o *DB) InsertProduct(product *context.ProductInfo) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO ipblock.product(product_title, product_thumbnail_url, product_price, product_type, token_type," +
		"create_ts, creator, description, content, quantity_total, quantity_remaining, state) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, product.Title, product.Thumbnail, product.Price, product.ProductType, product.TokenType,
		product.CreateTs, product.Creator, product.Desc, product.Content, product.QuantityTotal, product.QuantityRemaining, product.State)
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

func (o *DB) UpdateProductState(product *context.ProductUpdateState) (int64, error) {
	sqlQuery := "UPDATE ipblock.product set state=? WHERE product_id=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, product.State, product.ProductId)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return cnt, nil
}
