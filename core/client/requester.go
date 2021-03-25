package client

import "encoding/json"

type (
	Getter interface {
		Get(req Request) ([]byte, error)
	}

	StreamGetter interface {
		GetStream(req Request) (chan []byte, error)
	}

	// Requester abstracts the requesting details of supported protocols
	Requester interface {
		Getter
		StreamGetter
		Post(req Request) ([]byte, error)
		Put(req Request) ([]byte, error)
		Delete(req Request) ([]byte, error)
	}

	// Request is a api request abstracting the protocol differences
	Request struct {
		Method  string                 `json:"method,omitempty"`
		Action  string                 `json:"action,omitempty"`
		Node    string                 `json:"node,omitempty"`
		Options map[string]interface{} `json:"options,omitempty"`
	}
)

// NewRequest allocates an unconfigured RequestOptions and returns its
// address.
func NewRequest() *Request {
	r := &Request{}
	r.Options = make(map[string]interface{})
	return r
}

func (t Request) String() string {
	b, _ := json.Marshal(t)
	return "Request" + string(b)
}
