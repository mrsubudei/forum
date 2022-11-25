package v1

import (
	"forum/internal/entity"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func (h *Handler) UserPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if r.URL.Path != "/users/"+path[len(path)-1] || err != nil || id <= 0 {
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	user, err := h.usecases.Users.GetById(int64(id))
	if err != nil {
		if strings.Contains(err.Error(), entity.ErrUserNotFound.Error()) {
			errors.Code = http.StatusBadRequest
			errors.Message = UserNotExist
			h.Errors(w, errors)
			return
		}
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	if content.User.Id == int64(id) && content.Authorized || content.Admin {
		user.Owner = true
	}
	content.User = user

	html, err := template.ParseFiles("templates/user.html")
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	err = html.Execute(w, content)
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) AllUsersPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	if r.URL.Path != "/all_users_page/" {
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	users, err := h.usecases.Users.GetAllUsers()
	content.Users = users
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	html, err := template.ParseFiles("templates/all_users.html")
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
	err = html.Execute(w, content)
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) SignInPageHandler(w http.ResponseWriter, r *http.Request) {

	html, err := template.ParseFiles("templates/login.html")
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	err = html.Execute(w, nil)
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	html, err := template.ParseFiles("templates/registration.html")
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	err = html.Execute(w, nil)
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) EditProfilePageHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if r.URL.Path != "/edit_profile_page/"+path[len(path)-1] || err != nil || id <= 0 {
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	if content.User.Id != 1 && content.User.Id != int64(id) {
		errors.Code = http.StatusForbidden
		errors.Message = ErrLowAccessLevel
		h.Errors(w, errors)
		return
	}

	html, err := template.ParseFiles("templates/edit_profile.html")
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	err = html.Execute(w, content)
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	r.ParseForm()

	if len(r.Form["id"]) != 0 && r.Form["id"][0] != "" {
		id, err := strconv.Atoi(r.Form["id"][0])
		if err != nil {
			errors.Code = http.StatusInternalServerError
			errors.Message = ErrInternalServer
			h.Errors(w, errors)
			return
		}
		content.User.Id = int64(id)
	}

	existUser, err := h.usecases.Users.GetById(content.User.Id)
	if err != nil {
		errors.Code = http.StatusBadRequest
		errors.Message = ErrBadRequest
		h.Errors(w, errors)
		return
	}

	if len(r.Form["date_of_birth"]) != 0 && r.Form["date_of_birth"][0] != "" {
		existUser.DateOfBirth = r.Form["date_of_birth"][0]
	}
	if len(r.Form["city"]) != 0 && r.Form["city"][0] != "" {
		existUser.City = r.Form["city"][0]
	}
	if len(r.Form["gender"]) != 0 && r.Form["gender"][0] != "" {
		existUser.Gender = r.Form["gender"][0]
	}
	if len(r.Form["sign"]) != 0 && r.Form["sign"][0] != "" {
		existUser.Sign = r.Form["sign"][0]
	}
	if len(r.Form["role"]) != 0 && r.Form["role"][0] != "" {
		existUser.Role = r.Form["role"][0]
	}

	err = h.usecases.Users.UpdateUserInfo(existUser, UpdateQueryInfo)

	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	http.Redirect(w, r, "/users/"+strconv.Itoa(int(content.User.Id)), http.StatusFound)
}

func (h *Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}
	r.ParseForm()

	if len(r.Form["user"]) == 0 || len(r.Form["password"]) == 0 {
		errors.Code = http.StatusBadRequest
		errors.Message = ErrBadRequest
		h.Errors(w, errors)
		return
	}

	data := r.Form["user"][0]
	password := r.Form["password"][0]

	user := entity.User{
		Password: password,
	}
	if checkEmail(data) {
		user.Email = data
	} else {
		user.Name = data
	}

	valid := true
	content := Content{}

	err := h.usecases.Users.SignIn(user)

	if err == entity.ErrUserNotFound {
		content.ErrorMsg.Message = UserNotExist
		valid = false
	} else if err == entity.ErrUserPasswordIncorrect {
		content.ErrorMsg.Message = UserPassWrong
		valid = false
	}

	if !valid {
		html, err := template.ParseFiles("templates/login.html")
		if err != nil {
			errors.Code = http.StatusInternalServerError
			errors.Message = ErrInternalServer
			h.Errors(w, errors)
			return
		}
		err = html.Execute(w, content)
		if err != nil {
			errors.Code = http.StatusInternalServerError
			errors.Message = ErrInternalServer
			h.Errors(w, errors)
			return
		}
	} else {
		id, err := h.usecases.Users.GetIdBy(user)
		if err != nil {
			errors.Code = http.StatusBadRequest
			errors.Message = ErrBadRequest
			h.Errors(w, errors)
			return
		}

		userWithSession, err := h.usecases.Users.GetSession(id)
		if err != nil {
			errors.Code = http.StatusInternalServerError
			errors.Message = ErrInternalServer
			h.Errors(w, errors)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   userWithSession.SessionToken,
			Expires: userWithSession.SessionTTL,
			Path:    "/",
			Domain:  SessionDomain,
		})
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}
	r.ParseForm()
	var dateOfBirth string
	var city string
	var gender string
	var err error

	if len(r.Form["user"]) == 0 || len(r.Form["password"]) == 0 ||
		len(r.Form["email"]) == 0 || len(r.Form["confirm_password"]) == 0 {
		errors.Code = http.StatusBadRequest
		errors.Message = ErrBadRequest
		h.Errors(w, errors)
		return
	}

	name := r.Form["user"][0]
	email := r.Form["email"][0]
	password := r.Form["password"][0]
	confirmPassword := r.Form["confirm_password"][0]

	if len(r.Form["date_of_birth"]) != 0 {
		dateOfBirth = r.Form["date_of_birth"][0]
	}
	if len(r.Form["city"]) != 0 {
		city = r.Form["city"][0]
	}
	if len(r.Form["gender"]) != 0 {
		gender = r.Form["gender"][0]
	}

	errorMessage := ErrMessage{}
	valid := true

	if !checkEmail(email) {
		errorMessage.Message = EmailFormatWrong
		valid = false
	}
	if password != confirmPassword {
		errorMessage.Message = PasswordsNotSame
		valid = false
	}

	user := entity.User{
		Name:        name,
		Password:    password,
		Email:       email,
		City:        city,
		Gender:      gender,
		DateOfBirth: dateOfBirth,
	}
	err = h.usecases.Users.SignUp(user)
	if err != nil {
		if err == entity.ErrUserEmailAlreadyExists {
			errorMessage.Message = UserEmailExist
			valid = false
		} else if err == entity.ErrUserNameAlreadyExists {
			errorMessage.Message = UserNameExist
			valid = false
		}
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	if !valid {
		html, err := template.ParseFiles("templates/registration.html")
		if err != nil {
			errors.Code = http.StatusInternalServerError
			errors.Message = ErrInternalServer
			h.Errors(w, errors)
			return
		}
		err = html.Execute(w, errorMessage)
		if err != nil {
			errors.Code = http.StatusInternalServerError
			errors.Message = ErrInternalServer
			h.Errors(w, errors)
			return
		}
	} else {
		http.Redirect(w, r, "/signin_page", http.StatusFound)
	}
}

func (h *Handler) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	err := h.usecases.Users.DeleteSession(content.User)
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
	}
	time.Sleep(time.Second)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) FindReactedUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	query := path[len(path)-3]
	reaction := path[len(path)-2]
	id, err := strconv.Atoi(path[len(path)-1])

	if (query != QueryPost && query != QueryComment) ||
		(reaction != QueryLiked && reaction != QueryDisliked) {
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	if r.URL.Path != "/find_reacted_users/"+query+"/"+reaction+"/"+path[len(path)-1] ||
		err != nil || id <= 0 {
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	switch query {
	case QueryPost:
		if reaction == QueryLiked {
			content.Users, err = h.usecases.Posts.GetReactions(int64(id), QueryLiked)
			if err != nil {
				errors.Code = http.StatusBadRequest
				errors.Message = ErrBadRequest
				h.Errors(w, errors)
				return
			}
			content.Message = ReactionMessageLike
		} else if reaction == QueryDisliked {
			content.Users, err = h.usecases.Posts.GetReactions(int64(id), QueryDisliked)
			if err != nil {
				errors.Code = http.StatusBadRequest
				errors.Message = ErrBadRequest
				h.Errors(w, errors)
				return
			}
			content.Message = ReactionMessageDislike
		}
	case QueryComment:
		if reaction == QueryLiked {
			content.Users, err = h.usecases.Comments.GetReactions(int64(id), QueryLiked)
			if err != nil {
				errors.Code = http.StatusBadRequest
				errors.Message = ErrBadRequest
				h.Errors(w, errors)
				return
			}
			content.Message = ReactionMessageLike
		} else if reaction == QueryDisliked {
			content.Users, err = h.usecases.Comments.GetReactions(int64(id), QueryDisliked)
			if err != nil {
				errors.Code = http.StatusBadRequest
				errors.Message = ErrBadRequest
				h.Errors(w, errors)
				return
			}
			content.Message = ReactionMessageDislike
		}
	}

	html, err := template.ParseFiles("templates/reacted_users.html")
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	err = html.Execute(w, content)
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) GetExistedSession(w http.ResponseWriter, r *http.Request) entity.User {
	foundUser := entity.User{}
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return foundUser
		}
		w.WriteHeader(http.StatusBadRequest)
		return foundUser
	}
	token := cookie.Value
	user := entity.User{
		SessionToken: token,
	}
	id, err := h.usecases.Users.GetIdBy(user)
	if err != nil {
		return foundUser
	}

	foundUser.Id = id
	foundUser.SessionToken = token

	return foundUser
}

func checkEmail(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}
