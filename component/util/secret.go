package util

import (
	"fmt"
	"log"
	"crypto/rand"
	"github.com/rs/xid"
)

// MakeUniqueID 는 12bytes 20chars a globally unique ID
// with timestamp.
func MakeUniqueID() string {
	return xid.New().String()
}

func MakeUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
    	log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",	b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}
