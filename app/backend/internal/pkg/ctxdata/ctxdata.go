package ctxdata

import (
	"context"
)

var CtxKeyUid = "uid"

func GetUid(ctx context.Context) (int64, error) {
	uidFloat := ctx.Value(CtxKeyUid).(float64)
	return int64(uidFloat), nil
}
