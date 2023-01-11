package models

import (
	"strings"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	FirstName   string                `bson:"firstName,omitempty"`
	LastName    string                `bson:"lastName,omitempty"`
	Username    string                `bson:"username,omitempty"`
	Email       string                `json:"email" bson:"email, omitempty"`
	Password    string                `json:"password" bson:"password, omitempty"`
	Uid         uuid.UUID             `json:"uid" bson:"uid,omitempty"`
	EntryId     primitive.ObjectID    `bson:"_id,omitempty"`
	Credentials []webauthn.Credential `json:"credentials" bson:"credentials, omitempty"`
	Credential  webauthn.Credential   `json:"credential" bson:"credential, omitempty"`
	Session     webauthn.SessionData  `json:"session" bson:"session, omitempty"`
	// webauthn.User
	// AuthUser  webauthn.User      `json:"authUser" bson:"authUser, omitempty"`
}

type NewUserAccountRequest struct {
	Username string
	Email    string
	Password []byte
	Uid      uuid.UUID
}

func (user User) WebAuthnID() []byte {
	return []byte(user.Uid.String())
}

func (user User) WebAuthnName() string {
	return user.Username
}

func (user User) WebAuthnDisplayName() string {
	if user.FirstName != "" || user.LastName != "" {
		return strings.TrimSpace(user.FirstName + " " + user.LastName)
	}
	return user.Username
}

func (user User) WebAuthnIcon() string {
	return ""
	// return fmt.Sprintf("https://avatars.dicebear.com/api/adventurer/%v.svg", user.Uid)
}

// WebAuthnCredentials returns credentials owned by the user
func (user User) WebAuthnCredentials() []webauthn.Credential {
	if len(user.Credentials) == 0 {
		return []webauthn.Credential{user.Credential}
	}

	return user.Credentials
}

// AddCredential associates the credential to the user
func (user *User) AddCredential(cred *webauthn.Credential) {
	user.Credentials = append(user.Credentials, *cred)
}
