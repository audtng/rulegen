package rules

import (
	"github.com/pquerna/otp/totp"
)

func test() {
	_ = totp.Validate("123456", "secret")
}
