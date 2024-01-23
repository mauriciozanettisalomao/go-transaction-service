package restapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/service"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {

	tests := []struct {
		create             func(context.Context, *domain.Transaction) error
		name               string
		payload            string
		expectedStatusCode int
	}{
		{
			name: "valid scenario",
			payload: `{
				"amount": 999.99,
				"currency": "dollar",
				"origin": "mobile",
				"user": {
					"id": "999"
				},
				"operationType": "debit"
			}`,
			create:             func(context.Context, *domain.Transaction) error { return nil },
			expectedStatusCode: 201,
		},
		{
			name: "bad request - no amout",
			payload: `{
				"currency": "dollar",
				"origin": "mobile",
				"user": {
					"id": "999"
				},
				"operationType": "debit"
			}`,
			create:             func(context.Context, *domain.Transaction) error { return nil },
			expectedStatusCode: 400,
		},
		{
			name: "error on create transaction",
			payload: `{
				"amount": 999.99,
				"currency": "dollar",
				"origin": "mobile",
				"user": {
					"id": "999"
				},
				"operationType": "debit"
			}`,
			create: func(context.Context, *domain.Transaction) error {
				return domain.ErrDataAlreadyExists
			},
			expectedStatusCode: 409,
		},
		{
			name: "unexpected error on create transaction",
			payload: `{
				"amount": 999.99,
				"currency": "dollar",
				"origin": "mobile",
				"user": {
					"id": "999"
				},
				"operationType": "debit"
			}`,
			create: func(context.Context, *domain.Transaction) error {
				return errors.New("unexpected error")
			},
			expectedStatusCode: 500,
		},
	}

	assertion := assert.New(t)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			mock := &service.TransactionServiceMock{
				CreateMock: tc.create,
			}

			api := &transactionAPI{
				svcTransaction: mock,
			}

			request, _ := gin.CreateTestContext(httptest.NewRecorder())
			request.Request = &http.Request{
				Body: io.NopCloser(strings.NewReader(tc.payload)),
			}

			api.CreateTransaction(request)

			assertion.Equal(tc.expectedStatusCode, request.Writer.Status())
		})
	}

}
