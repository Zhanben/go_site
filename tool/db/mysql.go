package db

import (
    "database/sql"
    "sync"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

const DEFAULT_MAX_OPEN_CONNS = 500
const DEFAULT_MAX_IDLE_CONNS = 500
const DEFAULT_MAX_LIFE_TIME = 3600

type MysqlConn struct {
    Db *sql.DB
}

var (
    mysqlConnPoolMu sync.Mutex
    mysqlConnPool   = make(map[string]*MysqlConn)
)

func initMysqlConn(connStr string, maxOpenConns int, maxIdleConns int, maxLifetime int) (mysqlConn *MysqlConn, err error) {
    db, err := sql.Open("mysql", connStr)
    if err != nil {
        return
    }

    db.SetMaxOpenConns(maxOpenConns)
    db.SetMaxIdleConns(maxIdleConns)
    db.SetConnMaxLifetime(time.Second * time.Duration(maxLifetime))

    // 实际连接一次, 判断连接是否成功
    err = db.Ping()
    if err != nil {
        return
    }

    mysqlConn = &MysqlConn{
        Db: db,
    }
    mysqlConnPoolMu.Lock()
    mysqlConnPool[connStr] = mysqlConn
    mysqlConnPoolMu.Unlock()
    return
}

// 获取mysql连接实例
func GetMysqlInstance(connStr string) (mysqlConn *MysqlConn, err error) {
    mysqlConnPoolMu.Lock()
    conn, ok := mysqlConnPool[connStr]
    mysqlConnPoolMu.Unlock()
    if ok {
        mysqlConn = conn
        return
    }
    return initMysqlConn(connStr, DEFAULT_MAX_OPEN_CONNS, DEFAULT_MAX_IDLE_CONNS, DEFAULT_MAX_LIFE_TIME)
}

// 获取mysql连接实例(带连接池信息)
func GetMysqlInstance2(connStr string, maxOpenConns int, maxIdleConns int, maxLifetime int) (mysqlConn *MysqlConn, err error) {
    mysqlConnPoolMu.Lock()
    conn, ok := mysqlConnPool[connStr]
    mysqlConnPoolMu.Unlock()
    if ok {
        mysqlConn = conn
        return
    }
    return initMysqlConn(connStr, maxOpenConns, maxIdleConns, maxLifetime)
}


// 新增
func (c *MysqlConn) Insert(sql string, args ...interface{}) (lastInsertId int64, err error) {
    res, err := c.Db.Exec(sql, args...)
    if err != nil {
        return
    }
    return res.LastInsertId()
}

// 修改
func (c *MysqlConn) Update(sql string, args ...interface{}) (rowsAffected int64, err error) {
    res, err := c.Db.Exec(sql, args...)
    if err != nil {
        return
    }
    return res.RowsAffected()
}

// 删除
func (c *MysqlConn) Delete(sql string, args ...interface{}) (rowsAffected int64, err error) {
    res, err := c.Db.Exec(sql, args...)
    if err != nil {
        return
    }
    return res.RowsAffected()
}

// 查询
func (c *MysqlConn) Select(sql string, args ...interface{}) (results []map[string]string, err error) {
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

func (c *MysqlConn) GetSqlDB() *sql.DB {
    return c.Db
}