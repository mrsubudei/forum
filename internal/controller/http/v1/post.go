package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/entity"
)

func (h *Handler) PostPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - PostPageHandler - Atoi: %w", err))
	}
	if r.URL.Path != "/posts/"+path[len(path)-1] || err != nil || id <= 0 {
		h.Errors(w, http.StatusNotFound)
		return
	}

	post, err := h.Usecases.Posts.GetById(int64(id))
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - PostPageHandler - GetById: %w", err))
		h.Errors(w, http.StatusNotFound)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - PostPageHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	post.ContentWeb = strings.Split(post.Content, "\\n")
	content.Post = post

	err = h.ParseAndExecute(w, content, "templates/post.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - PostPageHandler - ParseAndExecute - %w", err))
	}
}

func (h *Handler) CreatePostPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - CreatePostPageHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	categories, err := h.Usecases.Posts.GetAllCategories()
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - CreatePostPageHandler - GetAllCategories: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return
	}
	content.Post.Categories = categories

	err = h.ParseAndExecute(w, content, "templates/create_post.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - CreatePostPageHandler - ParseAndExecute - %w", err))
	}
}

func (h *Handler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.Errors(w, http.StatusInternalServerError)
	}

	if len(r.Form["title"][0]) == 0 || len(r.Form["content"][0]) == 0 {
		h.Errors(w, http.StatusBadRequest)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - CreatePostHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
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
	content.Post = newPost

	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		categories, err := h.Usecases.Posts.GetAllCategories()
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - CreatePostHandler - GetAllCategories: %w", err))
			h.Errors(w, http.StatusInternalServerError)
			return
		}

		content.Post.Categories = categories
		err = h.ParseAndExecute(w, content, "templates/create_post.html")
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - CreatePostHandler - ParseAndExecute - %w", err))
		}

	} else {
		err := h.Usecases.Posts.CreatePost(newPost)
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - CreatePostHandler - CreatePost: %w", err))
			h.Errors(w, http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (h *Handler) PostPutLikeHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - PostPutLikeHandler - Atoi: %w", err))
	}
	if r.URL.Path != "/put_post_like/"+path[len(path)-1] || err != nil || id <= 0 {
		h.Errors(w, http.StatusNotFound)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - PostPutLikeHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	post := entity.Post{
		Id: int64(id),
	}
	post.User.Id = content.User.Id

	err = h.Usecases.Posts.MakeReaction(post, CommandPutLike)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - PostPutLikeHandler - MakeReaction: %w", err))
		h.Errors(w, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func (h *Handler) PostPutDislikeHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - PostPutDislikeHandler - Atoi: %w", err))
	}
	if r.URL.Path != "/put_post_dislike/"+path[len(path)-1] || err != nil || id <= 0 {
		h.Errors(w, http.StatusNotFound)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - PostPutDislikeHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	post := entity.Post{
		Id: int64(id),
	}
	post.User.Id = content.User.Id

	err = h.Usecases.Posts.MakeReaction(post, CommandPutDislike)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - PostPutDislikeHandler - MakeReaction: %w", err))
		h.Errors(w, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func (h *Handler) FindPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	query := path[len(path)-2]
	userId, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - FindPostsHandler - Atoi: %w", err))
	}

	if r.URL.Path != "/find_posts/"+query+"/"+path[len(path)-1] ||
		err != nil || userId <= 0 {
		h.Errors(w, http.StatusNotFound)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - FindPostsHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	user := entity.User{Id: int64(userId)}

	posts, err := h.Usecases.Posts.GetPostsByQuery(user, query)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - FindPostsHandler - GetPostsByQuery: %w", err))
		h.Errors(w, http.StatusNotFound)
		return
	}
	content.Posts = posts

	err = h.ParseAndExecute(w, content, "templates/index.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - FindPostsHandler - ParseAndExecute - %w", err))
	}
}
