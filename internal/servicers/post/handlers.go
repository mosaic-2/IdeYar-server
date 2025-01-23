package postImpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mosaic-2/IdeYar-server/internal/dbutil"
	"github.com/mosaic-2/IdeYar-server/internal/model"
	"github.com/mosaic-2/IdeYar-server/internal/servicers/util"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const UploadDir = "images"

const (
	postTitleKey        = "title"
	postImageKey        = "image"
	postDescriptionKey  = "description"
	postMinimumFundKey  = "minimumFund"
	postDeadlineKey     = "deadline"
	postCategoryKey     = "category"
	postDetailOrderKey  = "order"
	postDetailPostIDKey = "postID"
)

func HandlePostCreate(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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

	image, header, err := r.FormFile(postImageKey)
	if err != nil {
		http.Error(w, "Could not get uploaded file.", http.StatusBadRequest)
		return
	}
	defer image.Close()

	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		http.Error(w, "The uploaded file must be an image.", http.StatusBadRequest)
		return
	}

	tx := dbutil.GormDB(r.Context())

	err = tx.Transaction(func(tx *gorm.DB) error {
		post, err := getPostFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil
		}
		post.UserID = userID

		err = validatePost(post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil
		}

		err = tx.Create(&post).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return nil
		}

		// Create and write the file
		dst, err := os.Create(fmt.Sprintf("%s/%s", UploadDir, post.Image))
		if err != nil {
			http.Error(w, "Could not create a file.", http.StatusInternalServerError)
			return nil
		}
		defer dst.Close()

		if _, err := io.Copy(dst, image); err != nil {
			http.Error(w, "Failed to save the uploaded file.", http.StatusInternalServerError)
			return nil
		}

		PostID := struct {
			ID int64 `json:"id"`
		}{ID: post.ID}

		resp, err := json.Marshal(PostID)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return nil
		}

		_, err = w.Write(resp)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return nil
		}

		return nil
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func HandlePostDetailsCreate(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	hasImage := true
	image, header, err := r.FormFile(postImageKey)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			hasImage = false
		} else {
			http.Error(w, "Could not get uploaded file.", http.StatusBadRequest)
			return
		}
	} else {
		defer image.Close()
		if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
			http.Error(w, "The uploaded file must be an image.", http.StatusBadRequest)
			return
		}
	}

	tx := dbutil.GormDB(r.Context())

	err = tx.Transaction(func(tx *gorm.DB) error {
		postDetail, err := getPostDetailFromRequest(r, hasImage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil
		}

		hasAccess, err := hasCreateAccessPostDetail(tx, postDetail, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return nil
		}

		if !hasAccess {
			w.WriteHeader(http.StatusUnauthorized)
			return nil
		}

		err = validatePostDetail(postDetail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil
		}

		err = tx.Create(&postDetail).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return nil
		}

		// Create and write the file
		if hasImage {
			dst, err := os.Create(fmt.Sprintf("%s/%s", UploadDir, postDetail.Image))
			if err != nil {
				http.Error(w, "Could not create a file.", http.StatusInternalServerError)
				return err
			}
			defer dst.Close()

			if _, err := io.Copy(dst, image); err != nil {
				http.Error(w, "Failed to save the uploaded file.", http.StatusInternalServerError)
				return err
			}
		}

		w.WriteHeader(http.StatusOK)

		return nil
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getPostFromRequest(r *http.Request) (model.Post, error) {
	title := r.PostFormValue(postTitleKey)
	description := r.PostFormValue(postDescriptionKey)
	minimumFundStr := r.PostFormValue(postMinimumFundKey)
	deadlineStr := r.PostFormValue(postDeadlineKey)
	category := r.PostFormValue(postCategoryKey)
	imageFilename := util.GenerateFileName()

	minimumFund, err := decimal.NewFromString(minimumFundStr)
	if err != nil {
		return model.Post{}, errors.New("invalid minimum fund format")
	}

	deadlineDate, err := time.ParseInLocation(time.DateOnly, deadlineStr, time.Local)
	if err != nil {
		return model.Post{}, errors.New("invalid date format")
	}

	return model.Post{
		Title:        title,
		Description:  description,
		MinimumFund:  minimumFund,
		DeadlineDate: deadlineDate,
		Image:        imageFilename,
		Category:     category,
	}, nil
}

func getPostDetailFromRequest(r *http.Request, hasImage bool) (model.PostDetail, error) {
	title := r.PostFormValue(postTitleKey)
	description := r.PostFormValue(postDescriptionKey)
	orderStr := r.PostFormValue(postDetailOrderKey)
	postIDStr := r.PostFormValue(postDetailPostIDKey)

	var imageFilename string
	if hasImage {
		imageFilename = util.GenerateFileName()
	}

	order, err := strconv.ParseInt(orderStr, 10, 32)
	if err != nil {
		return model.PostDetail{}, errors.New("invalid order format")
	}

	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		return model.PostDetail{}, errors.New("invalid post id format")
	}

	return model.PostDetail{
		Title:       title,
		Description: description,
		Image:       imageFilename,
		Order:       int32(order),
		PostID:      postID,
	}, nil
}

func HandleImage(w http.ResponseWriter, r *http.Request, params map[string]string) {
	filePath := UploadDir + "/" + params["image"]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Image not found.", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}
