package helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Pagination struct {
	Page int
	PageSize int
	Total int64
	WithoutCount bool
}
func (p *Pagination) Build(db *gorm.DB, c *gin.Context) *gorm.DB {
	if c != nil {
		p.Page, _ = strconv.Atoi(c.DefaultQuery("pagination[page]", "1"))	
		p.PageSize, _ = strconv.Atoi(c.DefaultQuery("pagination[pageSize]", "20"))	
		p.WithoutCount, _ = strconv.ParseBool(c.DefaultQuery("pagination[withoutCount]", "0"))
	}

	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 30
	}
	// 查询总数
	if !p.WithoutCount {
		db.Count(&p.Total)
	}
	return db.Limit(p.PageSize).Offset(p.PageSize * (p.Page - 1))
}