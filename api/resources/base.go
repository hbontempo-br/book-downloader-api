package resources

type pagination struct {
	CurrentPage  int  `json:"current_page"`
	NextPage     *int `json:"next_page"`
	PreviousPage *int `json:"previous_page"`
	MaxPage      int  `json:"max_page"`
	RowsPerPage  int  `json:"rows_per_page"`
	TotalRows    int  `json:"total_rows"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination pagination  `json:"pagination"`
}

func formatPaginatedResponse(data []interface{}, pageSize, page, totalCount int) PaginatedResponse {
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
	pag := pagination{CurrentPage: page, NextPage: nextPage, PreviousPage: prevPage, MaxPage: maxPage, RowsPerPage: pageSize, TotalRows: totalCount}
	resp := PaginatedResponse{Data: data, Pagination: pag}
	return resp
}
