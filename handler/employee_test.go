package handler

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockdomain "mekari-test/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestEmployeeHandler_Register(t *testing.T) {
	var timeout = time.Duration(5) * time.Second

	testCases := []struct {
		name          string
		reqBody       string
		wantCode      int
		wantResp      string
		doMockService func(orderService *mockdomain.MockEmployeeService)
	}{
		{
			name:          "Invalid validate request body",
			reqBody:       "",
			wantCode:      http.StatusBadRequest,
			wantResp:      `{"status":"fail","data":"missing required field: first_name"}`,
			doMockService: func(empService *mockdomain.MockEmployeeService) {},
		},
		{
			name:          "Invalid parse time",
			reqBody:       `first_name=John&last_name=Doe&email=johndoe@example.com&hire_date=2022-12-122`,
			wantCode:      http.StatusBadRequest,
			wantResp:      `{"status":"fail","data":"invalid parsing date, err:parsing time \"2022-12-122\": extra text: \"2\""}`,
			doMockService: func(empService *mockdomain.MockEmployeeService) {},
		},
		{
			name:     "Success",
			reqBody:  `first_name=John&last_name=Doe&email=johndoe@example.com&hire_date=2022-12-12`,
			wantCode: http.StatusOK,
			wantResp: `{"status":"success"}`,
			doMockService: func(empService *mockdomain.MockEmployeeService) {
				empService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:     "Failure",
			reqBody:  `first_name=John&last_name=Doe&email=johndoe@example.com&hire_date=2022-12-12`,
			wantCode: http.StatusInternalServerError,
			wantResp: `{"status":"fail","data":"internal server error"}`,
			doMockService: func(empService *mockdomain.MockEmployeeService) {
				empService.EXPECT().Register(gomock.Any(), gomock.Any()).
					Return(errors.New("internal server error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				mockCtrl    = gomock.NewController(t)
				mockService = mockdomain.NewMockEmployeeService(mockCtrl)
				handler     = EmployeeHandler{
					employeeService: mockService,
					timeout:         timeout,
				}
				req = httptest.NewRequest("POST", "/employees",
					strings.NewReader(tc.reqBody))
				rr = httptest.NewRecorder()
			)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			tc.doMockService(mockService)

			// Create controller and handle request
			handler.Register(rr, req, nil)

			assert.Equal(t, tc.wantCode, rr.Code)
			assert.Equal(t, strings.TrimSpace(tc.wantResp), strings.TrimSpace(rr.Body.String()))
		})
	}
}
