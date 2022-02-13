package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/s14t284/simplebank/db/ent"
	"github.com/s14t284/simplebank/ent"
)

type createTransferRequest struct {
	FromAccountID int    `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int    `json:"to_account_id" binding:"required,min=1"`
	Amount        int    `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// validation
	if !server.validateAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}
	if !server.validateAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	result, err := server.store.TransferTx(ctx, db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validateAccount(ctx *gin.Context, accountID int, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if ent.IsNotFound(err) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
