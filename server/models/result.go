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

	return nil, errors.New("Uknown value of ResultEnum: " + r.String())
}

func (r ResultEnum) UnmarshalJSON(data []byte) error {
	sdata := string(data[:])
	switch sdata {
	case "success":
		r = SUCCESS
		return nil
	case "fail":
		r = FAIL
		return nil
	}

	return errors.New("Uknown value for ResultEnum: " + sdata)
}

func (r ResultEnum) String() string {
	switch r {
	case SUCCESS:
		return "success"
	case FAIL:
		return "fail"
	default:
		return "!INVALID VALUE (value: " + string(r) + ")!"
	}
}
