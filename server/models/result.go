package models

import "errors"

type ResultEnum int

const (
	SUCCESS ResultEnum = iota
	FAIL
)

func (r ResultEnum) MarshalJSON() ([]byte, error) {
	switch r {
	case SUCCESS:
		return []byte("success"), nil
	case FAIL:
		return []byte("fail"), nil
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
