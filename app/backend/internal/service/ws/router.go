package ws

import (
	"encoding/json"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type Req struct {
	Event     string `json:"event"`
	BotId     string `json:"bot_id,omitempty"`
	Direction int    `json:"direction,omitempty"`
}

var Clt *Client

func Router(client *Client, message string) {
	Clt = client
	var data Req
	err := json.Unmarshal([]byte(message), &data)
	if err != nil {
		log.Error(err.Error())
		return
	}

	event := data.Event

	if strings.EqualFold(event, "start-matching") {
		starMatching(data.BotId)
	} else if strings.EqualFold(event, "stop-matching") {
		stopMatching()
	} else if strings.EqualFold(event, "move") {
		move(client, data.Direction)
	}
}

func starMatching(botId string) {
	// TODO: 通过rpc去访问matchingSystem
	log.Debug("start-matching")
	// intBotId, _ := strconv.Atoi(botId)
	// intUserId, _ := strconv.Atoi(Clt.UserId)

	// Stubbed
}

func stopMatching() {
	// TODO: 通过rpc去访问matchingSystem
	log.Debug("stop-matching")

	// Stubbed
}

// 获得前端传来的移动信息(下一步的移动信息，把信息发送到公共区域)
func move(client *Client, d int) {
	////TODO: 把移动的方向发送给game_system
	/*
		snake.Space.ClientDirection <- shape.Pair{
			PlayerId:  client.UserId,
			Direction: strconv.Itoa(d),
		}
	*/
	log.Infof("move direction: %d", d)
}
