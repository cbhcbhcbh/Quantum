package util

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/cbhcbhcbh/Quantum/pkg/common/domain"
)

func GetServerAddrs(addrs string) []string {
	return strings.Split(addrs, ",")
}

func DecodeToMessage(data []byte) (*domain.Message, error) {
	var msg domain.Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func DecodeToMessagePresenter(data []byte) (*domain.MessagePresenter, error) {
	var msg domain.MessagePresenter
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func GetDurationInMillseconds(start time.Time) float64 {
	end := time.Now()
	duration := end.Sub(start)
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	return rounded
}

func ConstructKey(prefix string, id uint64) string {
	return Join(prefix, ":", strconv.FormatUint(id, 10))
}

func Join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}
