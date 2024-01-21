package restapi

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/service"
)

const (
	xIdempontencyKey = "x-idempontency-key"
)

func CreateTransaction(ctx *gin.Context) {
	var req domain.Transaction
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	req.IdempontencyKey = ctx.GetHeader(xIdempontencyKey)
	err := service.NewTransactionHandler().Create(ctx, &req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleCreatedSuccess(ctx, req, nil)
}

func ListTransactions(ctx *gin.Context) {

	fmt.Println(" limit", ctx.Query("limit"))
	limit := 10
	limitParam, ok := ctx.GetQuery("limit")
	if ok {
		limitConv, err := strconv.Atoi(limitParam)
		if err != nil {
			fmt.Println(err)

		}
		limit = limitConv
	}

	response, err := service.NewTransactionHandler().List(ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, response, newMeta(limit))
}
