package util

import (
	"github.com/rs/xid"
)

//GenNextUUID return a new UUID
func GenNextUUID() string {
	guid := xid.New()
	return guid.String()
}
