package utils

const (
	InvalidPageError          = PaginationError("invalid page number, page number can not be negative")
	InvalidPageSizeError      = PaginationError("invalid page size, page size can not be negative")
	InvalidTotalCountError    = PaginationError("invalid total count, total count can not be negative")
	PageSizeDataMismatchError = PaginationError("informed page size smaller than the data provided")
)

type PaginationError string

func (e PaginationError) Error() string {
	return string(e)
}

type Pagination struct {
	CurrentPage  int  `json:"current_page"`
	NextPage     *int `json:"next_page"`
	PreviousPage *int `json:"previous_page"`
	MaxPage      int  `json:"max_page"`
	RowsPerPage  int  `json:"rows_per_page"`
	TotalRows    int  `json:"total_rows"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"Pagination"`
}

func FormatPaginatedResponse(data []interface{}, pageSize, page, totalCount int) (*PaginatedResponse, error) {
	if page < 1 {
		return nil, InvalidPageError
	}
	if pageSize < 1 {
		return nil, InvalidPageSizeError
	}
	if totalCount < 0 {
		return nil, InvalidTotalCountError
	}
	if pageSize < len(data) {
		return nil, PageSizeDataMismatchError
	}
	maxPage := ((totalCount - 1) / pageSize) + 1
	var nextPage *int
	if maxPage > page {
		np := page + 1
		nextPage = &np
	}
	var prevPage *int
	if page != 1 && page-1 <= maxPage {
		pp := page - 1
		prevPage = &pp
	}
	pag := Pagination{CurrentPage: page, NextPage: nextPage, PreviousPage: prevPage, MaxPage: maxPage, RowsPerPage: pageSize, TotalRows: totalCount}
	resp := PaginatedResponse{Data: data, Pagination: pag}
	return &resp, nil
}
