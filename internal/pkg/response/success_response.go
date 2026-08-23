package response

import "github.com/gin-gonic/gin"

type Option[T any] func(*WebResponse[T])

func WithMeta[T any](meta *Metadata) Option[T] {
	return func(r *WebResponse[T]) {
		r.Meta = meta
	}
}

func WithToken[T any](token string) Option[T] {
	return func(r *WebResponse[T]) {
		r.Token = token
	}
}

func SuccessResponse[T any](c *gin.Context, code int, message string, data T, opts ...Option[T]) {
	res := &WebResponse[T]{
		Success: true,
		Message: message,
		Data:    data,
	}

	for _, opt := range opts {
		opt(res)
	}

	c.JSON(code, res)
}
