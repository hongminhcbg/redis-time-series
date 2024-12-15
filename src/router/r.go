package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hongminhcbg/velocity-rule/src/service"
)

func InitGin(e *gin.Engine, s *service.Service) {
	e.POST("/api/v1/users", s.CreateUser)
	e.POST("/api/v1/velocity-rule-in", s.VelocityInput)
	e.POST("/api/v1/run-realtime-rule", s.RunRealtimeRule)
}
