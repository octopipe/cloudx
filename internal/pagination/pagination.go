package pagination

type SimplePaginationRequest struct {
	Limit int64  `json:"limit"`
	Page  string `json:"page"`
}

type SimplePaginationResponse[T any] struct {
	Items []T   `json:"items"`
	Total int32 `json:"total"`
}

type ChunkingPaginationRequest struct {
	Limit int64  `json:"limit"`
	Chunk string `json:"chunk"`
}

type ChunkingPaginationResponse[T any] struct {
	Items []T    `json:"items"`
	Chunk string `json:"chunk"`
	Total int32  `json:"total"`
}
