package utils

import "regexp"

const (
	PHONE = `^(13[0-9]|14[5|7]|15[0|1|2|3|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9]|17[7])\d{8}$`
)

//手机号验证
func VerifyPhone(phone string) bool {
	match, err := regexp.Match(PHONE, []byte(phone))
	if err != nil {
		return false
	}
	return match
}
