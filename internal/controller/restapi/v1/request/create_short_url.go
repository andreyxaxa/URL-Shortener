package request

type CreateShortURLRequest struct {
	URL         string `json:"url"`
	CustomAlias string `json:"custom_alias,omitempty"`
}
