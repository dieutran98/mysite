package utils

const (
	defaultSize = 10
	limitSize   = 100
)

// Pagination Object
type Pagination struct {
	Next          *int `json:"nextPage"`
	Previous      *int `json:"previousPage"`
	RecordPerPage int  `json:"recordPerPage"`
	CurrentPage   int  `json:"currentPage"`
	TotalPage     int  `json:"totalPage"`
}

func PageBasePagination(page, size, totalRecords int) *Pagination {
	switch {
	case totalRecords < 1, page < 1, size < 1:
		return &Pagination{
			CurrentPage:   1,
			Next:          nil,
			Previous:      nil,
			RecordPerPage: defaultSize,
			TotalPage:     0,
		}
	}
	p := Pagination{
		CurrentPage:   calculatePage(page),
		Previous:      calculatePreviousPage(page),
		RecordPerPage: calculateSize(size),
		TotalPage:     calculateTotalPage(size, totalRecords),
	}
	p.Next = calculateNextPage(page, p.TotalPage)

	return &p
}

func ToOffsetLimit(page, size int) (offset, limit int) {
	if size < 1 {
		return 0, 0
	}

	if page < 1 {
		return 0, 0
	}
	return (page - 1) * size, size
}

func calculateTotalPage(size, totalRecords int) int {
	totalPage := totalRecords / size
	if totalRecords%size > 0 {
		totalPage = totalPage + 1
	}
	return totalPage
}

func calculatePreviousPage(page int) *int {
	if page < 2 {
		return nil
	}

	previous := page - 1
	return &previous
}

func calculateNextPage(page, totalPage int) *int {
	var next int
	if page < 1 || totalPage < 2 {
		return nil
	}

	next = page + 1
	if next > totalPage {
		return nil
	}

	return &next
}

func calculateSize(size int) int {
	switch {
	case size < 0:
		return defaultSize
	case size > limitSize:
		return limitSize
	}
	return size
}

func calculatePage(page int) int {
	if page < 1 {
		return 1
	}
	return page
}
