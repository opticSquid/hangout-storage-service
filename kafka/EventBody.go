package kafka

type eventBody struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	UserId      int32  `json:"userId"`
}
