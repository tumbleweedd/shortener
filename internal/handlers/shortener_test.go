package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/tumbleweedd/shortener/internal/services"
	mock_services "github.com/tumbleweedd/shortener/internal/services/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_saveURL(t *testing.T) {
	type mockBehavior func(s *mock_services.MockShortener, url string, expectedResult string)

	testTable := []struct {
		name                string
		inputURL            string
		expectedShortCode   string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:              "OK",
			inputURL:          "https://github.com/golang/mock",
			expectedShortCode: "abcd1234",
			mockBehavior: func(s *mock_services.MockShortener, url string, result string) {
				s.EXPECT().ShortenURL(gomock.Any(), url).Return(result, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"code":"abcd1234"}`,
		},
		{
			name:              "Save error",
			inputURL:          "https://github.com/golang/mock",
			expectedShortCode: "",
			mockBehavior: func(s *mock_services.MockShortener, url string, expectedResult string) {
				s.EXPECT().ShortenURL(gomock.Any(), url).Return("", errors.New("invalid save url"))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"invalid save url"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mock_services.NewMockShortener(ctrl)
			testCase.mockBehavior(mockClient, testCase.inputURL, testCase.expectedShortCode)

			service := &services.Service{Shortener: mockClient}
			handler := NewHandler(service)

			// Test server
			r := gin.New()
			r.POST("/a", handler.saveURL)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/a?url=%s", testCase.inputURL), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_redirect(t *testing.T) {
	type mockBehavior func(s *mock_services.MockShortener, code string, expectedResult string)
	testTable := []struct {
		name               string
		code               string
		mockBehavior       mockBehavior
		expectedURL        string
		expectedStatusCode int
	}{
		{
			name: "OK",
			code: "abcd1234",
			mockBehavior: func(s *mock_services.MockShortener, code string, expectedResult string) {
				s.EXPECT().GetLongURL(gomock.Any(), code).Return(expectedResult, nil)
			},
			expectedURL:        "https://example.com",
			expectedStatusCode: http.StatusFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mock_services.NewMockShortener(ctrl)
			testCase.mockBehavior(mockClient, testCase.code, testCase.expectedURL)

			service := &services.Service{Shortener: mockClient}
			handler := NewHandler(service)

			r := gin.New()
			r.GET("/s/:code", handler.redirect)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/s/%s", testCase.code), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedURL, w.Header().Get("Location"))
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}
