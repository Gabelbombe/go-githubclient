package main

import (
	"fmt"
	"net/http"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)


const htmlIndex = `<html><body>
Logged in with <a href="/login">GitHub</a>
</body></html>
`