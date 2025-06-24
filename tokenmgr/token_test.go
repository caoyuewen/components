package tokenmgr

import (
	"fmt"
	"github.com/caoyuewen/components/util/gen"
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	var tkInfo = TokenInfo{
		Id:          gen.IdString(),
		AppName:     "a",
		Permissions: 3,
		ExpireAt:    time.Now().AddDate(3, 0, 0).Unix(),
	}
	aes := "d7f1a9b3c4e68f12"
	token, err := CreateToken(aes, tkInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
}
