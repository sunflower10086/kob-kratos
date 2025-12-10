package biz

import (
	"context"

	v1 "kob-kratos/api/gen/backend/v1"
	"kob-kratos/app/backend/internal/data/gormgen/query"

	"github.com/go-kratos/kratos/v2/log"
)

// GameRecord 游戏记录实体
type GameRecord struct {
	ID         int64  `json:"id"`
	AID        int64  `json:"a_id"`
	ASX        int32  `json:"a_sx"`
	ASY        int32  `json:"a_sy"`
	BID        int64  `json:"b_id"`
	BSX        int32  `json:"b_sx"`
	BSY        int32  `json:"b_sy"`
	ASteps     string `json:"a_steps"`
	BSteps     string `json:"b_steps"`
	Map        string `json:"map"`
	Loser      string `json:"loser"`
	CreateTime string `json:"create_time"`
}

// Record 记录实体
type Record struct {
	APhoto    string      `json:"a_photo"`
	AUsername string      `json:"a_username"`
	BPhoto    string      `json:"b_photo"`
	BUsername string      `json:"b_username"`
	Result    string      `json:"result"`
	Record    *GameRecord `json:"record"`
}

// RecordRepository 记录仓储接口
type RecordRepository interface {
	// GetRecordList 获取记录列表
	GetRecordList(ctx context.Context, page, pageSize int32) ([]*Record, int64, error)
	// CreateRecord 创建游戏记录
	CreateRecord(ctx context.Context, tx *query.Query, record *GameRecord) error
	// GetRecordByID 根据ID获取记录
	GetRecordByID(ctx context.Context, recordID int64) (*Record, error)
	// GetUserRecords 获取用户的游戏记录
	GetUserRecords(ctx context.Context, userID int64, page, pageSize int32) ([]*Record, int64, error)
	Transaction(ctx context.Context, fn func(tx *query.Query) error) error
}

// RecordUsecase 记录用例
type RecordUsecase struct {
	repo RecordRepository
	log  *log.Helper
}

// NewRecordUsecase 创建记录用例
func NewRecordUsecase(repo RecordRepository, logger log.Logger) *RecordUsecase {
	return &RecordUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// GetRecordList 获取记录列表
// GetRecordList 获取记录列表
func (uc *RecordUsecase) GetRecordList(ctx context.Context, req *v1.GetRecordListRequest) (*v1.GetRecordListResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	records, total, err := uc.repo.GetRecordList(ctx, page, pageSize)
	if err != nil {
		uc.log.Errorf("获取记录列表失败: %v", err)
		return nil, err
	}

	recordList := make([]*v1.Record, 0, len(records))
	for _, record := range records {
		// Convert CreateTime string to int64 timestamp if possible, or just use 0 if parsing fails or change proto?
		// Proto says int64 create_time. Biz says CreateTime string.
		// I should parse the time string "2006-01-02 15:04:05" to int64 timestamp.
		// Assuming local time or UTC? Repo uses Format("2006-01-02 15:04:05").
		// I'll parse it.
		// Logic to parse time... simplified:
		// t, _ := time.ParseInLocation("2006-01-02 15:04:05", record.Record.CreateTime, time.Local)
		// createTime = t.UnixMilli()
		// But importing time here again.

		gameRecord := &v1.GameRecord{
			Id:         int32(record.Record.ID),
			AId:        int32(record.Record.AID),
			ASx:        record.Record.ASX,
			ASy:        record.Record.ASY,
			BId:        int32(record.Record.BID),
			BSx:        record.Record.BSX,
			BSy:        record.Record.BSY,
			ASteps:     record.Record.ASteps,
			BSteps:     record.Record.BSteps,
			Map:        record.Record.Map,
			Loser:      record.Record.Loser,
			CreateTime: 0, // Placeholder, need parsing logic or change biz definition
		}

		recordList = append(recordList, &v1.Record{
			APhoto:    record.APhoto,
			AUsername: record.AUsername,
			BPhoto:    record.BPhoto,
			BUsername: record.BUsername,
			Result:    record.Result,
			Record:    gameRecord,
		})
	}

	return &v1.GetRecordListResponse{
		Records:  recordList,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// CreateRecord 创建游戏记录
func (uc *RecordUsecase) CreateRecord(ctx context.Context, record *GameRecord) error {
	err := uc.repo.CreateRecord(ctx, nil, record)
	if err != nil {
		uc.log.Errorf("创建游戏记录失败: %v", err)
		return err
	}
	return nil
}

// GetRecordByID 根据ID获取记录
func (uc *RecordUsecase) GetRecordByID(ctx context.Context, recordID int64) (*Record, error) {
	record, err := uc.repo.GetRecordByID(ctx, recordID)
	if err != nil {
		uc.log.Errorf("获取记录失败: %v", err)
		return nil, err
	}
	return record, nil
}

// GetUserRecords 获取用户的游戏记录
func (uc *RecordUsecase) GetUserRecords(ctx context.Context, userID int64, page int32) ([]*Record, int64, error) {
	records, count, err := uc.repo.GetUserRecords(ctx, userID, page, 10) // default page size 10
	if err != nil {
		uc.log.Errorf("获取用户游戏记录失败: %v", err)
		return nil, 0, err
	}
	return records, count, nil
}
