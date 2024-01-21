package restapi

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
)

func CreateTransaction(ctx *gin.Context) {
	var req domain.Transaction
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

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

	// _, err := ph.svc.CreatePayment(ctx, &payment)
	// if err != nil {
	// 	handleError(ctx, err)
	// 	return
	// }

	handleCreatedSuccess(ctx, map[string]string{"test": "test"}, newMeta(limit))
}

func ListTransactions(ctx *gin.Context) {
	// var req listPaymentsRequest
	// var paymentsList []paymentResponse

	// if err := ctx.ShouldBindQuery(&req); err != nil {
	// 	validationError(ctx, err)
	// 	return
	// }

	// payments, err := ph.svc.ListPayments(ctx, req.Skip, req.Limit)
	// if err != nil {
	// 	handleError(ctx, err)
	// 	return
	// }

	// for _, payment := range payments {
	// 	paymentsList = append(paymentsList, newPaymentResponse(&payment))
	// }

	// total := uint64(len(paymentsList))
	// meta := newMeta(total, req.Limit, req.Skip)
	// rsp := toMap(meta, paymentsList, "payments")

	//handleSuccess(ctx, rsp)
}
