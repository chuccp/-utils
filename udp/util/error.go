package util

import (
	"errors"
)



var ProtocolError = errors.New("ProtocolError")

var RetryVerifyError = errors.New("RetryVerifyError")