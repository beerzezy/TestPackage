package null

import "fmt"

type Int64 struct {
	Val   int64
	valid bool
}

func NewInt64(i int64) Int64 {
	return Int64{valid: true, Val: i}
}

func (i Int64) String() string {
	if i.valid == false {
		return "null"
	}
	return fmt.Sprintf("%d", i.Val)
}

func (i Int64) MarshalJSON() ([]byte, error) {
	if !i.valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%d", i.Val)), nil
}
