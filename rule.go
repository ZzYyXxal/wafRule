package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
)

const dbName string = "waf"
const mysqlUsername string = "root"
const mysqlPassword string = "12345678"
const mysqlServer string = "127.0.0.1:3306"

//create table rule ( id varchar(255), name varchar(255), rex varchar(255), timestamp datetime(6), type varchar(255), risk_level varchar(255));

type RuleInfo struct {
	ID        string `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Rex       string `db:"rex" json:"rex"`
	Timestamp string `db:"timestamp" json:"timestamp"`
	RuleType  string `db:"type" json:"type"`
	RiskLevel string `db:"risk_level" json:"risk_level"`
}

func findAll() []RuleInfo {
	db, err := sqlConn()
	defer db.Close()

	if err != nil {
		fmt.Println("connect error:", err)
		return nil
	}

	var ss []RuleInfo

	rows, err := db.Query("SELECT * FROM rule")
	defer rows.Close()

	if err != nil {
		fmt.Println("find err:", err)
	}

	for rows.Next() {
		var s RuleInfo
		err := rows.Scan(&s.ID, &s.Name, &s.Rex, &s.Timestamp, &s.RuleType, &s.RiskLevel)
		if err != nil {
			fmt.Println("rows Scan error:", err)
			return ss
		}
		//fmt.Println(s.id, s.name, s.rex, s.timestamp, s.ruleType, s.riskLevel)
		ss = append(ss, s)
	}
	return ss
}

func insert(name, rex, ruleType, riskLevel string) (bool, string) {
	var levels []string = []string{"low", "mid", "high"}
	if !inArray(riskLevel, levels) {
		return false, "参数错误"
	}

	db, err := sqlConn()
	defer db.Close()

	if err != nil {
		fmt.Println("connect error:", err)
		return false, "连接数据库失败"
	}

	uid := uuid.NewV4().String()
	fmt.Println("UUID:", uid)

	timestamp := time.Now().Format("2006-01-02 15:04:05.000000")

	result, err := db.Exec("INSERT INTO rule VALUES (?, ?, ?, ?, ?, ?)", uid, name, rex, timestamp, ruleType, riskLevel)

	if err != nil {
		fmt.Println("insert error:", err)
		return false, "插入失败"
	}

	line, _ := result.RowsAffected()
	fmt.Println(line)
	// // return line
	return true, "OK"
}

func deleteByID(id string) (bool, string) {
	db, err := sqlConn()
	defer db.Close()

	if err != nil {
		fmt.Println("connect error:", err)
		return false, "连接数据库失败"
	}

	result, err := db.Exec("DELETE FROM rule WHERE id = ?", id)

	if err != nil {
		fmt.Println("delete error:", err)
		return false, "删除失败"
	}

	line, _ := result.RowsAffected()
	fmt.Println(line)
	return true, strconv.FormatInt(line, 10) + " lines deleted."
}

func update(id, newName, newRex, newRuleType, newRiskLevel string) (bool, string) {
	db, err := sqlConn()
	defer db.Close()

	if err != nil {
		fmt.Println("connect error:", err)
		return false, "连接数据库失败"
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000000")
	result, err := db.Exec("UPDATE rule SET `name` = ?, `rex` = ?, `type` = ?, `risk_level` = ?, `timestamp` = ? WHERE id = ?",
		newName, newRex, newRuleType, newRiskLevel, timestamp, id)
	if err != nil {
		fmt.Println("update error:", err)
		return false, "更新失败"
	}

	line, _ := result.RowsAffected()
	fmt.Println(line)
	return true, strconv.FormatInt(line, 10) + " lines deleted."
}

func sqlConn() (*sql.DB, error) {
	mysqlConnect := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", mysqlUsername, mysqlPassword, mysqlServer, dbName)
	db, err := sql.Open("mysql", mysqlConnect)

	return db, err
}

func inArray(s string, ss []string) bool {
	for _, e := range ss {
		if s == e {
			return true
		}
	}
	return false
}
