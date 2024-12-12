package media

type ConsumeUploadMsgInput struct {
	ID         string
	UserID     string
	ShopID     string
	FolderName string
	FileName   string
	File       []byte
}
