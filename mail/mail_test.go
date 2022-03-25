
package mail

import (
	"fmt"
	"testing"
)

func TestSendMail(t *testing.T) {
	err := SendMail("123", "******@qq.com", "123", "123", "")
	if err != nil {
		fmt.Println(err)
		return
	}
}
