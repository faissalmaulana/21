package utils

import "github.com/faissalmaulana/21/api/internal/model"

func BoolPtr(b bool) *bool {
	return &b
}

// string values that available are (open,done)
// if the input not one of them will return the default (open)
func ToStatus(val string) model.Status {
	mapStatus := map[string]model.Status{
		"open": model.Status(0),
		"done": model.Status(1),
	}

	if status, ok := mapStatus[val]; ok {
		return status
	}

	return model.Status(0)
}
