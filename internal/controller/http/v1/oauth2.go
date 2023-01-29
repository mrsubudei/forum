package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"forum/internal/entity"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var oauthStateString = "pseudo-random"

type OauthContent struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type OauthToken struct {
	AccessToken string `json:"access_token,omitempty"`
	Tokens      struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"google_tokens,omitempty"`
}

func (h *Handler) OauthSignHandler(w http.ResponseWriter, r *http.Request) {
	configFile, err := os.Open("tokens.json")
	if err != nil {
		log.Fatal(err)
	}
	data, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()
	tokens := OauthToken{}
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		log.Fatal(err)
	}
	h.Cfg.GoogleTokens.ClientID = tokens.Tokens.ClientID
	h.Cfg.GoogleTokens.ClientSecret = tokens.Tokens.ClientSecret

	err = h.apiIdentify(r)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - OauthSignHandler - apiIdentify: %w", err))
		h.Errors(w, http.StatusNotFound)
		return
	}

	var buf bytes.Buffer
	buf.WriteString(h.Cfg.Oauth.AuthURL)
	v := url.Values{"response_type": {"code"}, "client_id": {h.Cfg.GoogleTokens.ClientID}}
	v.Set("redirect_uri", "http://localhost:8087/oauth2_callback")
	v.Set("scope", "https://www.googleapis.com/auth/userinfo.email")
	v.Set("state", oauthStateString)
	buf.WriteByte('?')
	buf.WriteString(v.Encode())
	url := buf.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) apiIdentify(r *http.Request) error {
	path := strings.Split(r.URL.Path, "/")
	endpoint := path[len(path)-1]
	if r.URL.Path != "/oauth2_signin/"+endpoint {
		return fmt.Errorf("wrong endpoint")
	}

	switch endpoint {
	case "google":
		h.Cfg.Oauth.AuthURL = GoogleAuthURL
		h.Cfg.Oauth.TokenURL = GoogleTokenURL
	default:
		return fmt.Errorf("wrong endpoint")
	}

	return nil
}

func (h *Handler) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {

	token, err := h.exchageCode(r)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - GoogleCallbackHandler - exchageCode: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	content, err := h.tokenToCall(token)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	userData := OauthContent{}
	json.Unmarshal(content, &userData)
	user := entity.User{
		Email: userData.Email,
	}

	id, err := h.Usecases.Users.GetIdBy(user)
	if id == 0 {
		user.Name = userData.Name
		h.Usecases.Users.SignUp(user)
		id, err = h.Usecases.Users.GetIdBy(user)
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - GoogleCallbackHandler - GetIdBy#1: %w", err))
			h.Errors(w, http.StatusInternalServerError)
			return
		}
	}
	if err != nil {
		if !errors.Is(err, entity.ErrUserNotFound) {
			h.l.WriteLog(fmt.Errorf("v1 - GoogleCallbackHandler - GetIdBy#2: %w", err))
			h.Errors(w, http.StatusInternalServerError)
			return
		}
	}

	err = h.Usecases.Users.SignIn(user)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - GoogleCallbackHandler - SignIn: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	userWithSession, err := h.Usecases.Users.GetSession(id)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - GoogleCallbackHandler - GetSession: %w", err))
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

func (h *Handler) exchageCode(r *http.Request) (OauthToken, error) {
	token := OauthToken{}
	state := r.FormValue("state")
	code := r.FormValue("code")
	if state != oauthStateString {
		return token, fmt.Errorf("invalid oauth state")
	}

	var buf bytes.Buffer
	buf.WriteString(h.Cfg.Oauth.TokenURL)
	v := url.Values{"grant_type": {"authorization_code"}, "code": {code}}
	v.Set("redirect_uri", "http://localhost:8087/oauth2_callback")
	v.Set("client_id", h.Cfg.GoogleTokens.ClientID)
	v.Set("client_secret", h.Cfg.GoogleTokens.ClientSecret)
	buf.WriteString(v.Encode())
	url := buf.String()
	req, _ := http.NewRequest("POST", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return token, err
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &token)
	return token, nil
}

func (h *Handler) tokenToCall(token OauthToken) ([]byte, error) {
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
