package encrypt_test

import (
	"fmt"
	"testing"

	"kob-kratos/pkg/helper/encrypt"
)

func TestPasswordHash(t *testing.T) {
	fmt.Println(encrypt.PasswordHash("123456"))
}
