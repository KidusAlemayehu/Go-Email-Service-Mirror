package responses

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type SuccessResponse struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"status_code"`
}

func NewSuccessResponse(message string, data interface{}, statusCode int, w http.ResponseWriter) {
	if data == nil {
		data = make(map[string]interface{})
	}
	response := SuccessResponse{
		Message:    message,
		Data:       data,
		StatusCode: statusCode,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
	StatusCode   int    `json:"status_code"`
}

func (e ErrorResponse) Error() string {
	return e.ErrorMessage
}

func NewErrorResponse(errorMessage error, statusCode int, w http.ResponseWriter) {
	response := ErrorResponse{
		ErrorMessage: errorMessage.Error(),
		StatusCode:   statusCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

type PageData struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int   `json:"total_pages"`
	TotalCount int64 `json:"total_count"`
}

type PaginatedResponse struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	PageData   PageData    `json:"page_data"`
	StatusCode int         `json:"status_code"`
}

func GetPageParams(r *http.Request) (page int, pageSize int) {
	defaultPageSize := 10
	defaultPage := 1

	pageSize = defaultPageSize
	if l, err := strconv.Atoi(r.URL.Query().Get("per_page")); err == nil && l > 0 {
		pageSize = l
	}

	page = defaultPage
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}

	return page, pageSize
}

func NewPageData(page, pageSize int, total int64) PageData {
	pageCount := int((total + int64(pageSize) - 1) / int64(pageSize))
	return PageData{
		Page:       page,
		PerPage:    pageSize,
		TotalCount: int64(total),
		TotalPages: pageCount,
	}
}

func NewPaginatedResponse(message string, data interface{}, statusCode, page, perPage, totalPages int, totalCount int64, w http.ResponseWriter) {
	response := PaginatedResponse{
		Message: message,
		Data:    data,
		PageData: PageData{
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
			TotalCount: totalCount,
		},
		StatusCode: statusCode,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
