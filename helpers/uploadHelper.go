package helper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func SingleImageUpload(w http.ResponseWriter, r *http.Request, avatar string) (string, error) {

	file, fileHeader, err := r.FormFile(avatar)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return "", err
	}

	defer file.Close()

	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	avatarName := time.Now().UnixNano()
	avatarExtention := filepath.Ext(fileHeader.Filename)

	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", avatarName, avatarExtention))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	return fmt.Sprintf("%d%s", avatarName, avatarExtention), nil

}
