package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/hbontempo-br/book-downloader-api/api/DTOs"
	"github.com/hbontempo-br/book-downloader-api/api/controllers"
	"github.com/hbontempo-br/book-downloader-api/utils"
	"github.com/jinzhu/gorm"
	"net/http"
)

type BookStatusResource struct {
	DB *gorm.DB
}

func (bsr BookStatusResource) GetAll(c *gin.Context) {
	tx := bsr.DB.BeginTx(c, nil)
	defer tx.Rollback()

	bookStatusController := cotrollers.NewBookStatusController(tx)
	bookStatus, err := bookStatusController.GetAllStatus()
	if err != nil {
		utils.DefaultErrorMessage(c, http.StatusInternalServerError, nil)
		return
	}
	response := make([]interface{}, 0)
	for _, singleBookStatus := range bookStatus {
		response = append(response, DTOs.NewBookStatusDTO(*singleBookStatus))
	}
	c.JSON(http.StatusOK, response)
}
