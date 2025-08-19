package codex

var codeMsgMap = map[ResCode]string{
	CodeSuccess:           "操作成功",
	CodeInvalidParams:     "请求参数错误",
	CodeInvalidTags:       "非法的 tagsId",
	CodeInvalidCategoryId: "非法的类别Id",

	CodeNeedLogin:               "请先登陆",
	CodeCasbinError:             "没有该API接口访问权限",
	CodeBusinessTypeErr:         "业务类型错误",
	CodeUserIsExpired:           "用户体验时间已结束",
	CodeNonVipLimit:             "操作受限制，需要开通会员服务",
	CodeInvalidToken:            "无效的token",
	CodeWrongPassword:           "密码错误",
	CodeWrongUserNameOrPassword: "用户名或密码错误",
	CodePasswordNotEqual:        "两次密码不一致",
	CodeConfirmPasswordError:    "两次密码输入不一致",
	CodeOldPasswordError:        "旧密码错误",
	CodeUserNotExist:            "该用户不存在",
	CodeUserExist:               "该用户已存在",
	CodeInternalErr:             "服务器开小差啦，稍后再来试一试",
	CodeGenTokenErr:             "生成Token异常",
	CodeQueryUserErr:            "查询用户信息出错",
	CodeImgTypeErr:              "图片格式错误",
	CodeImgSizeExceedLimit:      "图片大小超过限制",
	CodeNotBanSelf:              "不能对自己进行操作",
	CodeEmailTypeErr:            "邮箱格式错误",
	CodeNeedAcceptAgreement:     "请先接受协议",
	CodePostNotExist:            "帖子不存在",
	CodeInValidImage:            "图片中没有识别到题目",
	CodeRepeatPayment:           "请勿重复支付",

	CodeTplNameIsExist:  "同年级同学科下模板名不可重复",
	CodeTemplateEditing: "不可使用未完成的模板",

	CodePwdLenErr:                   "密码长度小于6或大于20",
	CodePwdContainSpace:             "密码输入包含空格",
	CodePwdTooSimple:                "密码为6-20位字符，必须字母、数字两种以上组合，区分大小写",
	CodeNoReadingDesc:               "请勾选用户同意协议",
	CodeRecQuestionCountExceeded50:  "推题数量超过50道",
	CodeRecQuestionCountExceeded100: "推题数量超过100道",
	CodeModelTimeout:                "模型响应超时",
	CodeModelOverLoad:               "模型服务过载",
}

func (code ResCode) Msg() string {
	msg, ok := codeMsgMap[code]
	if !ok {
		msg = codeMsgMap[CodeInternalErr]
	}
	return msg
}
