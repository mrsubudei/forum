package v1

import (
	"fmt"
	"forum/internal/entity"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func (h *Handler) PostPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		log.Println(fmt.Errorf("v1 - PostPageHandler - Atoi: %w", err))
	}
	if r.URL.Path != "/posts/"+path[len(path)-1] || err != nil || id <= 0 {
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	html, err := template.ParseFiles("templates/post.html")
	if err != nil {
		log.Println(fmt.Errorf("v1 - PostPageHandler - ParseFiles: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	post, err := h.usecases.Posts.GetById(int64(id))
	if err != nil {
		log.Println(fmt.Errorf("v1 - PostPageHandler - GetById: %w", err))
		errors.Code = http.StatusBadRequest
		errors.Message = ErrBadRequest
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		log.Printf("v1 - PostPageHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content)
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	post.ContentWeb = strings.Split(post.Content, "\\n")
	content.Post = post

	err = html.Execute(w, content)
	if err != nil {
		log.Println(fmt.Errorf("v1 - PostPageHandler - Execute: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) CreatePostPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		log.Printf("v1 - CreatePostPageHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content)
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	categories, err := h.usecases.Posts.GetAllCategories()
	if err != nil {
		log.Println(fmt.Errorf("v1 - CreatePostPageHandler - GetAllCategories: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
	content.Post.Categories = categories
	html, err := template.ParseFiles("templates/create_post.html")
	if err != nil {
		log.Println(fmt.Errorf("v1 - CreatePostPageHandler - ParseFiles: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	err = html.Execute(w, content)
	if err != nil {
		log.Println(fmt.Errorf("v1 - CreatePostPageHandler - Execute: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	r.ParseForm()
	if len(r.Form["title"]) == 0 || len(r.Form["content"]) == 0 {
		errors.Code = http.StatusBadRequest
		errors.Message = ErrBadRequest
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		log.Printf("v1 - CreatePostHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content)
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	valid := true
	postTitle := r.Form["title"][0]
	postContent := r.Form["content"][0]
	categories := r.Form["categories"]

	if len(categories) == 0 {
		content.ErrorMsg.Message = PostCategoryRequired
		valid = false
	}

	newPost := entity.Post{}
	newPost.Title = postTitle
	newPost.Content = strings.ReplaceAll(postContent, "\r\n", "\\n")
	newPost.Categories = categories
	newPost.User = content.User

	if !valid {
		html, err := template.ParseFiles("templates/create_post.html")
		if err != nil {
			log.Println(fmt.Errorf("v1 - CreatePostHandler - ParseFiles: %w", err))
			errors.Code = http.StatusInternalServerError
			errors.Message = ErrInternalServer
			h.Errors(w, errors)
			return
		}
		categories, err := h.usecases.Posts.GetAllCategories()
		if err != nil {
			log.Println(fmt.Errorf("v1 - CreatePostHandler - GetAllCategories: %w", err))
			errors.Code = http.StatusInternalServerError
			errors.Message = ErrInternalServer
			h.Errors(w, errors)
			return
		}
		content.Post.Categories = categories
		err = html.Execute(w, content)
		if err != nil {
			log.Println(fmt.Errorf("v1 - CreatePostHandler - Execute: %w", err))
			errors.Code = http.StatusInternalServerError
			errors.Message = ErrInternalServer
			h.Errors(w, errors)
			return
		}
	} else {
		err := h.usecases.Posts.CreatePost(newPost)
		if err != nil {
			log.Println(fmt.Errorf("v1 - CreatePostHandler - CreatePost: %w", err))
			errors.Code = http.StatusInternalServerError
			errors.Message = ErrInternalServer
			h.Errors(w, errors)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (h *Handler) PostPutLikeHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		log.Println(fmt.Errorf("v1 - PostPutLikeHandler - Atoi: %w", err))
	}
	if r.URL.Path != "/put_post_like/"+path[len(path)-1] || err != nil || id <= 0 {
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		log.Printf("v1 - PostPutLikeHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content)
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	post := entity.Post{
		Id: int64(id),
	}
	post.User.Id = content.User.Id

	err = h.usecases.Posts.MakeReaction(post, CommandPutLike)
	if err != nil {
		log.Println(fmt.Errorf("v1 - PostPutLikeHandler - MakeReaction: %w", err))
		errors.Code = http.StatusBadRequest
		errors.Message = ErrBadRequest
		h.Errors(w, errors)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func (h *Handler) PostPutDislikeHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		log.Println(fmt.Errorf("v1 - PostPutDislikeHandler - Atoi: %w", err))
	}
	if r.URL.Path != "/put_post_dislike/"+path[len(path)-1] || err != nil || id <= 0 {
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		log.Printf("v1 - PostPutDislikeHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content)
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	post := entity.Post{
		Id: int64(id),
	}
	post.User.Id = content.User.Id

	err = h.usecases.Posts.MakeReaction(post, CommandPutDislike)
	if err != nil {
		log.Println(fmt.Errorf("v1 - PostPutDislikeHandler - MakeReaction: %w", err))
		errors.Code = http.StatusBadRequest
		errors.Message = ErrBadRequest
		h.Errors(w, errors)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func (h *Handler) FindPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	query := path[len(path)-2]
	userId, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		log.Println(fmt.Errorf("v1 - FindPostsHandler - Atoi: %w", err))
	}

	if r.URL.Path != "/find_posts/"+query+"/"+path[len(path)-1] ||
		err != nil || userId <= 0 {
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		log.Printf("v1 - FindPostsHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content)
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	user := entity.User{Id: int64(userId)}

	posts, err := h.usecases.Posts.GetPostsByQuery(user, query)
	if err != nil {
		log.Println(fmt.Errorf("v1 - FindPostsHandler - GetPostsByQuery: %w", err))
		errors.Code = http.StatusBadRequest
		errors.Message = ErrBadRequest
		h.Errors(w, errors)
		return
	}
	content.Posts = posts

	html, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println(fmt.Errorf("v1 - FindPostsHandler - ParseFiles: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	err = html.Execute(w, content)
	if err != nil {
		log.Println(fmt.Errorf("v1 - FindPostsHandler - Execute: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
}
