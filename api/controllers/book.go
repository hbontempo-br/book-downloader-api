package cotrollers

import (
	uuid "github.com/gofrs/uuid"
	"github.com/hbontempo-br/book-downloader-api/api/models"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func NewBookController(transaction *gorm.DB) BookController {
	return BookController{transaction: transaction}
}

type BookController baseController

func (bc *BookController) Delete(bookModel *models.BookModel) error {
	zap.S().Debugw("Executing BookController.Delete", "bookModel", bookModel.String())
	if errSlice := bc.transaction.Delete(&bookModel).GetErrors(); len(errSlice) != 0 {
		zap.S().Errorw("Generic error on controller", "errors", errSlice)
		return ErrGeneric
	}

	return nil
}

func (bc *BookController) GetBook(bookKey string) (*models.BookModel, error) {
	zap.S().Debugw("Executing BookController.GetBook", "bookKey", bookKey)
	query := bc.transaction.Model(&models.BookModel{}).Preload("Status")
	query = filter(query, "book_key", bookKey)

	var book models.BookModel
	resultQuery := query.First(&book)
	if resultQuery.RecordNotFound() {
		zap.S().Infow("Book not found", "bookKey", bookKey)
		return nil, ErrNotFound
	}
	if errSlice := resultQuery.GetErrors(); len(errSlice) != 0 {
		zap.S().Errorw("Generic error on controller", "errors", errSlice)
		return nil, ErrGeneric
	}

	return &book, nil
}

func (bc *BookController) GetBooks(nameLike string, page, pageSize int) ([]*models.BookModel, int, error) {
	zap.S().Debugw("Executing BookController.GetBooks", "nameLike", nameLike, "page", page, "pageSize", pageSize)
	query := bc.transaction.Model(&models.BookModel{}).Preload("Status")

	if nameLike != "" {
		query = bc.filterNameLike(query, nameLike)
	}

	var count int
	if errSlice := query.Count(&count).GetErrors(); len(errSlice) != 0 {
		zap.S().Errorw("Generic error on controller", "errors", errSlice)
		return nil, 0, ErrGeneric
	}

	query = query.Order("created_at desc")

	var books []*models.BookModel
	if errSlice := paginatedQuery(query, page, pageSize).Find(&books).GetErrors(); len(errSlice) != 0 {
		zap.S().Errorw("Generic error on controller", "errors", errSlice)
		return nil, 0, ErrGeneric
	}
	return books, count, nil
}

func (bc *BookController) filterNameLike(currentQuery *gorm.DB, name string) *gorm.DB {
	return filterLike(currentQuery, "name", name)
}

func (bc *BookController) Create(name, mask, status string) (*models.BookModel, error) {
	zap.S().Debugw("Executing BookController.Create", "name", name, "mask", mask, "status", status)

	bookKey, uuidErr := uuid.NewV4()
	if uuidErr != nil {
		zap.S().Errorw("Generic error on controller", "errors", uuidErr)
		return nil, ErrGeneric
	}

	bookStatusController := NewBookStatusController(bc.transaction)
	pendingStatus, _ := bookStatusController.GetStatus(status)
	book := models.BookModel{
		BookKey: bookKey,
		Name:    name,
		Mask:    mask,
		Status:  *pendingStatus,
	}
	if errSlice := bc.transaction.Create(&book).GetErrors(); len(errSlice) != 0 {
		zap.S().Errorw("Generic error on controller", "errors", errSlice)
		return nil, ErrGeneric
	}
	return &book, nil
}

func (bc *BookController) Update(bookModel *models.BookModel, name, mask, status string) error {
	zap.S().Debugw("Executing BookController.Update", "bookModel", bookModel, "name", name, "mask", mask, "status", status)
	bookStatusController := NewBookStatusController(bc.transaction)
	newStatus, statusErr := bookStatusController.GetStatus(status)
	if statusErr != nil {
		zap.S().Errorw("Generic error on controller", "errors", statusErr)
		return ErrGeneric
	}
	resultState := bc.transaction.Model(&bookModel).Updates(models.BookModel{Name: name, Mask: mask, Status: *newStatus, StatusID: newStatus.ID})
	if errSlice := resultState.GetErrors(); len(errSlice) != 0 {
		zap.S().Errorw("Generic error on controller", "errors", errSlice)
		return ErrGeneric
	}
	return nil
}
