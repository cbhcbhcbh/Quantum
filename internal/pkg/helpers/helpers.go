package helpers

import (
	"fmt"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
)

func GetUuid() string {

	u1 := uuid.NewV4()

	return u1.String()
}

func CreateEmailCode() string {
	return fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
}
