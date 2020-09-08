package resources

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hbontempo-br/book-downloader-api/config"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/hbontempo-br/book-downloader-api/api/DTOs"
	controllers "github.com/hbontempo-br/book-downloader-api/api/controllers"
	"github.com/hbontempo-br/book-downloader-api/api/models"
	"github.com/hbontempo-br/book-downloader-api/utils"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type CreateBookInput struct {
	Name string `json:"name" biding:"required"`
	Mask string `json:"mask" biding:"required"`
}

type GetBookQuery struct {
	Name     string `form:"name"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type GetBookLinkQuery struct {
	Expiry int `form:"expiry" biding:"numeric,max=300"`
}

type BookResource struct {
	DB          *gorm.DB
	FileStorage utils.MinioFileStorage
}

func (br *BookResource) GetOne(c *gin.Context) {
	tx := br.DB.BeginTx(c, nil)
	defer tx.Rollback()

	// Retrieve book_key from request
	bookKey := c.Param("book_key")

	// Search books
	bookController := controllers.NewBookController(tx)
	book, err := bookController.GetBook(bookKey)
	if err == controllers.ErrNotFound {
		utils.DefaultErrorMessage(c, http.StatusNotFound, "Unable to retrieve book with given key")
		return
	} else if err != nil {
		utils.DefaultErrorMessage(c, http.StatusInternalServerError, nil)
		return
	}

	// Format response
	response := DTOs.NewBookDTO(*book)

	// Success response
	c.PureJSON(http.StatusOK, response)
}

func (br *BookResource) DeleteOne(c *gin.Context) {
	tx := br.DB.BeginTx(c, nil)
	defer tx.Rollback()

	// Retrieve book_key from request
	bookKey := c.Param("book_key")

	// Search book
	bookController := controllers.NewBookController(tx)
	book, err := bookController.GetBook(bookKey)
	if err == controllers.ErrNotFound {
		utils.DefaultErrorMessage(c, http.StatusBadRequest, "Can't delete a non-existent book")
		return
	} else if err != nil {
		utils.DefaultErrorMessage(c, http.StatusInternalServerError, nil)
		return
	}

	// Delete book
	if err := bookController.Delete(book); err != nil {
		utils.DefaultErrorMessage(c, http.StatusInternalServerError, nil)
		return
	}

	// Commit
	tx.Commit()

	// Success response
	c.PureJSON(http.StatusNoContent, nil)
}

func (br *BookResource) GetList(c *gin.Context) {
	tx := br.DB.BeginTx(c, nil)
	defer tx.Rollback()

	// Load query string
	bookQuery := GetBookQuery{Page: 1, PageSize: 10} // TODO: remove this magic number, use environment variable
	if err := c.ShouldBindQuery(&bookQuery); err != nil {
		utils.DefaultErrorMessage(c, http.StatusBadRequest, err)
		return
	}

	// Search books
	bookController := controllers.NewBookController(tx)
	books, totalCount, err := bookController.GetBooks(bookQuery.Name, bookQuery.Page, bookQuery.PageSize)
	if err != nil {
		utils.DefaultErrorMessage(c, http.StatusInternalServerError, nil)
		return
	}

	// Format response
	data := make([]interface{}, 0)
	for _, book := range books {
		data = append(data, DTOs.NewBookDTO(*book))
	}
	response := utils.FormatPaginatedResponse(data, bookQuery.PageSize, bookQuery.Page, totalCount)

	// Success response
	c.PureJSON(http.StatusOK, response)
}

func (br *BookResource) Create(c *gin.Context) {
	tx := br.DB.BeginTx(c, nil)
	defer tx.Rollback()

	// Validate input
	var input CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.DefaultErrorMessage(c, http.StatusBadRequest, err)
		return
	}

	// TODO: Add validations, specially regarding book`s mask (a regex or something) and file extension

	bookController := controllers.NewBookController(tx)
	book, err := bookController.Create(input.Name, input.Mask, "pending")
	if err != nil {
		utils.DefaultErrorMessage(c, http.StatusBadRequest, nil)
		return
	}

	// Format response
	response := DTOs.NewBookDTO(*book)

	tx.Commit()

	// Success response
	c.PureJSON(http.StatusAccepted, response)

	// Start download routine
	go br.downloadRoutine(c.Copy(), book)
}

func (br *BookResource) Download(c *gin.Context) {
	tx := br.DB.BeginTx(c, nil)
	defer tx.Rollback()

	// Retrieve book_key from request
	bookKey := c.Param("book_key")

	// Search books
	bookController := controllers.NewBookController(tx)
	book, err := bookController.GetBook(bookKey)
	if err == controllers.ErrNotFound {
		utils.DefaultErrorMessage(c, http.StatusBadRequest, "Can't download a non-existent book")
		return
	} else if err != nil {
		utils.DefaultErrorMessage(c, http.StatusInternalServerError, nil)
		return
	}

	// TODO: check if status is valid before trying to download

	var bucketConfig config.BucketConfig
	if env, err := config.LoadEnvVars(); err != nil {
		// TODO: handle error
	} else {
		bucketConfig = env.BucketConfig
	}
	pdf, _ := br.FileStorage.Get(bucketConfig.Name, fmt.Sprintf("%v/%v", book.BookKey, book.Name)) // TODO: check for error

	if err := formatDownloadResponse(pdf, book.Name, c); err != nil {
		utils.DefaultErrorMessage(c, http.StatusInternalServerError, nil)
		return
	}
}

func (br *BookResource) DownloadLink(c *gin.Context) {
	tx := br.DB.BeginTx(c, nil)
	defer tx.Rollback()

	// Retrieve book_key from request
	bookKey := c.Param("book_key")

	// Load query string
	bookLinkQuery := GetBookLinkQuery{Expiry: 600} // TODO: remove this magic number, use environment variable
	if err := c.ShouldBindQuery(&bookLinkQuery); err != nil {
		utils.DefaultErrorMessage(c, http.StatusBadRequest, err)
		return
	}
	expiryDuration := time.Duration(bookLinkQuery.Expiry) * time.Second

	// Search books
	bookController := controllers.NewBookController(tx)
	book, err := bookController.GetBook(bookKey)
	if err == controllers.ErrNotFound {
		utils.DefaultErrorMessage(c, http.StatusBadRequest, "Can't download a non-existent book")
		return
	} else if err != nil {
		utils.DefaultErrorMessage(c, http.StatusInternalServerError, nil)
		return
	}

	var bucketConfig config.BucketConfig
	if env, err := config.LoadEnvVars(); err != nil {
		// TODO: handle error
	} else {
		bucketConfig = env.BucketConfig
	}
	bookLink, _ := br.FileStorage.GetLink(book.Name, bucketConfig.Name, fmt.Sprintf("%v/%v", book.BookKey, book.Name), expiryDuration)
	// TODO: check if status is valid before trying to download

	// Format response
	response := DTOs.NewBookLinkDTO(*bookLink)

	// Success response
	c.PureJSON(http.StatusOK, response)

}

func (br *BookResource) downloadRoutine(c context.Context, book *models.BookModel) {
	tx := br.DB.BeginTx(c, nil)
	defer tx.Rollback()

	urlMaskFunc := func(i int) string {
		s := strings.Split(book.Mask, "{page_number}")
		page := strconv.Itoa(i)
		return s[0] + page + s[1]
	}
	bookDownloader := utils.NewBookDownloader(urlMaskFunc)
	pdf, err := bookDownloader.CreatePDF()
	if err != nil {
		zap.S().Errorw("Unable to create PDF file", "errors", err)
		return
	}

	var bucketConfig config.BucketConfig
	if env, err := config.LoadEnvVars(); err != nil {
		// TODO: handle error
	} else {
		bucketConfig = env.BucketConfig
	}
	if err := br.FileStorage.Save(pdf, bucketConfig.Name, fmt.Sprintf("%v/%v", book.BookKey, book.Name)); err != nil {
		zap.S().Errorw("Unable to save book to file storage", "errors", err)
		return
	}

	bookController := controllers.NewBookController(tx)
	if err := bookController.Update(book, "", "", "finished"); err != nil {
		zap.S().Errorw("Unable to update book database record", "errors", err)
		return
	}

	tx.Commit()
}

func formatDownloadResponse(file io.Reader, filename string, c *gin.Context) error {
	var buffer bytes.Buffer
	contentLength, copyErr := io.Copy(&buffer, file)
	if copyErr != nil {
		zap.S().Errorw("Error copying downloaded file")
		return copyErr
	}

	responseFile := bytes.NewReader(buffer.Bytes())
	mt, err := mimetype.DetectReader(responseFile)
	if err != nil {
		zap.S().Errorw("Error detecting file`s mimetype")
		return err
	}
	contentType := mt.String()

	if _, err := responseFile.Seek(0, 0); err != nil {
		zap.S().Errorw("Internal error")
		return copyErr
	}

	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%v"`, filename),
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, responseFile, extraHeaders)

	return nil
}
