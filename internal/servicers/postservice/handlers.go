package postservice

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mosaic-2/IdeYar-server/internal/dbutil"
	"github.com/mosaic-2/IdeYar-server/internal/model"
	"github.com/mosaic-2/IdeYar-server/internal/servicers/util"
	"gorm.io/gorm"
)

const UploadDir = "images"

func HandlePostImage(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userID, ok := r.Context().Value(util.UserIDCtxKey{}).(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form, limiting file size to 10MB
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

	postIDStr := r.PostFormValue("postID")
	if postIDStr == "" {
		http.Error(w, "postID is required.", http.StatusBadRequest)
		return
	}
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid postID format.", http.StatusBadRequest)
		return
	}

	orderStr := r.PostFormValue("order")
	if orderStr == "" {
		http.Error(w, "order is required.", http.StatusBadRequest)
		return
	}
	order, err := strconv.ParseInt(orderStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid order format.", http.StatusBadRequest)
		return
	}

	filename := util.GenerateFileName()

	tx := dbutil.GormDB(r.Context())

	err = tx.Transaction(func(tx *gorm.DB) error {
		// validations
		if order < 0 || order > 9 {
			http.Error(w, "", http.StatusBadRequest)
			return errors.New("invalid order")
		}

		var ownerUserID int64

		err = tx.Table("post").
			Where("id = ?", postID).
			Select("user_id").
			Scan(&ownerUserID).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		if userID != ownerUserID {
			w.WriteHeader(http.StatusUnauthorized)
			return errors.New("Unauthorized")
		}

		var orderExists bool

		err = tx.Table("post_detail").
			Where("post_id = ? AND order_c = ?", postID, order).
			Select("count(*) > 0").
			Scan(&orderExists).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		postDetail := model.PostDetail{
			Image: filename,
			Order: int32(order),
		}

		if !orderExists {
			err := tx.Create(postDetail).Error
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return err
			}
		} else {
			err := tx.Table("post_detail").
				Where("post_id = ? AND order_c = ?", postID, order).
				Update("image", filename).Error
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return err
			}
		}

		// Create and write the file
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

		w.WriteHeader(http.StatusOK)

		return nil
	})

}

func HandleImage(w http.ResponseWriter, r *http.Request, params map[string]string) {
	filePath := UploadDir + "/" + params["image"]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Image not found.", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}
