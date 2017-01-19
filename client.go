package main


import (
	"fmt"
	"net/http"
	"math/rand"

	"stathat.com/c/jconfig"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	githuboauth "golang.org/x/oauth2/github"
)


const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const htmlIndex   = `<html><body>
Logged in with <a href="/login">GitHub</a>
</body></html>
`

var conf = jconfig.LoadConfig("config/config.json")
var (
	// You must register the app at https://github.com/settings/applications
	// Set callback to http://127.0.0.1:7000/github_oauth_cb
	oauthConf = &oauth2.Config{
		ClientID:     conf.GetString("ClientID"),
		ClientSecret: conf.GetString("ClientSecret"),
		Scopes:       []string{"user:email", "repo"},
		Endpoint:     githuboauth.Endpoint,
	}
	// random string for oauth2 API calls to protect against CSRF
	oauthStateString = genRandString(32)
)

// generate a random string
func genRandString(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}


// main handler
func handleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlIndex))
}


// login handler
func handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}


// callback handler , github_oauth_cb. Called by GH after Auth is granted
func handleGitHubCallback(w http.ResponseWriter, r *http.Request) {

	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)

	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)

	user, _, err := client.Users.Get("")

	if err != nil {
		fmt.Printf("client.Users.Get() faled with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Printf("Logged in as GitHub user: %s\n", *user.Login)
	fmt.Printf("\n%v\n", github.Stringify(user))

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}


func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGitHubLogin)
	http.HandleFunc("/github_oauth_cb", handleGitHubCallback)

	fmt.Print("Started running on http://127.0.0.1:7000\n")
	fmt.Println(http.ListenAndServe(":7000", nil))
}