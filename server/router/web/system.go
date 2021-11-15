package web

import (
	"github.com/gin-gonic/gin"
	"goweb-gin-demo/api"
	"goweb-gin-demo/middleware"
)

type SysRouter struct {
}

func (s *SysRouter) InitSystemRouter(Router *gin.RouterGroup) {
	sysRouter := Router.Group("system").Use(middleware.OperationRecord())
	var systemApi = api.ApiGroupApp.ApiGroup.SystemApi
	{
		sysRouter.POST("getSystemConfig", systemApi.GetSystemConfig) // 获取配置文件内容
		sysRouter.POST("setSystemConfig", systemApi.SetSystemConfig) // 设置配置文件内容
		sysRouter.POST("getServerInfo", systemApi.GetServerInfo)     // 获取服务器信息
		sysRouter.POST("reloadSystem", systemApi.ReloadSystem)       // 重启服务
	}
}
