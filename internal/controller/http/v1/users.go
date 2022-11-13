package v1

import (
	"forum/internal/entity"
	"net/http"
	"net/mail"
	"text/template"
	"time"
)

type signInErr struct {
	Code    int
	Message string
}

var (
	userNotExist     = "Такого пользователя не существует"
	userPassWrong    = "Неверный пароль, попробуйте ещё раз"
	passwordsNotSame = "Пароли не совпадают"
	emailFormatWrong = "Неправильный формат почты"
	userEmailExist   = "Пользователь с такой почтой уже существует"
	userNameExist    = "Пользователь с таким именем уже существует"
)

func (h *Handler) SignInPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signin_page/" {
		http.Error(w, "404: Page is Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "405: Method is not Allowed", http.StatusMethodNotAllowed)
		return
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
	if r.Method != http.MethodGet {
		http.Error(w, "405: Method is not Allowed", http.StatusMethodNotAllowed)
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

	s := signInErr{}
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
			Name:    "token",
			Value:   userWithSession.SessionToken,
			Expires: userWithSession.SessionTTL,
			Path:    "/",
			Domain:  "localhost",
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

	s := signInErr{}
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

func checkEmail(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}
