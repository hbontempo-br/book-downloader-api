package cotrollers

import (
	"github.com/hbontempo-br/book-downloader-api/api/models"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func NewBookStatusController(transaction *gorm.DB) BookStatusController {
	return BookStatusController{transaction: transaction}
}

type BookStatusController baseController

func (s *BookStatusController) GetAllStatus() ([]*models.BookStatusModel, error) {
	zap.S().Debug("Executing BookStatusController.GetAllStatus")
	var status []*models.BookStatusModel
	if errs := s.transaction.Find(&status).GetErrors(); len(errs) != 0 {
		zap.S().Errorw("Generic error on controller", "errors", errs)
		return nil, ErrGeneric
	}
	return status, nil
}

func (s *BookStatusController) GetStatus(enumerator string) (*models.BookStatusModel, error) {
	zap.S().Debugw("Executing BookStatusController.GetStatus", "enumerator", enumerator)

	var status models.BookStatusModel

	query := s.transaction
	query = filter(query, "enumerator", enumerator)
	resultState := query.First(&status)
	if resultState.RecordNotFound() {
		zap.S().Debugf("Record not found (enumerator=%v)", enumerator)
		return nil, ErrNotFound
	}
	if errs := query.First(&status).GetErrors(); len(errs) != 0 {
		zap.S().Errorw("Generic error on controller", "errors", errs)
		return nil, ErrGeneric
	}
	return &status, nil
}
