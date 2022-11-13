package v1

import (
	"fmt"
	"forum/internal/entity"
	"net/http"
	"net/mail"
	"text/template"
	"time"
)

type ErrMessage struct {
	Code    int
	Message string
}

var (
	sessionDomain    = "localhost"
	userNotExist     = "Такого пользователя не существует"
	userPassWrong    = "Неверный пароль, попробуйте ещё раз"
	passwordsNotSame = "Пароли не совпадают"
	emailFormatWrong = "Неправильный формат почты"
	userEmailExist   = "Пользователь с такой почтой уже существует"
	userNameExist    = "Пользователь с таким именем уже существует"
	pageNotFound     = "Страница не найдена"
)

func (h *Handler) SignInPageHandler(w http.ResponseWriter, r *http.Request) {
	authorized := h.checkSession(w, r)
	if authorized {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	html, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = html.Execute(w, nil)
	if err != nil {
		http.Error(w, "404: Not Found", 404)
		return
	}
}

func (h *Handler) SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup_page/" {
		http.Error(w, "404: Page is Not Found", http.StatusNotFound)
		return
	}

	html, err := template.ParseFiles("templates/registration.html")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = html.Execute(w, nil)
	if err != nil {
		http.Error(w, "404: Not Found", 404)
		return
	}
}

func (h *Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405: Method is not Allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
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
	err := h.usecases.Users.SignIn(user)

	s := ErrMessage{}
	valid := true
	if err == entity.ErrUserNotFound {
		s.Message = userNotExist
		valid = false
	} else if err == entity.ErrUserPasswordIncorrect {
		s.Message = userPassWrong
		valid = false
	}

	if !valid {
		html, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = html.Execute(w, s)
		if err != nil {
			http.Error(w, "404: Not Found", 404)
			return
		}
	} else {
		id, err := h.usecases.Users.GetIdBy(user)
		if err != nil {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}

		userWithSession, err := h.usecases.Users.GetSession(id)
		if err != nil {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   userWithSession.SessionToken,
			Expires: userWithSession.SessionTTL,
			Path:    "/",
			Domain:  sessionDomain,
		})
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405: Method is not Allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	var dateOfBirth string
	var city string
	var gender string
	var parsed time.Time
	var err error

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

	s := ErrMessage{}
	valid := true

	if !checkEmail(email) {
		s.Message = emailFormatWrong
		valid = false
	}
	if password != confirmPassword {
		s.Message = passwordsNotSame
		valid = false
	}
	if dateOfBirth != "" {
		parsed, err = time.Parse("2006-01-02", dateOfBirth)
		if err != nil {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	user := entity.User{
		Name:        name,
		Password:    password,
		Email:       email,
		City:        city,
		Gender:      gender,
		DateOfBirth: parsed,
	}
	err = h.usecases.Users.SignUp(user)

	if err == entity.ErrUserEmailAlreadyExists {
		s.Message = userEmailExist
		valid = false
	} else if err == entity.ErrUserNameAlreadyExists {
		s.Message = userNameExist
		valid = false
	}
	if !valid {
		html, err := template.ParseFiles("templates/registration.html")
		if err != nil {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = html.Execute(w, s)
		if err != nil {
			http.Error(w, "404: Not Found", 404)
			return
		}
	} else {
		http.Redirect(w, r, "/signin_page", http.StatusFound)
	}
}

func (h *Handler) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	foundUser := h.getExistedSession(w, r)

	err := h.usecases.Users.DeleteSession(foundUser)
	if err != nil {
		fmt.Println("errorchik")
	}
	time.Sleep(time.Second)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) checkSession(w http.ResponseWriter, r *http.Request) bool {
	foundUser := h.getExistedSession(w, r)
	authorized, err := h.usecases.Users.CheckSession(foundUser)
	if err != nil {
		return false
	}
	return authorized
}

func (h *Handler) getExistedSession(w http.ResponseWriter, r *http.Request) entity.User {
	foundUser := entity.User{}
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return foundUser
		}
		w.WriteHeader(http.StatusBadRequest)
		return foundUser
	}
	tokenExisted := cookie.Value
	user := entity.User{
		SessionToken: tokenExisted,
	}
	id, err := h.usecases.Users.GetIdBy(user)
	if err != nil {
		return foundUser
	}
	foundUser = entity.User{
		Id:           id,
		SessionToken: tokenExisted,
	}
	return foundUser
}

func checkEmail(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}
