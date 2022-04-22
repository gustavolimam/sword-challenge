package models

import "encoding/json"

type Event struct {
	// unique id of the call
	ID string `json:"id"`
	// name of the object to suffer the event
	Object string `json:"object"`
	//  AccountID name of the account who send the message
	AccountID string `json:"account_id"`
	HasError  bool   `json:"has_error"`
	// JSON body of the event
	Properties interface{} `json:"properties"`
}

func (e Event) Json() string {
	b, _ := json.Marshal(e)
	return string(b)
}
