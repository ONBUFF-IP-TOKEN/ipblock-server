package model

import (
	"github.com/ONBUFF-IP-TOKEN/basedb"
)

type DB struct {
	Mysql *basedb.Mysql
	Cache *basedb.Cache
}

var gDB *DB

func SetDB(db *basedb.Mysql, cache *basedb.Cache) {
	gDB = &DB{
		Mysql: db,
		Cache: cache,
	}
}

func GetDB() *DB {
	return gDB
}
