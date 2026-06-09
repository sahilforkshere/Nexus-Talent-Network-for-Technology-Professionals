package graph

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

func encodeFeedCursor(index int) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(index)))
}

func decodeFeedCursor(cursor string) (int, error) {
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, fmt.Errorf("invalid cursor")
	}
	return strconv.Atoi(string(b))
}
