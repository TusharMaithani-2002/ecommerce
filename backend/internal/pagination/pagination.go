package pagination

import "gorm.io/gorm"

type paginate struct {
	pageSize   int
	pageNumber int
}

func NewPaginate(limit int, page int) *paginate {
	return &paginate{pageSize: limit, pageNumber: page}
}

func (p *paginate) PaginatedResult(db *gorm.DB) *gorm.DB {

	offset := (p.pageNumber - 1) * p.pageSize

	return db.Offset(offset).Limit(p.pageSize)
}