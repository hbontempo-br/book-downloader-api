package utils

import (
	"bytes"
	"io"
	"net/http"

	"github.com/signintech/gopdf"
	"go.uber.org/zap"
)

const (
	ErrDownloadFailed      = DictError("error on image download")
	ErrCloseReqBodyFailed  = DictError("error on closing a requests response")
	ErrPDFGenerationFailed = DictError("error during pds generation")
)

type DictError string

func (e DictError) Error() string {
	return string(e)
}

func NewBookDownloader(urlMaskFunc func(int) string) BookDownloader {
	return BookDownloader{urlMaskFunc: urlMaskFunc, retry: 5} // TODO: remove this magic number, use environment variable
}

type BookDownloader struct {
	urlMaskFunc func(int) string
	pages       []io.ReadWriter
	retry       int
	book        bytes.Buffer
}

func (bd *BookDownloader) downloadPage(url string, retry int) (response bytes.Buffer, finished bool, err error) {
	zap.S().Debugf("Starting download (url: %v, retry: %v)", url, retry)

	finished = false

	client := func() *http.Client {
		client := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}

		return &client
	}()

	resp, errReq := client.Get(url)
	if errReq == nil {
		defer func() {
			if errReqClose := resp.Body.Close(); errReqClose != nil {
				zap.S().Errorw("Generic error on closing request body", "errors", errReqClose)
				err = ErrCloseReqBodyFailed
			}
		}()
	}

	switch {
	case errReq == nil && resp.StatusCode == http.StatusOK:
		if _, errCopy := io.Copy(&response, resp.Body); errCopy != nil {
			zap.S().Errorw("Generic error on coping request body", "errors", errCopy)
			err = ErrCloseReqBodyFailed
		}
		zap.S().Debugf("Successfully downloaded page from %v", url)
	case errReq == nil && resp.StatusCode == http.StatusForbidden:
		zap.S().Infof("Page forbidden -> book download finished")
		finished = true
	case retry > 0:
		zap.S().Warnf("Retrying download %v (%v retries available)\n", url, retry)
		return bd.downloadPage(url, retry-1)
	default:
		zap.S().Warnf("Download failed %v (%v retries available)\n", url, retry)
		err = ErrDownloadFailed
	}
	return
}

func (bd *BookDownloader) downloadAllPages() error {
	zap.S().Info("Starting book download")

	pageNumber := 1
	for {
		url := bd.urlMaskFunc(pageNumber)
		zap.S().Debugw("Starting page download", "pageNumber", pageNumber, "url", url)

		imgReader, finished, err := bd.downloadPage(url, bd.retry)
		if err != nil {
			return err
		}

		if finished {
			break
		}

		bd.pages = append(bd.pages, &imgReader)

		pageNumber++
	}

	zap.S().Debugw("Finished book download", "numberOfPages", pageNumber-1)
	return nil
}

func (bd *BookDownloader) generatePDF() (err error) {
	zap.S().Info("Starting book pdf generation")

	pdf := gopdf.GoPdf{}
	defer func() {
		if errPDFClose := pdf.Close(); errPDFClose != nil {
			zap.S().Errorw("Error closing pdf", "errors", errPDFClose)
			err = errPDFClose
		}
	}()

	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	for _, download := range bd.pages {
		imageHolder, _ := gopdf.ImageHolderByReader(download)
		pdf.AddPage()

		if errPDFImg := pdf.ImageByHolder(imageHolder, 0, 0, gopdf.PageSizeA4); errPDFImg != nil {
			zap.S().Errorw("Error creating new page", "errors", errPDFImg)
			return ErrPDFGenerationFailed
		}
	}

	if errPDFWrite := pdf.Write(&bd.book); errPDFWrite != nil {
		zap.S().Errorw("Error on pdf generation", "errors", errPDFWrite)
		return ErrPDFGenerationFailed
	}

	zap.S().Info("Successfully finished book pdf generation")
	return nil
}

func (bd *BookDownloader) CreatePDF() (io.Reader, error) {
	if err := bd.downloadAllPages(); err != nil {
		return nil, err
	}
	if err := bd.generatePDF(); err != nil {
		return nil, err
	}
	return &bd.book, nil
}
