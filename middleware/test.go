package middleware

import (
	"net/http"
)

func TestLoginMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read basic auth information
		// usr, hashPassword, ok := r.BasicAuth()

		// if there is no basic auth (no matter which credentials)
		// if !ok {
		// 	errMsg := "Authentication error!"
		// 	// return a 403 forbidden
		// 	http.Error(w, errMsg, http.StatusForbidden)
		// 	log.Println(errMsg)

		// 	// stop processing route
		// 	return
		// }

		// let's assume we check the user against a database to get
		// his admin-right and put this to the request context
		// context.Set(r, "isAdmin", true)
		// var isUserValid = database.ValidateDatabaseUser(usr, []byte(hashPassword))
		// fmt.Printf("User is valid %v \n", isUserValid)
		// if !isUserValid {
		// 	errMsg := "Authentication Failed"
		// 	// return a 403 forbidden
		// 	http.Error(w, errMsg, http.StatusForbidden)
		// 	log.Println(errMsg)
		// 	return
		// }

		// else continue processing
		// log.Printf("User %s logged in.", usr)
		next(w, r)
	}
}
