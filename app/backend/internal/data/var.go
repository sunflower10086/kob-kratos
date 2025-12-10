package data

import (
	"gorm.io/gen"
)

func Paginate(pageNum, pageSize int) func(dao gen.Dao) gen.Dao {
	return func(dao gen.Dao) gen.Dao {
		if pageNum <= 0 {
			pageNum = 1
		}

		offset := (pageNum - 1) * pageSize
		return dao.Offset(offset).Limit(pageSize)
	}
}
