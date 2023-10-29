package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"mekari-test/domain"
	mockdomain "mekari-test/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var timeout = time.Duration(5) * time.Second

func TestEmployeeHandler_Register(t *testing.T) {
	testCases := []struct {
		name          string
		reqBody       string
		wantCode      int
		wantResp      string
		doMockService func(orderService *mockdomain.MockEmployeeService)
	}{
		{
			name:          "Invalid Validate Request Body",
			reqBody:       "",
			wantCode:      http.StatusBadRequest,
			wantResp:      `{"status":"fail","data":"missing required field: first_name"}`,
			doMockService: func(empService *mockdomain.MockEmployeeService) {},
		},
		{
			name:          "Invalid Parse Time",
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
			name:     "Internal Server Error",
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

func TestEmployeeHandler_GetEmployeeById(t *testing.T) {
	var (
		employee = domain.Employee{
			Id:        1,
			FirstName: "Muhammad",
			LastName:  "Fadli",
			Email:     "asd@gmail.com",
			HireDate:  time.Date(2023, time.October, 27, 0, 0, 0, 0, time.UTC),
		}
	)

	testCases := []struct {
		name          string
		urlParam      string
		wantCode      int
		wantResp      string
		doMockService func(orderService *mockdomain.MockEmployeeService)
	}{
		{
			name:          "Invalid Parse URL Param",
			urlParam:      "test",
			wantCode:      http.StatusInternalServerError,
			wantResp:      `{"status":"fail","data":"strconv.Atoi: parsing \"test\": invalid syntax"}`,
			doMockService: func(empService *mockdomain.MockEmployeeService) {},
		},
		{
			name:     "Success",
			urlParam: "1",
			wantCode: http.StatusOK,
			wantResp: `{"status":"success","data":{"id":1,"first_name":"Muhammad","last_name":"Fadli","email":"asd@gmail.com","hire_date":"2023-10-27T00:00:00Z"}}`,
			doMockService: func(empService *mockdomain.MockEmployeeService) {
				empService.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(employee, nil)
			},
		},
		{
			name:     "Error Not Found",
			urlParam: "1",
			wantCode: http.StatusNotFound,
			wantResp: `{"status":"fail","data":"sql: no rows in result set"}`,
			doMockService: func(empService *mockdomain.MockEmployeeService) {
				empService.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(domain.Employee{},
					sql.ErrNoRows)
			},
		},
		{
			name:     "Internal Server Error",
			urlParam: "1",
			wantCode: http.StatusInternalServerError,
			wantResp: `{"status":"fail","data":"internal server error"}`,
			doMockService: func(empService *mockdomain.MockEmployeeService) {
				empService.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(domain.Employee{},
					errors.New("internal server error"))
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
				req = httptest.NewRequest("GET", fmt.Sprintf("/employees/%s", tc.urlParam),
					nil)
				rr = httptest.NewRecorder()
			)
			tc.doMockService(mockService)

			// Create a Params object with the URL parameter
			params := httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: tc.urlParam,
				},
			}

			// Create controller and handle request
			handler.GetEmployeeById(rr, req, params)

			assert.Equal(t, tc.wantCode, rr.Code)
			assert.Equal(t, strings.TrimSpace(tc.wantResp), strings.TrimSpace(rr.Body.String()))
		})
	}
}
