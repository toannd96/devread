package helper

import (
	"crypto"
	_ "crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
)

func Hash(objs ...interface{}) string {
	digester := crypto.MD5.New()
	for _, ob := range objs {
		fmt.Fprint(digester, reflect.TypeOf(ob))
		fmt.Fprint(digester, ob)
	}
	theHash := hex.EncodeToString(digester.Sum(nil))
	return theHash
}
