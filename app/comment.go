package app

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "github.com/zhanben/go_site/tool/db"
)

type Comment struct{
    DbConn *db.DbConn
    Logger *zap.SugaredLogger
}

func NewComment(dbConn *db.DbConn, logger *zap.SugaredLogger) (*Comment, error){
    comment := &Comment{
        DbConn: dbConn,
        Logger: logger,
    }
    return comment, nil
}

func (ct *Comment) initRouter(r *gin.RouterGroup) {
    //在此添加接口
    r.GET("/comments", ct.getAllComment) //根据所有的评论信息
}

func (ct *Comment) getAllComment(c *gin.Context) {
    //构建返回结构体
    res := map[string]interface{}{
        "Action":  "GetAllUserResponse",
        "RetCode": 0,
    }
    ct.Logger.Info("get all the comment successful!")
    c.JSON(http.StatusOK, res)
}
