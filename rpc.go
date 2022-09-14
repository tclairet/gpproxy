package gpproxy

import (
	"encoding/json"
)

// Request defines ethereum rpc request
type Request struct {
	ID      uint          `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// Bytes convert the request in bytes
func (r *Request) Bytes() []byte {
	b, _ := json.Marshal(r)
	return b
}

// MakeResponse returns a response for a request
func MakeResponse(r Request, result interface{}) *Response {
	return &Response{
		Jsonrpc: r.Jsonrpc,
		ID:      r.ID,
		Result:  result,
	}
}

// Response defines ethereum rpc response
type Response struct {
	ID      uint        `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error for RPC Error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
