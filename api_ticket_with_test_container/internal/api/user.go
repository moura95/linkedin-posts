package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github/moura95/ticket-api/pkg/ginx"
	"github/moura95/ticket-api/pkg/jsonplaceholder"
	"go.uber.org/zap"
)

func (u *UserRouter) list(ctx *gin.Context) {
	u.logger.Info("List All Users")

	response, err := jsonplaceholder.ListUsers()
	if err != nil {
		u.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, ginx.ErrorResponse("Error List Users"))
		return
	}

	ctx.JSON(http.StatusOK, ginx.SuccessResponse(response))
}

func (u *UserRouter) get(ctx *gin.Context) {

	u.logger.Info("Get By ID User")

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		u.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, ginx.ErrorResponse("Bad Request, Id Invalid"))
		return
	}
	response, err := jsonplaceholder.GetUserByID(int32(id))
	if err != nil {
		u.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, ginx.ErrorResponse("Error Get User"))
		return
	}
	if response == nil {
		ctx.JSON(http.StatusNotFound, ginx.SuccessResponse(response))
		return
	}

	ctx.JSON(http.StatusOK, ginx.SuccessResponse(response))
}

type IUser interface {
	SetupUserRoute(routers *gin.RouterGroup)
}

type UserRouter struct {
	logger *zap.SugaredLogger
}

func NewUserRouter(log *zap.SugaredLogger) *UserRouter {
	return &UserRouter{
		logger: log,
	}
}

func (t *UserRouter) SetupUserRoute(routers *gin.RouterGroup) {
	routers.GET("/users", t.list)
	routers.GET("/users/:id", t.get)
}
