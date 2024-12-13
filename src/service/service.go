package service

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hongminhcbg/velocity-rule/config"
	"github.com/hongminhcbg/velocity-rule/src/models"
	"github.com/hongminhcbg/velocity-rule/src/rules"
	"github.com/hongminhcbg/velocity-rule/src/store"
	"github.com/hongminhcbg/velocity-rule/src/timeseries"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/go-redis/redis/v9"
)

type Service struct {
	store   *store.UserStore
	log     logr.Logger
	redists timeseries.TimeSeries
}

func NewService(cfg *config.Config, store *store.UserStore, r *redis.Client, log logr.Logger) *Service {
	return &Service{
		store:   store,
		log:     log,
		redists: timeseries.New(),
	}
}

func (s *Service) createNewUser(ctx *gin.Context, req *models.CreateUserRequest) {
	if req.ReqId == "" {
		s.log.Info("req_id id empty, auto generate")
		req.ReqId = fmt.Sprint(time.Now().UnixMilli())
	}

	r := models.User{
		Name:      req.Name,
		ReqId:     req.ReqId,
		RetryTime: 0,
		Status:    models.USER_INIT,
	}

	err := s.store.Save(ctx.Request.Context(), &r)
	if err != nil {
		s.log.Error(err, "save to record error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	s.log.Info("create new user success")
	ctx.JSON(http.StatusOK, r)
}

func (s *Service) handleExistedReqId(ctx *gin.Context, u *models.User, req *models.CreateUserRequest) {
	ctx.JSON(http.StatusOK, u)
}

func (s *Service) CreateUser(ctx *gin.Context) {
	var req models.CreateUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.log.Error(err, "json unmarshal error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if req.ReqId == "" {
		s.log.Info("Will create new user")
		s.createNewUser(ctx, &req)
		return
	}

	u, _ := s.store.GetByReqId(ctx.Request.Context(), req.ReqId)
	if u == nil {
		s.createNewUser(ctx, &req)
		return
	}

	s.handleExistedReqId(ctx, u, &req)
}

func (s *Service) VelocityInput(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		s.log.Error(err, "get body error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return

	}

	rMsg := &rules.RuleMsg{
		Input:  b,
		UserId: "1",
	}
	err = rMsg.Parse()
	if err != nil {
		s.log.Error(err, "parse json error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = rules.MainEngine.Execute(ctx.Request.Context(), rMsg, s.log)
	if err != nil {
		s.log.Error(err, "execute rule error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	s.log.Info("execute rule success", "output", rMsg.VelocityOuput)
	for _, v := range rMsg.VelocityOuput {
		err = s.redists.AddAutoTs(fmt.Sprintf("%s%s", v.UserId, v.DataKey), v.Data)
		if err != nil {
			s.log.Error(err, "save to redis error")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "output": rMsg.VelocityOuput})
}
