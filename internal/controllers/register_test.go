package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/NodeFactoryIo/vedran/internal/configuration"
	"github.com/NodeFactoryIo/vedran/internal/models"
	"github.com/NodeFactoryIo/vedran/internal/repositories"
	"github.com/NodeFactoryIo/vedran/internal/whitelist"
	mocks "github.com/NodeFactoryIo/vedran/mocks/repositories"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestApiController_RegisterHandler(t *testing.T) {
	const TestTunnelServerAddress = "test-tunnel-url:5533"
	configuration.Config = configuration.Configuration{
		TunnelServerAddress: TestTunnelServerAddress,
	}

	// define test cases
	tests := []struct {
		name                          string
		registerRequest               RegisterRequest
		httpStatus                    int
		registerResponse              RegisterResponse
		isWhitelisted                 bool
		saveMockReturns               interface{}
		saveMockCalledNumber          int
	}{
		{
			name: "Valid registration test no whitelist",
			registerRequest: RegisterRequest{
				Id:            "1",
				ConfigHash:    "dadf2e32dwq12",
				PayoutAddress: "0xdafe2cdscdsa",
			},
			httpStatus: http.StatusOK,
			registerResponse: RegisterResponse{
				Token:               "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJub2RlX2lkIjoiMSJ9.LdQLi-cx5HZs6HvVzSFVx0WjXFTsGqDuO9FepXfYLlY",
				TunnelServerAddress: TestTunnelServerAddress,
			},
			isWhitelisted:                 false,
			saveMockReturns:               nil,
			saveMockCalledNumber:          1,
		},
		{
			name: "Valid registration test nodeId on whitelist",
			registerRequest: RegisterRequest{
				Id:            "1",
				ConfigHash:    "dadf2e32dwq12",
				PayoutAddress: "0xdafe2cdscdsa",
			},
			httpStatus: http.StatusOK,
			registerResponse: RegisterResponse{
				Token:               "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJub2RlX2lkIjoiMSJ9.LdQLi-cx5HZs6HvVzSFVx0WjXFTsGqDuO9FepXfYLlY",
				TunnelServerAddress: TestTunnelServerAddress,
			},
			isWhitelisted:                 true,
			saveMockReturns:               nil,
			saveMockCalledNumber:          1,
		},
		{
			name: "Invalid registration test nodeId not on whitelist",
			registerRequest: RegisterRequest{
				Id:            "2",
				ConfigHash:    "dadf2e32dwq12",
				PayoutAddress: "0xdafe2cdscdsa",
			},
			httpStatus:                    http.StatusBadRequest,
			registerResponse:              RegisterResponse{},
			isWhitelisted:                 true,
			saveMockReturns:               nil,
			saveMockCalledNumber:          0,
		},
	}
	_ = os.Setenv("AUTH_SECRET", "test-auth-secret")
	_, _ = whitelist.InitWhitelisting([]string{"1"}, "")

	// execute tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// create mock controller
			nodeRepoMock := mocks.NodeRepository{}
			pingRepoMock := mocks.PingRepository{}
			metricsRepoMock := mocks.MetricsRepository{}
			recordRepoMock := mocks.RecordRepository{}
			nodeRepoMock.On("Save", &models.Node{
				ID:            test.registerRequest.Id,
				ConfigHash:    test.registerRequest.ConfigHash,
				PayoutAddress: test.registerRequest.PayoutAddress,
				Token:         test.registerResponse.Token,
				LastUsed:      time.Now().Unix(),
			}).Return(test.saveMockReturns)

			apiController := NewApiController(test.isWhitelisted, repositories.Repos{
				NodeRepo:    &nodeRepoMock,
				PingRepo:    &pingRepoMock,
				MetricsRepo: &metricsRepoMock,
				RecordRepo:  &recordRepoMock,
			}, nil)

			handler := http.HandlerFunc(apiController.RegisterHandler)

			// create test request
			rb, _ := json.Marshal(test.registerRequest)
			req, err := http.NewRequest("POST", "/api/v1/node", bytes.NewReader(rb))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			// invoke test request
			handler.ServeHTTP(rr, req)

			// asserts
			assert.Equal(t, rr.Code, test.httpStatus, fmt.Sprintf("Response status code should be %d", test.httpStatus))

			var response RegisterResponse
			if rr.Code == http.StatusOK {
				_ = json.Unmarshal(rr.Body.Bytes(), &response)
				assert.Equal(t, response, test.registerResponse, fmt.Sprintf("Response should be %v", test.registerResponse))
			}

			assert.True(t, nodeRepoMock.AssertNumberOfCalls(t, "Save", test.saveMockCalledNumber))
		})
	}
	_ = os.Setenv("AUTH_SECRET", "")
}
