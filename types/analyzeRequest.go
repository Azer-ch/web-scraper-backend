package types

type AnalyzeRequest struct {
	URL string `json:"url" binding:"required"`
}

