package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/NodeFactoryIo/vedran/internal/models"
	"github.com/NodeFactoryIo/vedran/internal/rpc"
	mocks "github.com/NodeFactoryIo/vedran/mocks/models"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
)

type handleFnMock func(http.ResponseWriter, *http.Request)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
}

func teardown() {
	server.Close()
}

func TestApiController_RPCHandler(t *testing.T) {
	setup()
	defer teardown()

	nodeRepoMock := mocks.NodeRepository{}
	pingRepoMock := mocks.PingRepository{}
	metricsRepoMock := mocks.MetricsRepository{}
	apiController := NewApiController(false, &nodeRepoMock, &pingRepoMock, &metricsRepoMock)
	handler := http.HandlerFunc(apiController.RPCHandler)

	tests := []struct {
		name        string
		rpcRequest  string
		rpcResponse rpc.RPCResponse
		nodes       []models.Node
		handleFunc  handleFnMock
	}{
		{
			name:       "Returns parse error if json invalid",
			rpcRequest: `INVALID`,
			rpcResponse: rpc.RPCResponse{
				ID:      0,
				JSONRPC: "2.0",
				Error:   &rpc.RPCError{Code: -32700, Message: "Parse error"}},
		},
		{
			name:       "Returns invalid request if rpc request is invalid rpc request",
			rpcRequest: `{}`,
			rpcResponse: rpc.RPCResponse{
				ID:      0,
				JSONRPC: "2.0",
				Error:   &rpc.RPCError{Code: -32600, Message: "Invalid request"}},
		},
		{
			name:       "Returns server error if no available nodes",
			rpcRequest: `{"jsonrpc": "2.0", "id": 1, "method": "system"}`,
			rpcResponse: rpc.RPCResponse{
				ID:      1,
				JSONRPC: "2.0",
				Error:   &rpc.RPCError{Code: -32603, Message: "No available nodes"},
			},
			nodes: []models.Node{},
		},
		{
			name:       "Returns server error if all nodes return invalid rpc response",
			rpcRequest: `{"jsonrpc": "2.0", "id": 1, "method": "system"}`,
			rpcResponse: rpc.RPCResponse{
				ID:      1,
				JSONRPC: "2.0",
				Error:   &rpc.RPCError{Code: -32603, Message: "Internal Server Error"},
			},
			nodes: []models.Node{{ID: "test-id", NodeUrl: "invalid"}},
		},
		{
			name:       "Returns response if node returnes valid rpc response",
			rpcRequest: `{"jsonrpc": "2.0", "id": 1, "method": "system"}`,
			rpcResponse: rpc.RPCResponse{
				ID:      1,
				JSONRPC: "2.0",
				Error:   nil,
			},
			nodes: []models.Node{{ID: "test-id", NodeUrl: "valid"}},
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.WriteString(w, `{"id": 1, "jsonrpc": "2.0"}`)
			}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setup()

			if test.handleFunc != nil {
				mux.HandleFunc("/", test.handleFunc)
			}
			if len(test.nodes) > 0 && test.nodes[0].NodeUrl == "valid" {
				test.nodes[0].NodeUrl = server.URL
			} else if len(test.nodes) > 0 {
				test.nodes[0].NodeUrl = "INVALID"
			}

			nodeRepoMock.On("GetActiveNodes").Return(&test.nodes, nil)
			req, _ := http.NewRequest("POST", "/", bytes.NewReader([]byte(test.rpcRequest)))
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			var body rpc.RPCResponse
			_ = json.Unmarshal(rr.Body.Bytes(), &body)
			if !reflect.DeepEqual(body, test.rpcResponse) {
				t.Errorf("SendRequestToNode() body = %v, want %v", body, test.rpcResponse)
				return
			}

			teardown()
		})
	}
}