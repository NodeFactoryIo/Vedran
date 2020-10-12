package controllers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NodeFactoryIo/vedran/internal/auth"
	"github.com/NodeFactoryIo/vedran/internal/models"
	mocks "github.com/NodeFactoryIo/vedran/mocks/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApiController_PingHandler(t *testing.T) {
	tests := []struct {
		name                  string
		statusCode            int
		pingSaveCallCount     int
		downtimeSaveCallCount int
		pingSaveErr           error
		downtimeSaveErr       error
		calculateDowntimeErr  error
		downtimeDuration      time.Duration
	}{
		{
			name:                  "Returns 500 if downtime calculation fails",
			statusCode:            500,
			pingSaveCallCount:     0,
			pingSaveErr:           nil,
			downtimeSaveErr:       nil,
			downtimeSaveCallCount: 0,
			downtimeDuration:      time.Duration(0),
			calculateDowntimeErr:  fmt.Errorf("ERROR"),
		},
		{
			name:                  "Returns 500 if donwtime save fails",
			statusCode:            500,
			pingSaveCallCount:     0,
			pingSaveErr:           nil,
			downtimeSaveErr:       fmt.Errorf("ERROR"),
			downtimeSaveCallCount: 1,
			downtimeDuration:      time.Duration(time.Second * 11),
			calculateDowntimeErr:  nil,
		},
		{
			name:                  "Saves downtime if downtime duration more than 10 seconds",
			statusCode:            200,
			pingSaveCallCount:     1,
			pingSaveErr:           nil,
			downtimeSaveErr:       nil,
			downtimeSaveCallCount: 1,
			downtimeDuration:      time.Duration(time.Second * 11),
			calculateDowntimeErr:  nil,
		},
		{
			name:                  "Returns 500 if saving ping fails",
			statusCode:            500,
			pingSaveCallCount:     1,
			pingSaveErr:           fmt.Errorf("ERROR"),
			downtimeSaveErr:       nil,
			downtimeSaveCallCount: 0,
			downtimeDuration:      time.Duration(time.Second * 9),
			calculateDowntimeErr:  nil,
		},
		{
			name:                  "Returns 200 and does not save downtime if downtime duration less than 10 seconds",
			statusCode:            500,
			pingSaveCallCount:     1,
			pingSaveErr:           fmt.Errorf("ERROR"),
			downtimeSaveErr:       nil,
			downtimeSaveCallCount: 0,
			downtimeDuration:      time.Duration(time.Second * 9),
			calculateDowntimeErr:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			timestamp := time.Now()
			// create mock controller
			nodeRepoMock := mocks.NodeRepository{}
			pingRepoMock := mocks.PingRepository{}
			metricsRepoMock := mocks.MetricsRepository{}
			downtimeRepoMock := mocks.DowntimeRepository{}

			pingRepoMock.On("Save", &models.Ping{
				NodeId:    "1",
				Timestamp: timestamp,
			}).Return(test.pingSaveErr)
			pingRepoMock.On("CalculateDowntime", mock.Anything, mock.Anything).Return(
				time.Now(), test.downtimeDuration, test.calculateDowntimeErr)
			downtimeRepoMock.On("Save", mock.Anything).Return(test.downtimeSaveErr)

			apiController := NewApiController(false, &nodeRepoMock, &pingRepoMock, &metricsRepoMock, &downtimeRepoMock)
			handler := http.HandlerFunc(apiController.PingHandler)

			// create test request and populate context
			req, _ := http.NewRequest("POST", "/api/v1/node", bytes.NewReader(nil))
			c := &auth.RequestContext{
				NodeId:    "1",
				Timestamp: timestamp,
			}
			ctx := context.WithValue(req.Context(), auth.RequestContextKey, c)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			// invoke test request
			handler.ServeHTTP(rr, req)

			// asserts
			assert.Equal(t, rr.Code, test.statusCode, fmt.Sprintf("Response status code should be %d", test.statusCode))
			assert.True(t, pingRepoMock.AssertNumberOfCalls(t, "Save", test.pingSaveCallCount))
			assert.True(t, downtimeRepoMock.AssertNumberOfCalls(t, "Save", test.downtimeSaveCallCount))
		})
	}
}
