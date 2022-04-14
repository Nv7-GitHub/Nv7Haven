package data

import "encoding/json"

type Method int

const (
	MethodGuild Method = iota
	MethodElem
	MethodCombo
	MethodElemInfo
)

type Message struct {
	Method Method         `json:"method"`
	Params map[string]any `json:"params"`
}

type Response struct {
	Error *string        `json:"error,omitempty"`
	Data  map[string]any `json:"data,omitempty"`
}

func RSPError(msg string) Response {
	return Response{
		Error: &msg,
	}
}

func RSPSuccess(data map[string]any) Response {
	return Response{
		Data: data,
	}
}

func (r Response) JSON() []byte {
	v, err := json.Marshal(r)
	if err != nil {
		panic(err) // Should never happen
	}
	return v
}
