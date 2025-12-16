package server

import (
	"context"
	"log"
	"math"
	"sync"
	"time"
)

type Player struct {
	UserId   int32
	BotId    int32
	Rating   int32
	WaitTime int32
}

type MatchPool struct {
	players []*Player
	mu      sync.Mutex
	ch      chan struct{} // Close this channel to stop the pool
}

func NewMatchPool() *MatchPool {
	return &MatchPool{
		players: make([]*Player, 0),
		ch:      make(chan struct{}),
	}
}

func (p *MatchPool) AddPlayer(player *Player) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.players = append(p.players, player)
}

func (p *MatchPool) RemovePlayer(userId int32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	newPlayers := make([]*Player, 0)
	for i := 0; i < len(p.players); i++ {
		if p.players[i].UserId != userId {
			newPlayers = append(newPlayers, p.players[i])
		}
	}
	p.players = newPlayers
}

func (p *MatchPool) Stop(ctx context.Context) error {
	close(p.ch)

	return nil
}

func (p *MatchPool) Start(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.ch:
			return nil
		case <-ticker.C:
			p.matchPlayers()
		case <-ctx.Done():
			return nil
		}
	}

	return nil
}

// 尝试匹配所有玩家
func (p *MatchPool) matchPlayers() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.increaseWaitingTime()

	used := make([]bool, len(p.players))

	var newPlayers []*Player
	// Optimize: In-place filtering could be done, but for simplicity and safety against index shifts during matching loop,
	// we'll stick to marking used and then compacting.
	// Actually, we can just rebuild the slice.

	for i := 0; i < len(p.players); i++ {
		if p.players[i].UserId == 0 || used[i] {
			continue
		}
		for j := i + 1; j < len(p.players); j++ {
			if p.players[j].UserId == 0 || used[j] {
				continue
			}

			if p.players[i].UserId == p.players[j].UserId {
				continue
			}
			a, b := p.players[i], p.players[j]

			if p.checkMatched(a, b) {
				used[i] = true
				used[j] = true
				go p.sendResult(a, b)
				break
			}
		}
	}

	// Compact the slice
	for i := 0; i < len(p.players); i++ {
		if !used[i] {
			newPlayers = append(newPlayers, p.players[i])
		}
	}
	p.players = newPlayers
}

// 给所有玩家的匹配时间增加一秒
// Note: This must be called under Lock
func (p *MatchPool) increaseWaitingTime() {
	for _, player := range p.players {
		player.WaitTime += 1
	}
}

func (p *MatchPool) checkMatched(a, b *Player) bool {
	ratingDelta := math.Abs(float64(a.Rating) - float64(b.Rating))
	waitTime := math.Min(float64(a.WaitTime), float64(b.WaitTime))
	return ratingDelta <= waitTime*10
}

func (p *MatchPool) sendResult(a, b *Player) {
	log.Printf("Matched: Player %d and Player %d", a.UserId, b.UserId)
	// TODO : 把结果返回给backend层
	// EventType 0表示匹配的结果，1表示游戏的结果
	// req := pb.ResultReq{
	// 	EventType: 0,
	// 	MatchResult: &pb.MatchResult{
	// 		AId:    a.UserId,
	// 		ABotId: a.BotId,
	// 		BId:    b.UserId,
	// 		BBotId: b.BotId,
	// 	},
	// }
	// ctx := context.Background()
	// result, err := result.Result(ctx, &req)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println("sendResult", result)
}
