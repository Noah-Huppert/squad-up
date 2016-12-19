package models

import "errors"

type ResultEnum int

const (
	SUCCESS ResultEnum = iota
	FAIL
)

func (r ResultEnum) MarshalJSON() ([]byte, error) {
	if r == SUCCESS || r == FAIL {
		return []byte(r.String()), nil
	}

	return nil, errors.New("Uknown value of ResultEnum: " + r)
}

func (r ResultEnum) UnmarshalJSON(data []byte) error {
	switch data {
	case []byte("success"):
		r = SUCCESS
		return nil
	case []byte("fail"):
		r = FAIL
		return nil
	}

	return errors.New("Uknown value for ResultEnum: " + data)
}

func (r ResultEnum) String() string {
	switch r {
	case SUCCESS:
		return "success"
	case FAIL:
		return "fail"
	default:
		return "!INVALID VALUE (value: " + r + ")!"
	}
}