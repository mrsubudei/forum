package v1

import (
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"strings"

	"forum/internal/entity"
	"forum/pkg/auth"
)

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		h.Errors(w, http.StatusNotFound)
		return
	}

	key := Key("content")
	content, ok := r.Context().Value(key).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - IndexHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}
	posts, err := h.Usecases.Posts.GetAllPosts()
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - IndexHandler - GetAllPosts: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return
	}
	content.Posts = posts
	err = h.ParseAndExecute(w, content, "templates/index.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - IndexHandler - ParseAndExecute - %w", err))
	}
}

func (h *Handler) SearchPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - SearchPageHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	err := h.ParseAndExecute(w, content, "templates/search.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - SearchPageHandler - ParseAndExecute - %w", err))
	}
}

func (h *Handler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.Errors(w, http.StatusInternalServerError)
	}

	if len(r.Form["search"]) == 0 {
		h.Errors(w, http.StatusBadRequest)
		return
	}

	searchRequest := strings.ToLower(r.Form["search"][0])
	posts, err := h.Usecases.Posts.GetAllPosts()
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - SearchHandler - GetAllPosts: %w", err))
		h.Errors(w, http.StatusBadRequest)
		return
	}

	filtered := h.filterPosts(posts, searchRequest)
	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - SearchHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	content.Posts = filtered

	err = h.ParseAndExecute(w, content, "templates/index.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - SearchHandler - ParseAndExecute - %w", err))
	}
}

func (h *Handler) CreateCategoryPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - CreateCategoryPageHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	if !content.Admin {
		h.Errors(w, http.StatusForbidden)
		return
	}

	err := h.ParseAndExecute(w, content, "templates/create_category.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - CreateCategoryPageHandler - ParseAndExecute - %w", err))
	}
}

func (h *Handler) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.Errors(w, http.StatusInternalServerError)
	}

	if len(r.Form["category"]) == 0 || len(r.Form["category"][0]) == 0 {
		h.Errors(w, http.StatusBadRequest)
		return
	}
	data := r.Form["category"][0]
	categories := strings.Split(data, "\r\n")

	err := h.Usecases.Posts.CreateCategories(categories)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - CreateCategoryHandler - CreateCategories: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) SearchByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	category := path[len(path)-1]
	if r.URL.Path != "/categories/"+category {
		h.Errors(w, http.StatusNotFound)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - SearchByCategoryHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	posts, err := h.Usecases.Posts.GetAllByCategory(category)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - SearchByCategoryHandler - GetAllByCategory: %w", err))
		if strings.Contains(err.Error(), entity.ErrPostNotFound.Error()) {
			h.Errors(w, http.StatusBadRequest)
			return
		}
		h.Errors(w, http.StatusInternalServerError)
		return
	}
	content.Posts = posts

	err = h.ParseAndExecute(w, content, "templates/index.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - SearchByCategoryHandler - ParseAndExecute - %w", err))
	}
}

func (h *Handler) filterPosts(posts []entity.Post, request string) []entity.Post {
	var filtered []entity.Post
	if len(posts) == 0 {
		return filtered
	}
	for _, post := range posts {
		found := false
		if strings.Contains(strings.ToLower(post.Title), request) {
			filtered = append(filtered, post)
			continue
		} else if strings.Contains(strings.ToLower(post.Content), request) {
			filtered = append(filtered, post)
			continue
		} else if strings.Contains(strings.ToLower(post.User.Name), request) {
			filtered = append(filtered, post)
			continue
		}
		for _, val := range post.Categories {
			if strings.Contains(strings.ToLower(val), request) {
				filtered = append(filtered, post)
				found = true
				break
			}
		}
		if !found {
			for _, val := range post.Comments {
				if strings.Contains(strings.ToLower(val.Content), request) {
					filtered = append(filtered, post)
					break
				}
			}
		}
	}
	return filtered
}

func (h *Handler) GetImage(w http.ResponseWriter, r *http.Request) (string, error) {
	file, header, err := r.FormFile("image")
	if err != nil {
		return "", nil
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")
	typeSl := strings.Split(mimeType, "/")
	imageType := typeSl[1]
	if imageType != "jpeg" && imageType != "png" &&
		imageType != "gif" && imageType != "octet-stream" {
		return "", errors.New(imageTypeForbidden)
	}

	root := getRootPath()

	if _, err := os.Stat(root + "templates/img/storage"); os.IsNotExist(err) {
		os.MkdirAll(root+"templates/img/storage", os.ModePerm)
	}

	manager := auth.NewManager(h.Cfg)
	generated, err := manager.NewToken()
	if err != nil {
		return "", fmt.Errorf("newToken: %w", err)
	}

	path := root + "templates/img/storage/" + generated + "." + imageType
	targetFile, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("create: %w", err)
	}
	defer targetFile.Close()

	written, err := io.Copy(targetFile, file)
	if err != nil {
		return "", fmt.Errorf("copy: %w", err)
	}
	if written >= (ImageSizeInt << 20) {
		err := os.Remove(path)
		if err != nil {
			return "", fmt.Errorf("remove: %w", err)
		}
		return "", errors.New(imageTooLarge)
	}

	return path, nil
}

func (h *Handler) CheckSizeExceeded(path string) (bool, error) {
	if path == "" {
		return false, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("open: %w", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		return false, fmt.Errorf("decodeConfig: %w", err)
	}
	if image.Width > ImageWidthInt || image.Height > ImageHeightInt {
		return true, nil
	}
	return false, nil
}
