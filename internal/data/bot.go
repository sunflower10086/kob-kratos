package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"kob-kratos/internal/biz"
	"kob-kratos/internal/data/gormgen/model"
	"kob-kratos/internal/data/gormgen/query"
)

var _ biz.BotRepository = (*botRepo)(nil)

type botRepo struct {
	data *Data
	log  *log.Helper
}

// NewBotRepository 创建机器人仓储实例
func NewBotRepository(data *Data, logger log.Logger) biz.BotRepository {
	return &botRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/bot")),
	}
}

// GetBotList 获取机器人列表
func (r *botRepo) GetBotList(ctx context.Context, page, pageSize int32, userID int32) ([]*biz.Bot, int64, error) {
	b := r.data.DB.WithContext(ctx).Bot

	db := b.
		Where(r.data.DB.Bot.UserID.Eq(userID)).
		Order(r.data.DB.Bot.CreatedAt.Desc())
	total, err := db.Count()
	if err != nil {
		return nil, 0, errors.Wrap(err, "获取机器人列表失败")
	}

	modelBots, err := db.Find()
	if err != nil {
		return nil, 0, errors.Wrap(err, "获取机器人列表失败")
	}

	bots := make([]*biz.Bot, 0, len(modelBots))
	for _, modelBot := range modelBots {
		bot := r.modelToBiz(modelBot)
		bots = append(bots, bot)
	}

	return bots, total, nil
}

// GetBotByID 根据ID获取机器人
func (r *botRepo) GetBotByID(ctx context.Context, botID int32) (*biz.Bot, error) {
	modelBot, err := r.data.DB.WithContext(ctx).Bot.Where(r.data.DB.Bot.ID.Eq(botID)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Warnf("机器人不存在: botID=%d", botID)
			return nil, nil
		}
		r.log.Errorf("获取机器人失败: %v", err)
		return nil, err
	}

	return r.modelToBiz(modelBot), nil
}

// modelToBiz 将数据模型转换为业务实体
func (r *botRepo) modelToBiz(modelBot *model.Bot) *biz.Bot {
	bot := &biz.Bot{
		ID:     modelBot.ID,
		UserID: modelBot.UserID,
		Title:  modelBot.Title,
		Code:   modelBot.Code,
	}

	if modelBot.Description != nil {
		bot.Description = *modelBot.Description
	}

	if modelBot.CreatedAt != nil {
		bot.CreateTime = modelBot.CreatedAt.Format("2006-01-02 15:04:05")
	}

	if modelBot.UpdatedAt != nil {
		bot.ModifyTime = modelBot.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	return bot
}

// Insert 插入机器人（事务方法）
func (r *botRepo) Insert(ctx context.Context, tx *query.Query, bot *biz.Bot) error {
	// 如果tx为空，使用r.data.DB
	db := tx
	if db == nil {
		db = r.data.DB
	}

	modelBot := &model.Bot{
		UserID:      bot.UserID,
		Title:       bot.Title,
		Description: &bot.Description,
		Code:        bot.Code,
	}

	if err := db.WithContext(ctx).Bot.Create(modelBot); err != nil {
		r.log.Errorf("插入机器人失败: %v", err)
		return err
	}

	// 更新业务实体的ID
	bot.ID = modelBot.ID
	return nil
}

// Update 更新机器人（事务方法）
func (r *botRepo) Update(ctx context.Context, tx *query.Query, bot *biz.Bot) error {
	// 如果tx为空，使用r.data.DB
	db := tx
	if db == nil {
		db = r.data.DB
	}

	updateData := map[string]interface{}{
		"title":       bot.Title,
		"description": bot.Description,
		"code":        bot.Code,
		"updated_at":  time.Now(),
	}

	result, err := db.WithContext(ctx).Bot.
		Where(db.Bot.ID.Eq(bot.ID), db.Bot.UserID.Eq(bot.UserID)).
		Updates(updateData)
	if err != nil {
		r.log.Errorf("更新机器人失败: %v", err)
		return err
	}

	if result.RowsAffected == 0 {
		r.log.Warnf("机器人不存在或无权限更新: botID=%d, userID=%d", bot.ID, bot.UserID)
		return gorm.ErrRecordNotFound
	}

	return nil
}

// DeleteBot 删除机器人（事务方法）
func (r *botRepo) DeleteBot(ctx context.Context, tx *query.Query, botID int32) error {
	// 如果tx为空，使用r.data.DB
	db := tx
	if db == nil {
		db = r.data.DB
	}

	result, err := db.WithContext(ctx).Bot.
		Where(db.Bot.ID.Eq(botID)).
		Delete()
	if err != nil {
		r.log.Errorf("删除机器人失败: %v", err)
		return err
	}

	if result.RowsAffected == 0 {
		r.log.Warnf("机器人不存在: botID=%d", botID)
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Transaction 执行事务操作
func (r *botRepo) Transaction(ctx context.Context, fn func(tx *query.Query) error) error {
	return r.data.DB.Transaction(fn)
}
