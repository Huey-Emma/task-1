package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type PersonalInfo struct {
		SlackName 		string 	   `json:"slack_name"`
		ExampleName 	string 	   `json:"example_name"`
		CurrentDay 		string 	   `json:"current_day"`
		UTCTime 		time.Time  `json:"utc_time"`
		GithubFileURL 	string 	   `json:"github_file_url"`
		GithubRepoURL 	string 	   `json:"github_repo_url"`
		StatusCode 		int 	   `json:"status_code"`
}	

func day(t time.Time) string {
		fmt := t.Format("2006-01-02 15:04:05 Monday")
		parts := strings.Split(fmt, " ")
		return parts[len(parts) - 1]
}

func validstring(s string) bool {
		return len(strings.TrimSpace(s)) > 0
}

func queryParam(v url.Values, key string) string {
		return v.Get(key)
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(code)
		return json.NewEncoder(w).Encode(v)
}

type check struct {
		field string 
		cond  bool 
		msg   string
}

type validationError struct {
		Field  string `json:"field"`
		ErrMsg string `json:"errmsg"`
}

func (e validationError) Error() string {
		return e.ErrMsg
}

func validate(checks ...check) []validationError {
		errs := make([]validationError, 0)

		for _, chk := range checks {
				if !chk.cond {
						errs = append(errs, validationError{
								Field:  chk.field,
								ErrMsg: chk.msg,
						})
				}
		}

		if len(errs) == 0 {
				return nil
		}

		return errs
} 

func infoHandler(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query() 

		slackname := queryParam(query, "slack_name")
		exampName := queryParam(query, "example_name")

		checks := []check{
				{"slack_name", validstring(slackname), "slack_name cannot be blank"},
				{"example_name", validstring(exampName), "example_name cannot be blank"},
		}

		if errs := validate(checks...); errs != nil {
				writeJSON(w, http.StatusUnprocessableEntity, errs)
				return 
		}

		personalInfo := PersonalInfo{
				SlackName:     slackname,
				ExampleName:   exampName,
				CurrentDay:    day(time.Now()),
				UTCTime:       time.Now(),
				GithubFileURL: "https://github.com/Huey-Emma/task-1/blob/main/main.go",
				GithubRepoURL: "https://github.com/Huey-Emma/task-1",
				StatusCode:    http.StatusOK,
		}

		writeJSON(w, http.StatusOK, personalInfo)
}

func main() {
		mux := http.NewServeMux()

		mux.HandleFunc("/api", infoHandler)

		err := http.ListenAndServe("0.0.0.0:4000", mux)

		if err != nil {
				log.Fatal(err)
		}
}

