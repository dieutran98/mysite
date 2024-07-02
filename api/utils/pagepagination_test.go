package utils

import (
	"mysite/utils/ptrconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculatePage(t *testing.T) {
	// Test case 1: Positive page number
	{
		expected := 5
		require.Equal(t, expected, calculatePage(expected))
	}

	// Test case 2: Zero page number
	{
		expected := 1
		require.Equal(t, expected, calculatePage(0))
	}

	// Test case 3: Negative page number
	{
		expected := 1
		require.Equal(t, expected, calculatePage(-10))
	}
}

func TestCalculateSize(t *testing.T) {
	// Test case 1: Negative size
	{
		expected := defaultSize
		actual := calculateSize(-5)
		require.Equal(t, expected, actual)
	}

	// Test case 2: Size within limits
	{
		expected := 25
		actual := calculateSize(25)
		require.Equal(t, expected, actual)
	}

	// Test case 3: Size exceeding limit
	{
		expected := limitSize
		actual := calculateSize(120)
		require.Equal(t, expected, actual)
	}
}

func TestCalculateNextPage(t *testing.T) {
	// Test case 1: Valid page within limits
	{
		expected := 3
		page, totalPage := 2, 5
		actual := calculateNextPage(page, totalPage)
		require.Equal(t, expected, *actual)
	}

	// Test case 2: Page 1 with valid total pages
	{
		expected := 2
		page, totalPage := 1, 10
		actual := calculateNextPage(page, totalPage)
		require.Equal(t, expected, *actual)
	}

	// Test case 3: Last page
	{
		page, totalPage := 3, 3 // page == totalPage
		require.Equal(t, (*int)(nil), calculateNextPage(page, totalPage))
	}

	// Test case 4: Invalid page (negative)
	{
		page, totalPage := -1, 5
		actual := calculateNextPage(page, totalPage)
		require.Nil(t, actual)
	}

	// Test case 5: Invalid total pages (less than 2)
	{
		page, totalPage := 2, 1
		actual := calculateNextPage(page, totalPage)
		require.Nil(t, actual)
	}
}

func TestCalculatePreviousPage(t *testing.T) {
	// Test case 1: Valid page with previous page
	{
		expected := 2
		page := 3
		actual := calculatePreviousPage(page)
		require.Equal(t, expected, *actual)
	}

	// Test case 2: Page 1 (no previous page)
	{
		page := 1
		actual := calculatePreviousPage(page)
		require.Nil(t, actual)
	}

	// Test case 3: Negative page (should not affect logic)
	{
		page := -2
		actual := calculatePreviousPage(page)
		require.Nil(t, actual)
	}
}

func TestCalculateTotalPage(t *testing.T) {
	// Test case 1: Divisible records
	{
		expected := 5
		size, totalRecords := 10, 50
		actual := calculateTotalPage(size, totalRecords)
		require.Equal(t, expected, actual)
	}

	// Test case 2: Non-divisible records
	{
		expected := 3
		size, totalRecords := 15, 37
		actual := calculateTotalPage(size, totalRecords)
		require.Equal(t, expected, actual)
	}

	// Test case 3: Zero records
	{
		expected := 0
		size, totalRecords := 10, 0
		actual := calculateTotalPage(size, totalRecords)
		require.Equal(t, expected, actual)
	}

	// Test case 4: Zero size (expected panic)
	{
		size, totalRecords := 0, 10
		require.Panics(t, func() { _ = calculateTotalPage(size, totalRecords) })

	}
}

func TestPageBasePagination(t *testing.T) {
	// Define expected values for a test case (replace with actual data structure)
	expected := &Pagination{
		CurrentPage:   1,
		Previous:      nil,
		RecordPerPage: 10,
		TotalPage:     5,
		Next:          nil,
	}

	// Test case 1: Valid data (page 1, size 10, total records 50)
	{
		expected.Next = ptrconv.Ptr(2)
		page, size, totalRecords := 1, 10, 50
		actual := PageBasePagination(page, size, totalRecords)
		require.Equal(t, expected, actual)
	}

	// Test case 2: Non-first page (page 3, size 15, total records 42)
	{
		expected.CurrentPage = 3
		expected.Previous = ptrconv.Ptr(2)
		expected.RecordPerPage = 15
		expected.TotalPage = 3
		expected.Next = nil

		page, size, totalRecords := 3, 15, 42
		actual := PageBasePagination(page, size, totalRecords)
		require.Equal(t, expected, actual)
	}

	// Test case 3: Zero records (expected defaults for most fields)
	{
		expected.CurrentPage = 1
		expected.Previous = nil
		expected.RecordPerPage = defaultSize
		expected.TotalPage = 0
		expected.Next = nil

		page, size, totalRecords := 2, 20, 0
		actual := PageBasePagination(page, size, totalRecords)
		require.Equal(t, expected, actual)
	}

	// Test case 4: Invalid page (negative) (expected defaults for most fields)
	{
		expected.CurrentPage = 1
		expected.Previous = nil
		expected.RecordPerPage = 10
		expected.TotalPage = 0
		expected.Next = nil

		page, size, totalRecords := -5, 10, 25
		actual := PageBasePagination(page, size, totalRecords)
		require.Equal(t, expected, actual)
	}

	// Test case 5: Invalid size (negative) (expected defaults for most fields)
	{
		expected.CurrentPage = 1
		expected.Previous = nil
		expected.RecordPerPage = defaultSize
		expected.TotalPage = 0
		expected.Next = nil

		page, size, totalRecords := 1, -10, 25
		actual := PageBasePagination(page, size, totalRecords)
		require.Equal(t, expected, actual)
	}
}

func TestToOffsetLimit(t *testing.T) {

	// Test case 1: Valid page and size
	{
		expectedOffset := 20
		expectedLimit := 10
		page, size := 3, 10
		offset, limit := ToOffsetLimit(page, size)
		require.Equal(t, expectedOffset, offset)
		require.Equal(t, expectedLimit, limit)
	}

	// Test case 2: Page 1 with valid size
	{
		expectedOffset := 0
		expectedLimit := 10
		page, size := 1, 10
		offset, limit := ToOffsetLimit(page, size)
		require.Equal(t, expectedOffset, offset)
		require.Equal(t, expectedLimit, limit)
	}

	// Test case 3: Negative page (should return 0, 0)
	{
		expectedOffset := 0
		expectedLimit := 0
		page, size := -2, 15
		offset, limit := ToOffsetLimit(page, size)
		require.Equal(t, expectedOffset, offset)
		require.Equal(t, expectedLimit, limit)
	}

	// Test case 4: Zero size (should use default size)
	{
		expectedOffset := 0
		expectedLimit := 0
		page, size := 3, 0
		offset, limit := ToOffsetLimit(page, size)
		require.Equal(t, expectedOffset, offset)
		require.Equal(t, expectedLimit, limit)
	}
}
