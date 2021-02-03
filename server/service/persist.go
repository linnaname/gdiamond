package service

import (
	"database/sql"
	"errors"
	"gdiamond/server/common"
	"gdiamond/server/model"
	"strings"
	"time"
)

func addConfigInfo(config *model.ConfigInfo) error {
	stm, err := common.GDBConn.Prepare("INSERT INTO config_info (data_id,group_id,content,md5,gmt_create,gmt_modified) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	timestamp := time.Now()
	_, eErr := stm.Exec(config.DataID, config.Group, config.Content, config.MD5, timestamp, timestamp)
	if eErr != nil {
		return eErr
	}
	return nil
}

func updateConfigInfo(config *model.ConfigInfo) error {
	stm, err := common.GDBConn.Prepare("UPDATE config_info SET content=?,md5=?,gmt_modified=? WHERE data_id=? AND group_id=?")
	if err != nil {
		return err
	}
	_, eErr := stm.Exec(config.Content, config.MD5, config.LastModified, config.DataID, config.Group)
	if eErr != nil {
		return eErr
	}
	return nil
}

func findConfigInfo(dataID, group string) (*model.ConfigInfo, error) {
	stm, err := common.GDBConn.Prepare("select id,data_id,group_id,content,md5,gmt_modified from config_info where data_id=? and group_id=?")
	if err != nil {
		return nil, err
	}
	configInfo := &model.ConfigInfo{}
	err = stm.QueryRow(dataID, group).Scan(&configInfo.ID, &configInfo.DataID, &configInfo.Group, &configInfo.Content, &configInfo.MD5, &configInfo.LastModified)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return configInfo, nil
}

func findConfigInfoByID(id int64) (*model.ConfigInfo, error) {
	stm, err := common.GDBConn.Prepare("select id,data_id,group_id,content,md5,gmt_modified from config_info where id=?")
	if err != nil {
		return nil, err
	}
	configInfo := &model.ConfigInfo{}
	err = stm.QueryRow(id).Scan(&configInfo.ID, &configInfo.DataID, &configInfo.Group, &configInfo.Content, &configInfo.MD5, &configInfo.LastModified)
	if err != nil {
		return nil, err
	}
	return configInfo, nil
}

func findConfigInfoByDataID(pageNo, pageSize int, dataID string) (*model.Page, error) {
	page, err := FetchPage("select count(id) from config_info where data_id=?",
		"select id,data_id,group_id,content,md5,gmt_modified from config_info where data_id=?",
		pageNo, pageSize, dataID)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func findConfigInfoByGroup(pageNo, pageSize int, group string) (*model.Page, error) {
	page, err := FetchPage("select count(id) from config_info where group_id=?",
		"select id,data_id,group_id,content,md5,gmt_modified from config_info where group_id=?",
		pageNo, pageSize, group)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func findAllConfigInfo(pageNo, pageSize int) (*model.Page, error) {
	page, err := FetchPage("select count(id) from config_info",
		"select id,data_id,group_id,content,md5,gmt_modified from config_info order by id ",
		pageNo, pageSize)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func findAllConfigLike(pageNo, pageSize int, dataID, group string) (*model.Page, error) {
	if dataID == "" && group == "" {
		return findAllConfigInfo(pageNo, pageSize)
	}

	sqlCountRows := "select count(id) from config_info where "
	sqlFetchRows := "select id,data_id,group_id,content,md5,gmt_modified from config_info where "
	wasFirst := true
	if dataID != "" {
		sqlCountRows += "data_id like ? "
		sqlFetchRows += "data_id like ? "
		wasFirst = false
	}
	if group != "" {
		if wasFirst {
			sqlCountRows += "group_id like ? "
			sqlFetchRows += "group_id like ? "
		} else {
			sqlCountRows += "and group_id like ? "
			sqlFetchRows += "and group_id like ? "
		}
	}

	if dataID != "" && group != "" {
		return FetchPage(sqlCountRows, sqlFetchRows, pageNo, pageSize, generateLikeArgument(dataID), generateLikeArgument(group))
	} else if dataID != "" {
		return FetchPage(sqlCountRows, sqlFetchRows, pageNo, pageSize, generateLikeArgument(dataID))
	} else if group != "" {
		return FetchPage(sqlCountRows, sqlFetchRows, pageNo, pageSize, generateLikeArgument(group))
	}
	return nil, errors.New("don't know how to do")
}

func removeConfigInfo(config *model.ConfigInfo) error {
	stm, err := common.GDBConn.Prepare("DELETE FROM config_info WHERE data_id=? AND group_id=?")
	if err != nil {
		return err
	}
	_, eErr := stm.Exec(config.DataID, config.Group)
	if eErr != nil {
		return eErr
	}
	return nil
}

func generateLikeArgument(s string) string {
	if strings.Index(s, "*") >= 0 {
		return strings.ReplaceAll(s, "\\*", "%")
	}
	return "%" + s + "%"
}
