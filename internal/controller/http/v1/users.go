package v1

import (
	"forum/internal/entity"
	"net/http"
	"text/template"
	"time"
)

type signInErr struct {
	Code    int
	Message string
}

var (
	userNotExist     = "Такого пользователя не существует."
	userPassWrong    = "Неверный пароль, попробуйте ещё раз."
	passwordsNotSame = "Пароли не совпадают"
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
	name := r.Form["user"][0]
	password := r.Form["password"][0]

	user := entity.User{
		Name:     name,
		Password: password,
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
		s.Message = userPassWrong
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
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405: Method is not Allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	name := r.Form["user"][0]
	email := r.Form["email"][0]
	password := r.Form["password"][0]
	confirmPassword := r.Form["confirm_password"][0]
	dateOfBirth := r.Form["date_of_birth"][0]

	s := signInErr{}
	valid := true

	if password != confirmPassword {
		s.Message = passwordsNotSame
		valid = false
	}
	parsed, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	user := entity.User{
		Name:        name,
		Password:    password,
		Email:       email,
		DateOfBirth: parsed,
	}
	err = h.usecases.Users.SignUp(user)

	if err == entity.ErrUserEmailAlreadyExists {
		s.Message = userNotExist
		valid = false
	} else if err == entity.ErrUserPasswordIncorrect {
		s.Message = userPassWrong
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
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
