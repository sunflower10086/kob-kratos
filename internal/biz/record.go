package biz

import (
	"context"

	v1 "kob-kratos/api/backend/v1"
	"kob-kratos/internal/data/gormgen/query"

	"github.com/go-kratos/kratos/v2/log"
)

// GameRecord 游戏记录实体
type GameRecord struct {
	ID         int32  `json:"id"`
	AID        int32  `json:"a_id"`
	ASX        int32  `json:"a_sx"`
	ASY        int32  `json:"a_sy"`
	BID        int32  `json:"b_id"`
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
	GetRecordByID(ctx context.Context, recordID int32) (*Record, error)
	// GetUserRecords 获取用户的游戏记录
	GetUserRecords(ctx context.Context, userID int32, page, pageSize int32) ([]*Record, int64, error)
	Transaction(ctx context.Context, fn func(tx *query.Query) error) error
}

// RecordUsecase 记录用例
type RecordUsecase struct {
	repo RecordRepository
	log  *log.Helper
}

// NewRecordUsecase 创建记录用例
func NewRecordUsecase(logger log.Logger) *RecordUsecase {
	return &RecordUsecase{
		// repo: repo,
		log: log.NewHelper(logger),
	}
}

// GetRecordList 获取记录列表
func (uc *RecordUsecase) GetRecordList(ctx context.Context, req *v1.GetRecordListRequest) (*v1.GetRecordListResponse, error) {
	// page := parseStringToInt32(req.Page)
	// if page <= 0 {
	// 	page = 1
	// }

	// records, recordsCount, err := uc.repo.GetRecordList(ctx, page)
	// if err != nil {
	// 	uc.log.Errorf("获取记录列表失败: %v", err)
	// 	return nil, err
	// }

	// recordList := make([]*v1.Record, 0, len(records))
	// for _, record := range records {
	// 	gameRecord := &v1.GameRecord{
	// 		Id:         record.Record.ID,
	// 		AId:        record.Record.AID,
	// 		ASx:        record.Record.ASX,
	// 		ASy:        record.Record.ASY,
	// 		BId:        record.Record.BID,
	// 		BSx:        record.Record.BSX,
	// 		BSy:        record.Record.BSY,
	// 		ASteps:     record.Record.ASteps,
	// 		BSteps:     record.Record.BSteps,
	// 		Map:        record.Record.Map,
	// 		Loser:      record.Record.Loser,
	// 		CreateTime: record.Record.CreateTime,
	// 	}

	// 	recordList = append(recordList, &v1.Record{
	// 		APhoto:    record.APhoto,
	// 		AUsername: record.AUsername,
	// 		BPhoto:    record.BPhoto,
	// 		BUsername: record.BUsername,
	// 		Result:    record.Result,
	// 		Record:    gameRecord,
	// 	})
	// }

	// return &v1.GetRecordListResponse{
	// 	Records:      recordList,
	// 	RecordsCount: recordsCount,
	// }, nil
	panic("implement me")
}

// CreateRecord 创建游戏记录
func (uc *RecordUsecase) CreateRecord(ctx context.Context, record *GameRecord) error {
	// err := uc.repo.CreateRecord(ctx, record)
	// if err != nil {
	// 	uc.log.Errorf("创建游戏记录失败: %v", err)
	// 	return err
	// }
	// return nil
	panic("implement me")
}

// GetRecordByID 根据ID获取记录
func (uc *RecordUsecase) GetRecordByID(ctx context.Context, recordID int32) (*Record, error) {
	// record, err := uc.repo.GetRecordByID(ctx, recordID)
	// if err != nil {
	// 	uc.log.Errorf("获取记录失败: %v", err)
	// 	return nil, err
	// }
	// return record, nil
	panic("implement me")
}

// GetUserRecords 获取用户的游戏记录
func (uc *RecordUsecase) GetUserRecords(ctx context.Context, userID int32, page int32) ([]*Record, int64, error) {
	// records, count, err := uc.repo.GetUserRecords(ctx, userID, page)
	// if err != nil {
	// 	uc.log.Errorf("获取用户游戏记录失败: %v", err)
	// 	return nil, 0, err
	// }
	// return records, count, nil
	panic("implement me")
}
