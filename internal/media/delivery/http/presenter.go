package http

import (
	"mime/multipart"

	"github.com/pt010104/api-golang/internal/media"
)

type uploadRequest struct {
	Files []*multipart.FileHeader
}

func (r uploadRequest) toInput() media.UploadInput {
	files := make([][]byte, 0, len(r.Files))

	for _, file := range r.Files {
		f, err := file.Open()
		if err != nil {
			continue
		}
		defer f.Close()

		buf := make([]byte, file.Size)
		_, err = f.Read(buf)
		if err != nil {
			continue
		}

		files = append(files, buf)
	}

	return media.UploadInput{
		Files: files,
	}
}
