
package config

import (
	"fmt"
	"testing"
	"time"
)

func TestGetConfig(t *testing.T) {
	info := GetConfig("../config.yaml")
	fmt.Println(info.Users.Address)
	fmt.Println(info.Users.Username)
	fmt.Println(info.Users.Password)
	fmt.Println(time.Now().Format("2006-01-02"))
}
