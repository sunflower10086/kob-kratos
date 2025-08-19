package codex

type ResCode int

const (
	CodeSuccess ResCode = 200

	CodeInvalidParams     ResCode = 400
	CodeInvalidTags       ResCode = 400001
	CodeInvalidCategoryId ResCode = 400002

	CodeNeedLogin       ResCode = 401
	CodeCasbinError     ResCode = 401001
	CodeBusinessTypeErr ResCode = 401002
	CodeNonVipLimit     ResCode = 401003
	CodeInValidImage    ResCode = 401004
	CodeUserIsExpired   ResCode = 401005

	CodeInvalidToken                ResCode = 403
	CodeWrongPassword               ResCode = 403001
	CodeWrongUserNameOrPassword     ResCode = 403002
	CodePasswordNotEqual            ResCode = 403007
	CodeWrongEmailCode              ResCode = 403010
	CodeConfirmPasswordError        ResCode = 403011
	CodeOldPasswordError            ResCode = 403012
	CodeEmailCodeExpired            ResCode = 403013
	CodeImgTypeErr                  ResCode = 403014
	CodeImgSizeExceedLimit          ResCode = 403015
	CodeNeedAcceptAgreement         ResCode = 403016
	CodeRepeatPayment               ResCode = 403017
	CodeNotBanSelf                  ResCode = 403018
	CodeEmailTypeErr                ResCode = 403019
	CodePwdLenErr                   ResCode = 403020
	CodePwdContainSpace             ResCode = 403021
	CodePwdTooSimple                ResCode = 403022
	CodeRecQuestionCountExceeded50  ResCode = 403023
	CodeRecQuestionCountExceeded100 ResCode = 403024
	CodeTemplateEditing             ResCode = 403025

	CodeUserNotExist ResCode = 404001
	CodePostNotExist ResCode = 404002

	CodeTplNameIsExist ResCode = 405001

	CodeModelTimeout ResCode = 408001
	CodeUserExist    ResCode = 409002
	CodeQueryUserErr ResCode = 409006

	CodeInternalErr   ResCode = 500
	CodeGenTokenErr   ResCode = 500001
	CodeNoReadingDesc ResCode = 500009

	CodeModelOverLoad ResCode = 503001
)
