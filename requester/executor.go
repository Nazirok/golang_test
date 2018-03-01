package requester

import (
	"github.com/golang_test/handler"
	"net/http"
)

func RequestIssueExecutor (issues chan handler.RequestToChannel) {
	client := &http.Client{}
	for issue := range issues {
		resp, err := client.Do(issue.NewReq)
		defer resp.Body.Close()
		if err != nil {
			//http.Error(w, "Error during request", http.StatusInternalServerError)
			return
		}
	}

}
