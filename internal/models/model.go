package models

type Post struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

type UploadImageRequest struct {
	FileName string `json:"file_name"`
	FileData string `json:"file_data"`
}
