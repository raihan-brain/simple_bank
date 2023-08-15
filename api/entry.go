package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/raihan-brain/simple-bank/db/sqlc"
	"net/http"
)

func addEntryRoutes(rg *gin.Engine, server *Server) {
	entry := rg.Group("/entry")
	{
		entry.POST("/", server.createEntry)
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
