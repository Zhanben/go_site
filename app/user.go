package app

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"fmt"
	"github.com/zhanben/go_site/tool/db"
)

type User struct {
	db *db.DbConn
	Logger *zap.Logger
}

func NewUser(dbConn *db.DbConn, logger *zap.Logger) (*User, error) {
	user := &User{
		db: dbConn,
		Logger: logger,
	}
	return user, nil
}

func (u *User) initRouter(r *gin.RouterGroup) {
	//在此添加接口
	r.GET("/users", u.getAllUsers)      //根据账号获取所有用户信息
	r.GET("/users/:name", u.getOneUser) //根据用户名获取用户详细信息
}

func (u *User) getAllUsers(c *gin.Context) {
	//构建返回结构体
	res := map[string]interface{}{
		"Action":  "GetAllUserResponse",
		"RetCode": 0,
	}

	sql := "select * from user limit 10"
	result, err := u.db.Select(sql)
	if err != nil {
		u.Logger.Error("get user info from db error!")
		abortWithError(u.Logger, c, err)
	}
	res["UserInfo"] = result
	u.Logger.Info("get all the user info successful!")
	c.JSON(http.StatusOK, res)
}

func (u *User) getOneUser(c *gin.Context) {
	//构建返回结构体
	res := map[string]interface{}{
		"Action":  "GetOneUserResponse",
		"RetCode": 0,
	}
	username, ok :=c.Params.Get("name")
	if !ok {
        u.Logger.Error("parameter name must be fixed!")
    }
	u.Logger.Error(fmt.Sprintf("get user name from url:%s",username))

    sql := "select * from user where username=?"
    result, err := u.db.Select(sql,username)
    if err != nil {
        u.Logger.Error("get user info from db error!")
        res["RetCode"]= "-1"
        res["Error"] = "user not exist!"
    }else{
        res["UserInfo"] = result
	    u.Logger.Info("get one user info successful!")
    }
	c.JSON(http.StatusOK, res)
}
