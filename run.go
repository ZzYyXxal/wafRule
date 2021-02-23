package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//查询数据库中所有数据
	router.GET("/app/v1/sqlapi/findAll", func(c *gin.Context) {
		ss := findAll()
		if len(ss) == 0 {
			c.String(http.StatusInternalServerError, "Can't find data")
		}
		c.JSON(http.StatusOK, ss)
	})

	//插入数据
	router.POST("/app/v1/sqlapi/insert", func(c *gin.Context) {
		json := RuleInfo{}
		err := c.BindJSON(&json)
		if err != nil {
			c.String(http.StatusBadRequest, "Error!")
		}

		flg, s := insert(json.Name, json.Rex, json.RuleType, json.RiskLevel)
		if flg == false {
			c.String(http.StatusInternalServerError, s)
		} else {
			c.String(http.StatusOK, s)
		}
	})

	//删除数据
	router.POST("/app/v1/sqlapi/delete", func(c *gin.Context) {
		json := RuleInfo{}
		err := c.BindJSON(&json)
		if err != nil {
			c.String(http.StatusBadRequest, "Error!")
		}

		flg, s := deleteByID(json.ID)
		if flg == false {
			c.String(http.StatusInternalServerError, s)
		} else {
			c.String(http.StatusOK, s)
		}
	})

	//更新数据
	router.POST("/app/v1/sqlapi/update", func(c *gin.Context) {
		json := RuleInfo{}
		err := c.BindJSON(&json)
		if err != nil {
			c.String(http.StatusBadRequest, "Error!")
		}

		flg, s := update(json.ID, json.Name, json.Rex, json.RuleType, json.RiskLevel)
		if flg == false {
			c.String(http.StatusInternalServerError, s)
		} else {
			c.String(http.StatusOK, s)
		}
	})

	router.Run(":8000")
}
