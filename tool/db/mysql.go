package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"github.com/spf13/viper"
	"fmt"
)

type DbConn struct {
	Db *sql.DB
}


func InitDbConn()(*DbConn, error) {
	dbConn := viper.GetString("Database.DbConn")
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		return nil, fmt.Errorf("get db connection error")
	}
	maxOpenConns := viper.GetInt("Database.MaxOpenConn")
	maxIdleConns := viper.GetInt("Database.MaxIdleConn")
	maxLifetime := viper.GetInt("Database.MaxLifetime")

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Second * time.Duration(maxLifetime))

	// 实际连接一次, 判断连接是否成功
	err = db.Ping()
	if err != nil {
		return nil,fmt.Errorf("connection db error")
	}else{
		dbconn := &DbConn{Db:db}
		return dbconn,nil
	}
}

func (c *DbConn) Insert(sql string, args ...interface{}) (lastInsertId int64, err error) {
	res, err := c.Db.Exec(sql, args...)
	if err != nil {
		return
	}
	return res.LastInsertId()
}

func (c *DbConn) Update(sql string, args ...interface{}) (rowsAffected int64, err error) {
	res, err := c.Db.Exec(sql, args...)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

func (c *DbConn) Delete(sql string, args ...interface{}) (rowsAffected int64, err error) {
	res, err := c.Db.Exec(sql, args...)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

func (c *DbConn) Select(sql string, args ...interface{}) (results []map[string]string, err error) {
	rows, err := c.Db.Query(sql, args...)
	if err != nil {
		return
	}
	// 关闭数据集
	defer rows.Close()
	// 列信息
	columns, _ := rows.Columns()
	// 列值
	values := make([][]byte, len(columns))
	// 扫描器
	scans := make([]interface{}, len(columns))
	for i := range values {
		scans[i] = &values[i]
	}
	results = make([]map[string]string, 0)

	for rows.Next() {
		if err = rows.Scan(scans...); err != nil {
			return
		}
		row := make(map[string]string)
		for k, v := range values {
			key := columns[k]
			row[key] = string(v)
		}
		results = append(results, row)
	}
	return
}
