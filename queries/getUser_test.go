package queries

import (
	"testing"
)

func TestGetUser(t *testing.T) {

	var userUid = "562e9fee-66e5-44c6-b625-5150c8c3368d"

	user, failedFetchUser := GetUser(userUid)
	if failedFetchUser != nil {
		t.Errorf("GetUser(\"%v\") FAILED DATABASE FETCH", userUid)
		return
	}

	if user.Uid.String() != "" {
		t.Log("GetUser PASSED")
	}
	t.Log("GetUser PASSED")

}
