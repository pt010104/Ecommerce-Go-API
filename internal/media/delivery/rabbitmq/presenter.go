package rabbitmq

type UploadMessage struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	ShopID     string `json:"shop_id"`
	FileName   string `json:"file_name"`
	File       []byte `json:"file"`
	FolderName string `json:"folder_name"`
}
