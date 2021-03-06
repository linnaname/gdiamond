package service

import (
	"errors"
	"fmt"
	"gdiamond/server/internal/common"
	"gdiamond/server/internal/model"
)

//FetchPage util for fetch data by page
//page data need total count and total page so it  count rows sql, args will pass to count row  sql and fetch row sql
func FetchPage(sqlCountRows, sqlFetchRows string, pageNo, pageSize int, args ...interface{}) (*model.Page, error) {
	if pageSize <= 0 {
		return nil, errors.New("pageSize can't greater than 0")
	}

	stm, err := common.GDBConn.Prepare(sqlCountRows)
	defer stm.Close()
	if err != nil {
		return nil, err
	}
	rowCount := 0
	err = stm.QueryRow(args...).Scan(&rowCount)
	if err != nil {
		return nil, err
	}

	pageCount := rowCount / pageSize
	if rowCount > pageSize*pageCount {
		pageCount++
	}

	if pageNo > pageCount {
		return nil, errors.New("pageNo can't greater than pageCount")
	}

	startRow := (pageNo - 1) * pageSize
	selectSQL := fmt.Sprintf("%v limit %v , %v", sqlFetchRows, startRow, pageSize)
	sstm, err := common.GDBConn.Prepare(selectSQL)
	if err != nil {
		return nil, err
	}
	rows, err := sstm.Query(args...)
	if err != nil {
		return nil, err
	}
	pageItems := make([]interface{}, 0, pageSize)
	for rows.Next() {
		configInfo := &model.ConfigInfo{}
		err := rows.Scan(&configInfo.ID, &configInfo.DataID, &configInfo.Group, &configInfo.Content, &configInfo.MD5, &configInfo.LastModified)
		if err != nil {
			return nil, err
		}
		pageItems = append(pageItems, configInfo)
	}

	page := &model.Page{PageNO: pageNo, PageAvailable: pageCount, TotalCount: rowCount, PageItems: pageItems}
	return page, nil
}
