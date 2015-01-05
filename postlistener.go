package github

import (
	"encoding/json"
	"net/http"
)

type Payload struct {
	Ref, After, Before, Compare string
	Created, Deleted, Forced    bool
	Commits                     []Commit
	HeadCommit                  Commit
	Repository                  Repository
	Pusher                      Author
}

type Repository struct {
	Id, Watchers, Stargazers, Forks, Size, OpenIssues, CreatedAt, PushedAt int
	Name, Url, Description, Language, MasterBranch                         string
	Fork, Private, HasIssues, HasDownloads, HasWiki                        bool
}

type Author struct {
	Name, Email, Username string
}

type Commit struct {
	Id, Message, Timestamp, Url string
	Distinct                    bool
	Author, Commiter            Author
	Added, Removed, Modified    []string
}

func ReceiveHooks(addr string) (<-chan Payload, <-chan error) {
	ch := make(chan Payload, 1)
	ech := make(chan error, 1)
	go func() {
		err := http.ListenAndServe(addr, http.HandlerFunc(pushHandler(ch)))
		ech <- err
	}()
	return ch, ech
}

func pushHandler(ch chan<- Payload) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			decoder := json.NewDecoder(r.Body)
			out := &Payload{}
			err := decoder.Decode(out)
			if err == nil {
				ch <- *out
			} else {
				http.Error(w, err.Error(), 500)
			}
		default:
		}
	}
}
