package internal

import "errors"

const errMsg = "collision detected"

func PanicOnFound(found bool) {
	if found {
		panic(errors.New(errMsg))
	}
}

func CheckForLockOrReturnValue(value interface{}) (rtnValue interface{}) {
	switch value := value.(type) {
	case Lock:
		panic(errors.New(errMsg))
	default:
		rtnValue = value
	}
	return
}
