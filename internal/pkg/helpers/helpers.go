package helpers

import (
	uuid "github.com/satori/go.uuid"
)

func GetUuid() string {

	u1 := uuid.NewV4()

	return u1.String()
}
