package dto

type Link struct {
	OriginLink string `json:"origin_link" binding:"required" `
	Slug       string `json:"slug" `
}
