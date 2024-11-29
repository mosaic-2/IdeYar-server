package userProfileImpl

import (
	"fmt"
	"github.com/mosaic-2/IdeYar-server/internal/model"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mosaic-2/IdeYar-server/internal/dbutil"
	"github.com/mosaic-2/IdeYar-server/internal/servicers/util"
	"gorm.io/gorm"
)

const UploadDir = "images"

func HandleUserImage(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userID, ok := r.Context().Value(util.UserIDCtxKey{}).(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := os.MkdirAll(UploadDir, os.ModePerm); err != nil {
		http.Error(w, "Cannot create image directory.", http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "The uploaded file is too big.", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("uploadFile")
	if err != nil {
		http.Error(w, "Could not get uploaded file.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		http.Error(w, "The uploaded file must be an image.", http.StatusBadRequest)
		return
	}

	filename := util.GenerateFileName()

	tx := dbutil.GormDB(r.Context())

	err = tx.Transaction(func(tx *gorm.DB) error {
		// Save the file to the filesystem
		dst, err := os.Create(fmt.Sprintf("%s/%s", UploadDir, filename))
		if err != nil {
			http.Error(w, "Could not create a file.", http.StatusInternalServerError)
			return err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Failed to save the uploaded file.", http.StatusInternalServerError)
			return err
		}

		if err := tx.Model(&model.User{}).Where("id = ?", userID).Update("profile_image_url", filename).Error; err != nil {
			http.Error(w, "Failed to update profile image URL.", http.StatusInternalServerError)
			return err
		}

		w.WriteHeader(http.StatusOK)
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
