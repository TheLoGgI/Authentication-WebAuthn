package providers

import (
	"fmt"
	"net/http"
)

func RefreshAuthToken(w http.ResponseWriter, r *http.Request) {

	// CORS Headers
	w.Header().Set("Content-Type", "application/json")

	// getting isAdmin from context and convert to bool
	// adm := context.Get(r, "isAdmin").(bool)
	adm := true

	// creating response, depending on isAdmin status
	body := "<h1>Hello on secret route.</h1>\n<p>%s</p>"
	var response string
	if adm {
		response = fmt.Sprintf(body, "You are admin.")
	} else {
		response = fmt.Sprintf(body, "You are user.")
	}

	fmt.Fprintln(w, response)
}
