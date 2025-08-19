package data

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"kob-kratos/internal/biz"
	"kob-kratos/internal/data/gormgen/model"
	"kob-kratos/internal/data/gormgen/query"
)

var _ biz.RecordRepository = (*recordRepo)(nil)

type recordRepo struct {
	data *Data
	log  *log.Helper
}

// NewRecordRepository 创建记录仓储实例
func NewRecordRepository(data *Data, logger log.Logger) biz.RecordRepository {
	return &recordRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/record")),
	}
}

// GetRecordList 获取记录列表（支持分页大小参数）
func (r *recordRepo) GetRecordList(ctx context.Context, page, pageSize int32) ([]*biz.Record, int64, error) {
	// 如果pageSize为0，使用默认值
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 查询记录总数
	total, err := r.data.DB.WithContext(ctx).Record.Count()
	if err != nil {
		r.log.Errorf("获取记录总数失败: %v", err)
		return nil, 0, err
	}

	// 查询记录列表
	modelRecords, err := r.data.DB.WithContext(ctx).Record.
		Order(r.data.DB.Record.CreatedAt.Desc()).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Find()
	if err != nil {
		r.log.Errorf("获取记录列表失败: %v", err)
		return nil, 0, err
	}

	// 转换为业务实体
	records := make([]*biz.Record, 0, len(modelRecords))
	for _, modelRecord := range modelRecords {
		record, err := r.modelToBizWithUserInfo(ctx, modelRecord)
		if err != nil {
			r.log.Errorf("转换记录失败: %v", err)
			continue
		}
		records = append(records, record)
	}

	return records, total, nil
}

// CreateRecord 创建游戏记录（事务方法）
func (r *recordRepo) CreateRecord(ctx context.Context, tx *query.Query, record *biz.GameRecord) error {
	// 如果tx为空，使用r.data.DB
	db := tx
	if db == nil {
		db = r.data.DB
	}

	modelRecord := &model.Record{
		AID:       record.AID,
		ASx:       record.ASX,
		ASy:       record.ASY,
		BID:       record.BID,
		BSx:       record.BSX,
		BSy:       record.BSY,
		ASteps:    record.ASteps,
		BSteps:    record.BSteps,
		GameMap:   record.Map,
		LoserName: record.Loser,
	}

	// 根据loser名称确定loser_id
	if record.Loser == "A" {
		modelRecord.LoserID = int64(record.AID)
	} else if record.Loser == "B" {
		modelRecord.LoserID = int64(record.BID)
	}

	if err := db.WithContext(ctx).Record.Create(modelRecord); err != nil {
		r.log.Errorf("创建游戏记录失败: %v", err)
		return err
	}

	// 更新业务实体的ID
	record.ID = modelRecord.ID
	return nil
}

// GetRecordByID 根据ID获取记录
func (r *recordRepo) GetRecordByID(ctx context.Context, recordID int32) (*biz.Record, error) {
	modelRecord, err := r.data.DB.WithContext(ctx).Record.Where(r.data.DB.Record.ID.Eq(recordID)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Warnf("记录不存在: recordID=%d", recordID)
			return nil, nil
		}
		r.log.Errorf("获取记录失败: %v", err)
		return nil, err
	}

	return r.modelToBizWithUserInfo(ctx, modelRecord)
}

// GetUserRecords 获取用户的游戏记录（支持分页大小参数）
func (r *recordRepo) GetUserRecords(ctx context.Context, userID int32, page, pageSize int32) ([]*biz.Record, int64, error) {
	// 如果pageSize为0，使用默认值
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 查询用户记录总数
	total, err := r.data.DB.WithContext(ctx).Record.
		Where(r.data.DB.Record.AID.Eq(userID)).
		Or(r.data.DB.Record.BID.Eq(userID)).
		Count()
	if err != nil {
		r.log.Errorf("获取用户记录总数失败: %v", err)
		return nil, 0, err
	}

	// 查询用户记录列表
	modelRecords, err := r.data.DB.WithContext(ctx).Record.
		Where(r.data.DB.Record.AID.Eq(userID)).
		Or(r.data.DB.Record.BID.Eq(userID)).
		Order(r.data.DB.Record.CreatedAt.Desc()).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Find()
	if err != nil {
		r.log.Errorf("获取用户记录列表失败: %v", err)
		return nil, 0, err
	}

	// 转换为业务实体
	records := make([]*biz.Record, 0, len(modelRecords))
	for _, modelRecord := range modelRecords {
		record, err := r.modelToBizWithUserInfo(ctx, modelRecord)
		if err != nil {
			r.log.Errorf("转换用户记录失败: %v", err)
			continue
		}
		records = append(records, record)
	}

	return records, total, nil
}

// modelToBizWithUserInfo 将数据模型转换为业务实体（包含用户信息）
func (r *recordRepo) modelToBizWithUserInfo(ctx context.Context, modelRecord *model.Record) (*biz.Record, error) {
	// 获取玩家A的信息
	userA, err := r.data.DB.WithContext(ctx).User.Where(r.data.DB.User.ID.Eq(modelRecord.AID)).First()
	if err != nil {
		r.log.Errorf("获取玩家A信息失败: %v", err)
		return nil, err
	}

	// 获取玩家B的信息
	userB, err := r.data.DB.WithContext(ctx).User.Where(r.data.DB.User.ID.Eq(modelRecord.BID)).First()
	if err != nil {
		r.log.Errorf("获取玩家B信息失败: %v", err)
		return nil, err
	}

	// 构建游戏记录
	gameRecord := &biz.GameRecord{
		ID:     modelRecord.ID,
		AID:    modelRecord.AID,
		ASX:    modelRecord.ASx,
		ASY:    modelRecord.ASy,
		BID:    modelRecord.BID,
		BSX:    modelRecord.BSx,
		BSY:    modelRecord.BSy,
		ASteps: modelRecord.ASteps,
		BSteps: modelRecord.BSteps,
		Map:    modelRecord.GameMap,
		Loser:  modelRecord.LoserName,
	}

	if modelRecord.CreatedAt != nil {
		gameRecord.CreateTime = modelRecord.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// 构建记录实体
	record := &biz.Record{
		AUsername: userA.Username,
		BUsername: userB.Username,
		Record:    gameRecord,
	}

	// 设置头像
	if userA.Photo != nil {
		record.APhoto = *userA.Photo
	}
	if userB.Photo != nil {
		record.BPhoto = *userB.Photo
	}

	// 设置比赛结果
	if modelRecord.LoserName == "A" {
		record.Result = fmt.Sprintf("%s 胜利", userB.Username)
	} else if modelRecord.LoserName == "B" {
		record.Result = fmt.Sprintf("%s 胜利", userA.Username)
	} else {
		record.Result = "平局"
	}

	return record, nil
}

// Transaction 执行事务操作
func (r *recordRepo) Transaction(ctx context.Context, fn func(tx *query.Query) error) error {
	return r.data.DB.Transaction(fn)
}

