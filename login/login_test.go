
package login

import (
	"fmt"
	"testing"
)

func TestLogin_Login(t *testing.T) {
	l := New("*****","******!")
	err := l.Login()
	if err != nil {
		fmt.Println(err)
	}

}
