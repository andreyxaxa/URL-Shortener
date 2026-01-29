package response

type CreateShortURLResponse struct {
	URL      string `json:"original_url"`
	ShortURL string `json:"short_url"`
}
