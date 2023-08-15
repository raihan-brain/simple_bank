package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/raihan-brain/simple-bank/db/sqlc"
	"net/http"
)

func addEntryRoutes(rg *gin.Engine, server *Server) {
	entry := rg.Group("/entry")
	{
		entry.POST("/", server.createEntry)
		entry.GET("/:id", server.getEntry)
		entry.GET("/list", server.listEntries)
		entry.PUT("/update", server.updateEntry)
	}
}

type createEntryRequest struct {
	Amount    int64 `json:"amount" binding:"required"`
	AccountID int64 `json:"account_id" binding:"required,numeric"`
}

func (server *Server) createEntry(ctx *gin.Context) {
	var req createEntryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Print("internal issue cannot post entry request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateEntryParams{
		Amount:    req.Amount,
		AccountID: req.AccountID,
	}

	entry, err := server.store.CreateEntry(ctx, arg)

	if err != nil {
		fmt.Print("internal issue in db for entry")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

type getEntryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getEntry(ctx *gin.Context) {
	var req getEntryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fmt.Print("request error")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	entry, err := server.store.GetEntry(ctx, req.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		fmt.Print("internal issue")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

type listEntryRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=10"`
}

func (server *Server) listEntries(ctx *gin.Context) {
	var req listEntryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		fmt.Print("internal issue 2")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListEntryParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	entries, err := server.store.ListEntry(ctx, arg)

	if err != nil {
		fmt.Print("internal issue")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

type updateEntryRequest struct {
	ID        int64 `json:"id" binding:"required"`
	Amount    int64 `json:"amount" binding:"required"`
	AccountID int64 `json:"account_id" binding:"required"`
}

func (server *Server) updateEntry(ctx *gin.Context) {

	var req updateEntryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("internal issue updateAccount")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateEntryParams{
		ID:        req.ID,
		AccountID: req.AccountID,
		Amount:    req.Amount,
	}

	entry, err := server.store.UpdateEntry(ctx, arg)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		fmt.Print("internal issue")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}
