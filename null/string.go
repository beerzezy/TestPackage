package null

import (
	"encoding/json"
)

type String struct {
	Val   string
	valid bool
}

func NewString(s string) String {
	return String{valid: true, Val: s}
}

func (s String) String() string {
	if s.valid == false {
		return "null"
	}
	return s.Val
}

func (s String) MarshalJSON() ([]byte, error) {
	if s.valid == false {
		return []byte("null"), nil
	}

	return json.Marshal(s.Val)
}
