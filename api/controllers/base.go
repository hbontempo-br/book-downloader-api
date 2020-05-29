package cotrollers

import (
	"strings"

	"github.com/jinzhu/gorm"
)

const (
	ErrNotFound = ControllerErr("record nof found")
	ErrGeneric  = ControllerErr("generic error")
)

type ControllerErr string

func (e ControllerErr) Error() string {
	return string(e)
}

type baseController struct {
	transaction *gorm.DB // TODO: Change to a interface
}

func addWildcard(s string) string {
	if !strings.HasPrefix(s, "%") {
		s = "%" + s
	}
	if !strings.HasSuffix(s, "%") {
		s += "%"
	}

	return s
}

func filterLike(currentQuery *gorm.DB, column, s string) *gorm.DB {
	s = addWildcard(s)
	newQuery := currentQuery.Where(column+" LIKE ?", s)
	return newQuery
}

func paginatedQuery(currentQuery *gorm.DB, page, pageSize int) *gorm.DB {
	itemsPerPage := pageSize
	offset := (page - 1) * pageSize
	newQuery := currentQuery.Limit(itemsPerPage).Offset(offset)
	return newQuery
}

func filter(currentQuery *gorm.DB, column string, value interface{}) *gorm.DB {
	newQuery := currentQuery.Where(column+" = ?", value)
	return newQuery
}
