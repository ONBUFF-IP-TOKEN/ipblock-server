package model

import (
	"errors"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
)

func (o *DB) InsertOrder(order *context.OrderInfo) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO orders(order_date, purchase_tx_hash, order_state, product_id, product_price, quantity_index, quantity_total," +
		"customer_wallet_address, customer_email, token_id) VALUES (?,?,?,?,?,?,?,?,?,?)")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, order.Date, order.PurchaseTxHash, order.State, order.ProductId, order.Price, order.QuantityIndex,
		order.QuantityTotal, order.CustomerWalletAddr, order.CustomerEmail, order.TokenId)
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

func (o *DB) UpdateOrderState(tokenId, orderState int64) (int64, error) {
	sqlQuery := "UPDATE orders set order_state=? WHERE token_id=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, orderState, tokenId)
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
