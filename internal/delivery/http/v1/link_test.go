package v1

import (
	"bytes"
	"errors"
	"fmt"
	"link-shortener/internal/domain"
	"link-shortener/internal/service"
	mock_service "link-shortener/internal/service/mocks"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_createToken(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockLink, inp domain.CreateTokenRequest)

	tests := []struct {
		name                 string
		inputBody            string
		input                domain.CreateTokenRequest
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"original_url":"https://github.com/Atasik"}`,
			input: domain.CreateTokenRequest{
				OriginalURL: "https://github.com/Atasik",
			},
			mockBehaviour: func(r *mock_service.MockLink, inp domain.CreateTokenRequest) {
				link := domain.Link{
					OriginalURL: inp.OriginalURL,
				}
				r.EXPECT().CreateToken(link).Return("token", nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"token":"token"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"original_url":"https://github.com/Atasik"}`,
			input: domain.CreateTokenRequest{
				OriginalURL: "https://github.com/Atasik",
			},
			mockBehaviour: func(r *mock_service.MockLink, inp domain.CreateTokenRequest) {
				link := domain.Link{
					OriginalURL: inp.OriginalURL,
				}
				r.EXPECT().CreateToken(link).Return("", errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
		{
			name:      "Bad Payload",
			inputBody: `{"original_url":"wrong input"}`,
			input: domain.CreateTokenRequest{
				OriginalURL: "wrong input",
			},
			mockBehaviour:        func(r *mock_service.MockLink, inp domain.CreateTokenRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid URL"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			servLink := mock_service.NewMockLink(c)
			test.mockBehaviour(servLink, test.input)

			services := &service.Service{Link: servLink}

			h := &Handler{services: services}

			r := mux.NewRouter()
			r.HandleFunc("/link", h.createToken).Methods("POST")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/link",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", appJSON)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getOriginalURL(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockLink, inp domain.GetOriginalURLRequest)

	tests := []struct {
		name                 string
		inputBody            string
		input                domain.GetOriginalURLRequest
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"token":"token12345"}`,
			input: domain.GetOriginalURLRequest{
				Token: "token12345",
			},
			mockBehaviour: func(r *mock_service.MockLink, inp domain.GetOriginalURLRequest) {
				r.EXPECT().GetOriginalURL(inp.Token).Return("https://github.com/Atasik", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"link":"https://github.com/Atasik"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"token":"token12345"}`,
			input: domain.GetOriginalURLRequest{
				Token: "token12345",
			},
			mockBehaviour: func(r *mock_service.MockLink, inp domain.GetOriginalURLRequest) {
				r.EXPECT().GetOriginalURL(inp.Token).Return("", errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
		{
			name:      "Bad Payload",
			inputBody: `{"token":"wrong_input"}`,
			input: domain.GetOriginalURLRequest{
				Token: "wrong_input",
			},
			mockBehaviour:        func(r *mock_service.MockLink, inp domain.GetOriginalURLRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid token"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			servLink := mock_service.NewMockLink(c)
			test.mockBehaviour(servLink, test.input)

			services := &service.Service{Link: servLink}

			h := &Handler{services: services}

			r := mux.NewRouter()
			r.HandleFunc("/link/{token}", h.getOriginalURL).Methods("GET")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/link/%s", test.input.Token), nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
