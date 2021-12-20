package e

var MapFlags = map[int]string{
	SUCCESS:         "ok",
	ERROR:           "fail",
	INVALID_PARAMS:  "参数类型错误",
	TOKEN_INVALID:   "token无效",
	TOKEN_EXPIRED:   "token过期",
	USER_NOT_FOUND:  "不存在该用户",
	EMAIL_NOT_FOUND: "不存在该邮箱",
}

func GetMsg(code int) string {
	msg, ok := MapFlags[code]
	if ok {
		return msg
	}
	return MapFlags[ERROR]
}
