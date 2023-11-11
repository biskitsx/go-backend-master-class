package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/biskitsx/go-backend-master-class/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var reqBody createTransferRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !server.validAccount(ctx, reqBody.FromAccountID, reqBody.Currency) || !server.validAccount(ctx, reqBody.ToAccountID, reqBody.Currency) {
		return
	}

	arg := db.TranferTxParams{
		FromAccountID: reqBody.FromAccountID,
		ToAccountID:   reqBody.ToAccountID,
		Amount:        reqBody.Amount,
	}

	result, err := server.store.TranferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, AccountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch:", account.ID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}

// type getAccountRequest struct {
// 	ID uint64 `uri:"id" binding:"required,min=1"`
// }

// func (server *Server) getAccount(ctx *gin.Context) {
// 	var reqUri getAccountRequest
// 	if err := ctx.ShouldBindUri(&reqUri); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	account, err := server.store.GetAccount(ctx, int64(reqUri.ID))
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, account)
// }

// type listAccountRequest struct {
// 	PageID   uint `form:"page_id" binding:"required,min=1"`
// 	PageSize uint `form:"page_size" binding:"required,min=1"`
// }

// func (server *Server) listAccount(ctx *gin.Context) {
// 	var reqQuery listAccountRequest
// 	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}
// 	arg := db.ListAccountParams{
// 		Limit:  int32(reqQuery.PageSize),
// 		Offset: int32(reqQuery.PageID-1) * int32(reqQuery.PageSize),
// 	}
// 	accounts, err := server.store.ListAccount(ctx, arg)

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, accounts)
// }
