package model

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

func (o *DB) InsertNotice(notices *context_auc.NoticeRegister) error {
	for _, notice := range notices.Notices {
		sqlQuery := fmt.Sprintf("INSERT INTO auc_notices (title, description, urls) VALUES (?,?,?)")

		title, _ := json.Marshal(notice.Title)
		desc, _ := json.Marshal(notice.Desc)
		urls, _ := json.Marshal(notice.Urls)

		result, err := o.Mysql.PrepareAndExec(sqlQuery, string(title), string(desc), string(urls))

		if err != nil {
			log.Error(err)
			return err
		}
		insertId, err := result.LastInsertId()
		if err != nil {
			log.Error(err)
			return err
		}
		log.Debug("InsertAucAuction id:", insertId)
	}

	// notice list cache 전체 삭제
	o.DeleteNoticeList()
	return nil
}

func (o *DB) GetNotice(pageInfo *context_auc.NoticeList) (*[]context_auc.Notice, int64, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM auc_notices ORDER BY id DESC LIMIT %v,%v", pageInfo.PageSize*pageInfo.PageOffset, pageInfo.PageSize)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	defer rows.Close()

	notices := []context_auc.Notice{}
	for rows.Next() {
		var err error
		notice, err := o.ScanNotice(rows)
		if err != nil {
			log.Error(err)
			continue
		}
		notices = append(notices, *notice)
	}

	count, _ := o.GetTotalNoticeSize()

	return &notices, count, err
}

func (o *DB) DeleteNotice(noticeId int64) (bool, error) {
	var (
		sqlQuery string
		result   sql.Result
		err      error
	)

	if noticeId < 0 {
		//전체 삭제
		sqlQuery = "DELETE FROM auc_notices"
		result, err = o.Mysql.PrepareAndExec(sqlQuery)
	} else {
		sqlQuery = "DELETE FROM auc_notices WHERE id=?"
		result, err = o.Mysql.PrepareAndExec(sqlQuery, noticeId)
	}

	if err != nil {
		log.Error(err)
		return false, err
	}
	cnt, err := result.RowsAffected()
	if cnt == 0 {
		log.Error(err)
		return false, err
	}

	// notice list cache 전체 삭제
	o.DeleteNoticeList()

	return true, nil
}

func (o *DB) UpdateNotice(notices *context_auc.NoticeUpdate) error {
	for _, notice := range notices.Notices {
		sqlQuery := fmt.Sprintf("UPDATE auc_notices set title=?, description=?, urls=? WHERE id=?")

		title, _ := json.Marshal(notice.Title)
		desc, _ := json.Marshal(notice.Desc)
		urls, _ := json.Marshal(notice.Urls)

		result, err := o.Mysql.PrepareAndExec(sqlQuery, string(title), string(desc), string(urls), notice.Id)

		if err != nil {
			log.Error(err)
			return err
		}
		insertId, err := result.RowsAffected()
		if err != nil {
			log.Error(err)
			return err
		}
		log.Debug("UpdateNotice id:", insertId)
	}

	// notice list cache 전체 삭제
	o.DeleteNoticeList()
	return nil
}

func (o *DB) ScanNotice(rows *sql.Rows) (*context_auc.Notice, error) {
	var title, desc, urls sql.NullString

	notice := &context_auc.Notice{}
	if err := rows.Scan(&notice.Id, &title, &desc, &urls); err != nil {
		log.Error("ScanNotice error: ", err)
		return nil, err
	}

	aTitle := context_auc.Localization{}
	json.Unmarshal([]byte(title.String), &aTitle)
	notice.Title = aTitle

	aDesc := context_auc.Localization{}
	json.Unmarshal([]byte(desc.String), &aDesc)
	notice.Desc = aDesc

	aUrls := []string{}
	json.Unmarshal([]byte(urls.String), &aUrls)
	notice.Urls = aUrls

	return notice, nil
}

func (o *DB) GetTotalNoticeSize() (int64, error) {
	var dataCount int64
	if err := o.Mysql.QueryRow("SELECT COUNT(*) as count FROM auc_notices", &dataCount); err != nil {
		log.Error("GetTotalNoticeSize : ", err)
		return 0, err
	}
	return dataCount, nil
}
