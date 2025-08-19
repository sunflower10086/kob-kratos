package httpencoder

import (
	stdhttp "net/http"

	"kob-kratos/pkg/codex"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

const (
	baseContentType = "application"
)

// ContentType returns the content-type with base prefix.
func ContentType(subtype string) string {
	return baseContentType + "/" + subtype
}

func SuccessEncoder(w http.ResponseWriter, r *http.Request, resp interface{}) error {
	var body Response
	body.Code = stdhttp.StatusOK
	body.Msg = codex.CodeSuccess.Msg()

	codec, _ := http.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(resp)
	if err != nil {
		body = Response{
			Code: 500,
			Msg:  codex.CodeInternalErr.Msg(),
		}
		target := new(errors.Error)
		if errors.As(err, target) {
			body.Code = int(target.Code)
			body.Msg = target.Message
		}
		return err
	}
	w.WriteHeader(stdhttp.StatusOK)

	replyData, err := codec.Marshal(body)
	if err != nil {
		return err
	}

	newData := make([]byte, 0, len(replyData)+len(data)+8)
	newData = append(newData, replyData[:len(replyData)-1]...)
	newData = append(newData, []byte(`,"data":`)...)
	newData = append(newData, data...)
	newData = append(newData, '}')

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(newData)
	if err != nil {
		return err
	}
	return nil
}

func ErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	codec, _ := http.CodecForRequest(r, "Accept")
	w.Header().Set("Content-Type", ContentType(codec.Name()))
	// 返回码均是200
	w.WriteHeader(stdhttp.StatusOK)

	se := errors.FromError(err)
	body := Response{Code: int(codex.CodeInternalErr), Msg: codex.CodeInternalErr.Msg()}

	if se.Code != int32(codex.CodeInternalErr) {
		body.Code = int(se.GetCode())
		body.Msg = se.GetMessage()
	}

	data, err := codec.Marshal(&body)
	if err != nil {
		w.WriteHeader(stdhttp.StatusInternalServerError)
		return
	}

	_, _ = w.Write(data)
}
