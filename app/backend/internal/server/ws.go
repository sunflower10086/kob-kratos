package server

import (
	"fmt"
	"net/http"

	"kob-kratos/app/backend/internal/service/ws"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/mux"
)

func NewWsRouter(hub *ws.Hub, logger log.Logger) http.Handler {
	router := mux.NewRouter()

	h := log.NewHelper(log.With(logger, "module", "service/ws"))

	go hub.Run()

	router.HandleFunc("/ws/{token}", func(w http.ResponseWriter, r *http.Request) {
		// 处理路径参数中的 token
		token := mux.Vars(r)["token"]
		// 解析 token
		fmt.Println(token)
		userId := 0
		ws.WsHandler(hub, int64(userId), w, r, h)
	})

	return router
}
