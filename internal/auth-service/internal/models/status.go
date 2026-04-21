package models

import "fmt"

type Status string

const (
	StatusUnspecified Status = ""

	StatusActive   Status = "ACTIVE"
	StatusBlocked  Status = "BLOCKED"
	StatusLocked   Status = "LOCKED"
	StatusDisabled Status = "DISABLED"
)

func ToStatus(st string) (Status, error) {
	switch st {
	case string(StatusActive):
		return StatusActive, nil
	case string(StatusBlocked):
		return StatusBlocked, nil
	case string(StatusLocked):
		return StatusLocked, nil
	case string(StatusDisabled):
		return StatusDisabled, nil
	default:
		return StatusUnspecified, fmt.Errorf("invalid status: %s", st)
	}
}
