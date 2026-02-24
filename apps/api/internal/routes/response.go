package routes

type ApiError struct {
	Code       string           `json:"code,omitempty"`
	Target     *string          `json:"target,omitempty"`
	Message    string           `json:"message,omitempty"`
	InnerError *ApiError        `json:"innerError,omitempty"`
	Details    []ApiErrorDetail `json:"details,omitempty"`
}

type ApiErrorDetail struct {
	Code    *string `json:"code,omitempty"`
	Message string  `json:"message,omitempty"`
	Target  string  `json:"target,omitempty"`
}

type Response[T any] struct {
	Ok       bool      `json:"ok"`
	Error    *ApiError `json:"error,omitempty"`
	TraceId  *string   `json:"traceId,omitempty"`
	SpanId   *string   `json:"spanId,omitempty"`
	Value    T         `json:"value,omitempty"`
	NextLink *string   `json:"@nextLink,omitempty"`
}

func SuccessResponse[T any](value T) Response[T] {
	return Response[T]{
		Ok:    true,
		Value: value,
	}
}

func ErrorResponse(err *ApiError) Response[any] {
	return Response[any]{
		Ok:    false,
		Error: err,
	}
}
