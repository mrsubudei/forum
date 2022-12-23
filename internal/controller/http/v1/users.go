package v1

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"forum/internal/entity"
)

func (h *Handler) UserPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - UserPageHandler - Atoi: %w", err))
	}

	if r.URL.Path != "/users/"+path[len(path)-1] || err != nil || id <= 0 {
		h.Errors(w, http.StatusNotFound)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - UserPageHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}
	content.OwnerId = content.User.Id

	user, err := h.Usecases.Users.GetById(int64(id))
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - UserPageHandler - GetById: %w", err))
		if errors.Is(err, entity.ErrUserNotFound) {
			h.Errors(w, http.StatusNotFound)
			return
		}
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	if content.User.Id == int64(id) && content.Authorized || content.Admin {
		user.Owner = true
	}
	content.User = user

	err = h.ParseAndExecute(w, content, "templates/user.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - UserPageHandler - ParseAndExecute - %w", err))
	}
}

func (h *Handler) AllUsersPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/all_users_page/" {
		h.Errors(w, http.StatusNotFound)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - AllUsersPageHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	users, err := h.Usecases.Users.GetAllUsers()
	content.Users = users
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - AllUsersPageHandler - GetAllUsers: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	err = h.ParseAndExecute(w, content, "templates/all_users.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - AllUsersPageHandler - ParseAndExecute - %w", err))
	}
}

func (h *Handler) SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - AllUsersPageHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	err := h.ParseAndExecute(w, content, "templates/registration.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - SignUpPageHandler - ParseAndExecute - %w", err))
	}
}

func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		h.Errors(w, http.StatusInternalServerError)
	}
	var dateOfBirth string
	var city string
	var gender string
	var err error

	if len(r.Form["user"]) == 0 || len(r.Form["password"]) == 0 ||
		len(r.Form["email"]) == 0 || len(r.Form["confirm_password"]) == 0 {
		h.Errors(w, http.StatusBadRequest)
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

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - EditProfileHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	valid := true

	if !checkEmail(email) {
		content.ErrorMsg.Message = EmailFormatWrong
		valid = false
	}
	if password != confirmPassword {
		content.ErrorMsg.Message = PasswordsNotSame
		valid = false
	}

	user := entity.User{
		Name:        name,
		Password:    password,
		Email:       strings.ToLower(email),
		City:        city,
		Gender:      gender,
		DateOfBirth: dateOfBirth,
	}

	content.User = user

	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		err := h.ParseAndExecute(w, content, "templates/registration.html")
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - SignUpHandler - ParseAndExecute #1 - %w", err))
		}
		return
	}

	err = h.Usecases.Users.SignUp(user)
	if err != nil {
		if err == entity.ErrUserEmailAlreadyExists {
			content.ErrorMsg.Message = UserEmailAlreadyExist
			valid = false
		} else if err == entity.ErrUserNameAlreadyExists {
			content.ErrorMsg.Message = UserNameAlreadyExist
			valid = false
		} else {
			h.l.WriteLog(fmt.Errorf("v1 - SignUpHandler - SignUp: %w", err))
			h.Errors(w, http.StatusInternalServerError)
			return
		}
	}

	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		err := h.ParseAndExecute(w, content, "templates/registration.html")
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - SignUpHandler - ParseAndExecute #2 - %w", err))
		}
	} else {
		http.Redirect(w, r, "/signin_page", http.StatusFound)
	}
}

func (h *Handler) SignInPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	err := h.ParseAndExecute(w, Content{}, "templates/login.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - SignInPageHandler - ParseAndExecute - %w", err))
	}
}

func (h *Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.Errors(w, http.StatusInternalServerError)
	}

	if len(r.Form["user"]) == 0 || len(r.Form["password"]) == 0 {
		h.Errors(w, http.StatusBadRequest)
		return
	}

	data := r.Form["user"][0]
	password := r.Form["password"][0]

	user := entity.User{
		Password: password,
	}
	if checkEmail(data) {
		user.Email = strings.ToLower(data)
	} else {
		user.Name = data
	}

	valid := true
	content := Content{}

	err := h.Usecases.Users.SignIn(user)

	if err != nil && !strings.Contains(err.Error(), ErrNoRowsInResult) {
		h.l.WriteLog(fmt.Errorf("v1 - SignInHandler - SignIn: %w", err))
	}
	if err == entity.ErrUserNotFound {
		content.ErrorMsg.Message = UserNotExist
		valid = false
	} else if err == entity.ErrUserPasswordIncorrect {
		content.ErrorMsg.Message = UserPassWrong
		valid = false
	}

	if !valid {
		w.WriteHeader(http.StatusUnauthorized)

		err := h.ParseAndExecute(w, content, "templates/login.html")
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - SignInHandler - ParseAndExecute - %w", err))
		}
		return

	} else {
		id, err := h.Usecases.Users.GetIdBy(user)
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - SignInHandler - GetIdBy: %w", err))
			h.Errors(w, http.StatusBadRequest)
			return
		}

		userWithSession, err := h.Usecases.Users.GetSession(id)
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - SignInHandler - GetSession: %w", err))
			h.Errors(w, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   userWithSession.SessionToken,
			Expires: userWithSession.SessionTTL,
			Path:    "/",
			Domain:  h.Cfg.Server.Host,
		})
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (h *Handler) EditProfilePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - EditProfilePageHandler - Atoi: %w", err))
	}
	if r.URL.Path != "/edit_profile_page/"+path[len(path)-1] || err != nil || id <= 0 {
		h.Errors(w, http.StatusNotFound)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - EditProfilePageHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	if content.User.Id != 1 && content.User.Id != int64(id) {
		h.Errors(w, http.StatusForbidden)
		return
	}

	err = h.ParseAndExecute(w, content, "templates/edit_profile.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - EditProfilePageHandler - ParseAndExecute - %w", err))
	}
}

func (h *Handler) EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - EditProfileHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.Errors(w, http.StatusInternalServerError)
	}

	if len(r.Form["id"]) != 0 && r.Form["id"][0] != "" {
		id, err := strconv.Atoi(r.Form["id"][0])
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - EditProfileHandler - Atoi: %w", err))
			h.Errors(w, http.StatusInternalServerError)
			return
		}
		content.User.Id = int64(id)
	}

	existUser, err := h.Usecases.Users.GetById(content.User.Id)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - EditProfileHandler - GetById: %w", err))
		h.Errors(w, http.StatusNotFound)
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

	err = h.Usecases.Users.UpdateUserInfo(existUser, UpdateQueryInfo)

	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - EditProfileHandler - UpdateUserInfo: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users/"+strconv.Itoa(int(content.User.Id)), http.StatusFound)
}

func (h *Handler) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - SignOutHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	err := h.Usecases.Users.DeleteSession(content.User)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - SignOutHandler - DeleteSession: %w", err))
		h.Errors(w, http.StatusInternalServerError)
	}
	time.Sleep(time.Second)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) FindReactedUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	query := path[len(path)-3]
	reaction := path[len(path)-2]
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - FindReactedUsersHandler - Atoi: %w", err))
	}

	if (query != QueryPost && query != QueryComment) ||
		(reaction != QueryLiked && reaction != QueryDisliked) {
		h.Errors(w, http.StatusNotFound)
		return
	}

	if r.URL.Path != "/find_reacted_users/"+query+"/"+reaction+"/"+path[len(path)-1] ||
		err != nil || id <= 0 {
		h.Errors(w, http.StatusNotFound)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - FindReactedUsersHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	switch query {
	case QueryPost:
		if reaction == QueryLiked {
			content.Users, err = h.Usecases.Posts.GetReactions(int64(id), QueryLiked)
			if err != nil {
				h.l.WriteLog(fmt.Errorf("v1 - FindReactedUsersHandler - GetReactions #1: %w", err))
				h.Errors(w, http.StatusNotFound)
				return
			}
			content.Message = ReactionMessageLike
		} else if reaction == QueryDisliked {
			content.Users, err = h.Usecases.Posts.GetReactions(int64(id), QueryDisliked)
			if err != nil {
				h.l.WriteLog(fmt.Errorf("v1 - FindReactedUsersHandler - GetReactions #2: %w", err))
				h.Errors(w, http.StatusNotFound)
				return
			}
			content.Message = ReactionMessageDislike
		}
	case QueryComment:
		if reaction == QueryLiked {
			content.Users, err = h.Usecases.Comments.GetReactions(int64(id), QueryLiked)
			if err != nil {
				h.l.WriteLog(fmt.Errorf("v1 - FindReactedUsersHandler - GetReactions #3: %w", err))
				h.Errors(w, http.StatusNotFound)
				return
			}
			content.Message = ReactionMessageLike
		} else if reaction == QueryDisliked {
			content.Users, err = h.Usecases.Comments.GetReactions(int64(id), QueryDisliked)
			if err != nil {
				h.l.WriteLog(fmt.Errorf("v1 - FindReactedUsersHandler - GetReactions #4: %w", err))
				h.Errors(w, http.StatusBadRequest)
				return
			}
			content.Message = ReactionMessageDislike
		}
	}

	err = h.ParseAndExecute(w, content, "templates/reacted_users.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - FindReactedUsersHandler - ParseAndExecute - %w", err))
	}
}

func (h *Handler) GetExistedSession(w http.ResponseWriter, r *http.Request) entity.User {
	foundUser := entity.User{}
	cookie, err := r.Cookie(h.Cfg.TokenManager.TokenName)
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
	id, err := h.Usecases.Users.GetIdBy(user)
	if err != nil {
		if !strings.Contains(err.Error(), ErrNoRowsInResult) {
			h.l.WriteLog(fmt.Errorf("v1 - GetExistedSession - GetIdBy: %w", err))
		}
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
