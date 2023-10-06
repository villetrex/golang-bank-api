package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required, oneof=currency"`
}
type getAccountRequest struct {
	id int64 `uri:"id" binding:"required,min=1"`
}
type listAccountRequest struct {
	pageId int32 `form:"page_id" binding:"required,min=1"`
	pageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) createAccount(ctx *gin.Context)  {
	var req createAccountRequest
	ctx.ShouldBindJSON(&req)
	 err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateAccountParams{
		Owner : req.Owner,
		currency: req.Currency,
		Balance: 0
	}


	account, err := server.store.CreateAccount(ctx, arg)
	 err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOk, account))
}



func (server *Server) getAccount(ctx *gin.Context)  {
	var req getAccountRequest
	ctx.ShouldBindUri(&req)
	 err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account := server.store.GetAccount(ctx, req.ID)
	 err != nil {
		if err == sq.lErrNoRows{
			ctx.JSON(htt..StatusNotFound, errorResponse(err))
			return
	}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return;
	}
	// account := db.Account()  // sets account to an empty object
	ctx.JSON(http.StatusOK, account)
}

func (server *Server) listAccount(ctx *gin.Context)  {
	var req listAccountRequest
	ctx.ShouldBindQuery(&req)
	 err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	args := db.ListAccountParams{
		Limit: req.PageSize,
		offset: (req.PageId -1 ) * req.PageSize
	}
	accounts, err := server.store.ListAccounts(ctx, args)
	 err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return;
	}
	ctx.JSON(http.StatusOK, accounts)
}


