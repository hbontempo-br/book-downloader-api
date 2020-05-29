package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/hbontempo-br/book-downloader-api/api/resources"
	"github.com/hbontempo-br/book-downloader-api/utils"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func SetupRouter(db *gorm.DB, fileStorage utils.MinioFileStorage) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(zap.L(), true))

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Routes
	bookResource := resources.BookResource{DB: db, FileStorage: fileStorage}
	router.GET("/book", bookResource.GetList)
	router.GET("/book/:book_key", bookResource.GetOne)
	router.DELETE("book/:book_key", bookResource.DeleteOne)
	router.GET("book/:book_key/download", bookResource.Download)
	router.POST("/book", bookResource.Create)

	bookStatusResource := resources.BookStatusResource{DB: db}
	router.GET("/book_status", bookStatusResource.GetAll)

	return router
}

func main() {
	// Load env vars
	EnvConfig := utils.LoadEnvVars()

	// Setup logger
	logger, ErrLog := utils.SetupLog(EnvConfig.Environment)
	if ErrLog != nil {
		panic(ErrLog)
	}
	zap.ReplaceGlobals(logger)

	defer func() {
		if logger.Sync() != nil {
			logger.Sugar().Errorw("Error trying to sync logger")
		}
	}()

	// Setup DB
	dbConfig := EnvConfig.DBConfig
	mySQLConnector := utils.NewMySQLConnector(dbConfig.Address, dbConfig.Port, dbConfig.DBName, dbConfig.User, dbConfig.Password)
	db, errBb := mySQLConnector.Connect()
	if errBb != nil {
		panic(errBb)
	}

	// Setup file storage
	minioConfig := EnvConfig.MinioConfig
	fileStorage, errFs := utils.NewMinioFileStorage(minioConfig.Endpoint, minioConfig.AccessKey, minioConfig.SecretKey, minioConfig.SSL)
	if errFs != nil {
		panic(errBb)
	}

	// Server setup
	router := SetupRouter(db, fileStorage)
	port := fmt.Sprintf(":%v", EnvConfig.ServerPort)
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	// Starting server
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			zap.S().Fatalf("listen: %s\n", err)
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.S().Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.S().Fatal("Server forced to shutdown:", err)
	}

	zap.S().Info("Server exiting")
}
