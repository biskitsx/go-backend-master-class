package api

import (
	"database/sql"
	"errors"
	"fmt"

	"net/http"

	db "github.com/biskitsx/go-backend-master-class/db/sqlc"
	"github.com/biskitsx/go-backend-master-class/token"
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
	fromAccount, valid := server.validAccount(ctx, reqBody.FromAccountID, reqBody.Currency)
	if !valid {
		err := errors.New("from account doesn't exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	_, valid = server.validAccount(ctx, reqBody.ToAccountID, reqBody.Currency)
	if !valid {
		err := errors.New("from account doesn't exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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

func (server *Server) validAccount(ctx *gin.Context, AccountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false

	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch:", account.ID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false

	}
	return account, true

}
