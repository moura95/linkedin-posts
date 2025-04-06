package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github/moura95/ticket-api/internal/service"
	"github/moura95/ticket-api/internal/util"
	"github/moura95/ticket-api/pkg/errors"
	"github/moura95/ticket-api/pkg/ginx"
	"github/moura95/ticket-api/pkg/jsonplaceholder"
	"go.uber.org/zap"
)

type ticketResponse struct {
	Id            int32                 `json:"id"`
	Title         string                `json:"title"`
	Description   string                `json:"description"`
	Status        string                `json:"status"`
	SeverityId    int32                 `json:"severity_id"`
	CategoryId    int32                 `json:"category_id"`
	UserID        *int32                `json:"user_id"`
	SubCategoryId *int32                `json:"subcategory_id"`
	Category      *categoryResponse     `json:"category"`
	SubCategory   *categoryResponse     `json:"subcategory"`
	User          *jsonplaceholder.User `json:"user"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
	CompletedAt   *time.Time            `json:"completed_at"`
}

type listTicketRequest struct {
	Status     string `form:"status"`
	Title      string `form:"name"`
	SeverityId int32  `form:"severity_id"`
	Limit      int32  `form:"limit"`
	Page       int32  `form:"page"`
}

type listTicketResponse struct {
	Id            int32      `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Status        string     `json:"status"`
	SeverityId    int32      `json:"severity_id"`
	CategoryId    int32      `json:"category_id"`
	UserID        *int32     `json:"user_id"`
	SubCategoryId *int32     `json:"subcategory_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CompletedAt   *time.Time `json:"completed_at"`
}

// @Summary List all Tickets
// @Description Get a list of all tickets
// @Tags tickets
// @Accept json
// @Produce json
// @Success 200 {array} listTicketResponse
// @Router /tickets [get]
func (t *TicketRouter) list(ctx *gin.Context) {
	t.logger.Info("List All Tickets")

	var filters listTicketRequest
	err := ginx.ParseQuery(ctx, &filters)
	if err != nil {
		t.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, ginx.ErrorResponse(err.Error()))
		return
	}
	tickets, err := t.service.GetAll(ctx)
	if err != nil {
		t.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, ginx.ErrorResponse(errors.FailedToList("Tickets")))
		return
	}

	var response []listTicketResponse
	for _, ticket := range tickets {
		response = append(response, listTicketResponse{
			Id:            ticket.ID,
			Title:         ticket.Title,
			Description:   ticket.Description,
			Status:        ticket.Status,
			SeverityId:    ticket.SeverityID,
			CategoryId:    ticket.CategoryID,
			UserID:        util.NullInt32ToPtr(ticket.UserID),
			SubCategoryId: util.NullInt32ToPtr(ticket.SubcategoryID),
			CreatedAt:     ticket.CreatedAt,
			UpdatedAt:     ticket.UpdatedAt,
			CompletedAt:   util.NullDateToPtr(ticket.CompletedAt),
		})

	}

	ctx.JSON(http.StatusOK, ginx.SuccessResponse(response))
}

// @Summary Get a ticket by id
// @Description Get details of a ticket by its ID
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} ticketResponse
// @Router /tickets/{id} [get]
func (t *TicketRouter) get(ctx *gin.Context) {

	t.logger.Info("Get By ID Ticket")

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		t.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, ginx.ErrorResponse("Bad Request, Id Invalid"))
		return
	}

	ticket, err := t.service.GetByID(ctx, int32(id))
	if err != nil {
		t.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, ginx.ErrorResponse(errors.FailedToGet("Ticket")))
		return
	}

	response := ticketResponse{
		Id:            ticket.ID,
		Title:         ticket.Title,
		Description:   ticket.Description,
		Status:        ticket.Status,
		SeverityId:    ticket.SeverityID,
		CategoryId:    ticket.CategoryID,
		SubCategoryId: util.NullInt32ToPtr(ticket.SubcategoryID),
		UserID:        util.NullInt32ToPtr(ticket.UserID),
		CreatedAt:     ticket.CreatedAt,
		UpdatedAt:     ticket.UpdatedAt,
		CompletedAt:   util.NullDateToPtr(ticket.CompletedAt),
	}

	// Add Category
	response.Category = &categoryResponse{
		Id:   &ticket.CategoryID,
		Name: &ticket.CategoryName,
	}
	// Add SubCategory
	response.SubCategory = &categoryResponse{
		Id:   util.NullInt32ToPtr(ticket.SubcategoryID),
		Name: util.NullStringToPtr(ticket.SubcategoryName),
	}
	// Add User
	user, err := jsonplaceholder.GetUserByID(ticket.UserID.Int32)
	if err != nil {
		t.logger.Error(err)
	}
	response.User = user

	ctx.JSON(http.StatusOK, ginx.SuccessResponse(response))
}

type createTicketRequest struct {
	Title         string `json:"title" validate:"required"`
	Description   string `json:"description" validate:"required"`
	SeverityId    int32  `json:"severity_id" validate:"gte=1,lte=4"`
	CategoryId    int32  `json:"category_id" validate:"required"`
	SubCategoryId int32  `json:"subcategory_id"`
}

// @Summary Add a new Ticket
// @Description Add a new ticket
// @Tags tickets
// @Accept json
// @Produce json
// @Param receiver body createTicketRequest true "Ticket"
// @Success 201 {object} ticketResponse
// @Failure 400 {object} object{error=string}
// @Router /tickets [post]
func (t *TicketRouter) create(ctx *gin.Context) {
	var req createTicketRequest
	t.logger.Info("Create Ticket")

	err := ginx.ParseJSON(ctx, &req)
	if err != nil {
		t.logger.Info("Bad Request %s", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if req.SeverityId == 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request: You can't create ticket with severity issue high(1)"})
		return
	}

	// Validate  struct
	if err = t.validate.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ti, err := t.service.Create(ctx, req.Title, req.Description, req.SeverityId, req.CategoryId, req.SubCategoryId)
	if err != nil {
		t.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, ginx.ErrorResponse(err.Error()))
		return
	}
	ticket, _ := t.service.GetByID(ctx, ti.ID)

	ctx.JSON(http.StatusCreated, ginx.SuccessResponse(ticket))
}

type updateTicketRequest struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	UserID        int32  `json:"user_id"`
	SeverityId    int32  `json:"severity_id"`
	CategoryId    int32  `json:"category_id"`
	SubCategoryId int32  `json:"subcategory_id"`
}

// @Summary Update a ticket
// @Description Update a ticket with the given ID
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param receiver body updateTicketRequest true "Ticket"
// @Success 204
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /tickets/{id} [patch]
func (t *TicketRouter) update(ctx *gin.Context) {
	var req updateTicketRequest

	err := ginx.ParseJSON(ctx, &req)
	if err != nil {
		t.logger.Info("Bad Request %s", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	t.logger.Info("Update Ticket")

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		t.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, ginx.ErrorResponse("Bad Request, Id Invalid"))
		return
	}

	if req.UserID > 0 {
		user, err := jsonplaceholder.GetUserByID(req.UserID)
		if err != nil {
			t.logger.Error(err)
			ctx.JSON(http.StatusBadRequest, ginx.ErrorResponse(err.Error()))
			return
		}

		if user == nil {
			t.logger.Info(user)
			ctx.JSON(http.StatusNotFound, ginx.ErrorResponse("Not Found User"))
			return
		}
	}

	err = t.service.Update(ctx, int32(id), req.UserID, req.SeverityId, req.CategoryId, req.SubCategoryId, req.Title, req.Description, req.Status)
	if err != nil {
		t.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, ginx.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusNoContent, ginx.SuccessResponse(""))
}

// @Summary delete a ticket by ID
// @Description delete with the given ID
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200
// @Failure 404 {object} object{error=string}
// @Router /tickets/{id} [delete]
func (t *TicketRouter) hardDelete(ctx *gin.Context) {

	t.logger.Info("Delete ID Ticket")

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		t.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, ginx.ErrorResponse("Bad Request, Id Invalid"))
		return
	}

	err = t.service.Delete(ctx, int32(id))
	if err != nil {
		t.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, ginx.ErrorResponse(errors.FailedToDelete("Ticket")))
		return
	}

	ctx.JSON(http.StatusOK, ginx.SuccessResponse("Ok"))
}

type ITicket interface {
	SetupTicketRoute(routers *gin.RouterGroup)
}

type TicketRouter struct {
	service  service.TicketService
	logger   *zap.SugaredLogger
	validate validator.Validate
}

func NewTicketRouter(s service.TicketService, log *zap.SugaredLogger) *TicketRouter {
	return &TicketRouter{
		service:  s,
		logger:   log,
		validate: *validator.New(),
	}
}

func (t *TicketRouter) SetupTicketRoute(routers *gin.RouterGroup) {
	routers.GET("/tickets", t.list)
	routers.GET("/tickets/:id", t.get)
	routers.DELETE("/tickets/:id", t.hardDelete)
	routers.POST("/tickets", t.create)
	routers.PATCH("/tickets/:id", t.update)
}
