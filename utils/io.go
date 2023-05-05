package utils

import (
	"io"
	"mime/multipart"
	"os"
)

func SaveToLocalFile(fileLocation string, file *multipart.File) error {
	out, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		errClose := out.Close()
		if errClose != nil {
			err = errClose
		}
	}(out)
	_, err = io.Copy(out, *file)
	return err
}
