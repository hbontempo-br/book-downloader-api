package DTOs

import (
	"net/url"
)

type BookLinkDTO struct {
	DownloadLink string `json:"download_link"`
}

func NewBookLinkDTO(link url.URL) BookLinkDTO {
	return BookLinkDTO{DownloadLink: link.String()}
}
