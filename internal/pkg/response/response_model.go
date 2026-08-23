package response

type WebResponse[T any] struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    T         `json:"data"`
	Meta    *Metadata `json:"meta,omitempty"`
	Token   string    `json:"token,omitempty"`
}

type Metadata struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
}

type WebErrorResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Details []ErrorDetails `json:"details,omitempty"`
}

type ErrorDetails struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
