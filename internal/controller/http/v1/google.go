package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"forum/internal/entity"
)

var oauthStateString = "pseudo-random"

type GoogleContent struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func SignInGoogleHandler(w http.ResponseWriter, r *http.Request) {
	url := AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func AuthCodeURL(state string) string {
	var buf bytes.Buffer
	buf.WriteString("https://accounts.google.com/o/oauth2/auth")
	v := url.Values{"response_type": {"code"}, "client_id": {"952515746557-39jc67t7bb5olv9shjr6sp9v5he3mnt7.apps.googleusercontent.com"}}
	v.Set("redirect_uri", "http://localhost:8087/google_callback")
	v.Set("scope", "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile")
	v.Set("state", state)
	buf.WriteByte('?')
	buf.WriteString(v.Encode())
	return buf.String()
}

func (h *Handler) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	c := GoogleContent{}
	json.Unmarshal(content, &c)
	user := entity.User{
		Name:  c.Name,
		Email: c.Email,
	}

	err = h.Usecases.Users.SignUp(user)
	if err != nil {
		fmt.Println(err)
	}
	// if err != nil {
	// 	if err == entity.ErrUserEmailAlreadyExists {
	// 		content.ErrorMsg.Message = UserEmailAlreadyExist
	// 		valid = false
	// 	} else if err == entity.ErrUserNameAlreadyExists {
	// 		content.ErrorMsg.Message = UserNameAlreadyExist
	// 		valid = false
	// 	} else {
	// 		h.l.WriteLog(fmt.Errorf("v1 - SignUpHandler - SignUp: %w", err))
	// 		h.Errors(w, http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	fmt.Println(c)
}

func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := exchange(code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}

func exchange(code string) (Token, error) {
	var buf bytes.Buffer
	buf.WriteString("https://oauth2.googleapis.com/token?")
	v := url.Values{"grant_type": {"authorization_code"}, "code": {code}}
	v.Set("redirect_uri", "http://localhost:8087/google_callback")
	v.Set("client_id", "952515746557-39jc67t7bb5olv9shjr6sp9v5he3mnt7.apps.googleusercontent.com")
	v.Set("client_secret", "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")
	buf.WriteString(v.Encode())
	url := buf.String()
	req, _ := http.NewRequest("POST", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Token{}, err
	}
	defer resp.Body.Close()
	var token Token
	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &token)
	return token, nil
}

type Token struct {
	AccessToken string `json:"access_token"`
}
