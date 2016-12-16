package models

import "errors"

type ResultEnum int

const (
	SUCCESS ResultEnum = iota
	FAIL
)

func (this ResultEnum) MarshalJSON() ([]byte, error) {
	switch this {
	case SUCCESS:
		return []byte("success"), nil
	case FAIL:
		return []byte("fail"), nil
	}

	return nil, errors.New("Uknown value of ResultEnum: " + this)
}

func (this ResultEnum) UnmarshalJSON(data []byte) error {
	switch data {
	case []byte("success"):
		this = SUCCESS
		return nil
	case []byte("fail"):
		this = FAIL
		return nil
	}

	return errors.New("Uknown value for ResultEnum: " + data)
}
