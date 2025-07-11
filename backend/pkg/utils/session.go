package utils

import (
	"net/http"
	"time"

	"social-network/pkg/models"
)

type contextKey string

// key for using context / accessing user_id
var UserKey = contextKey("UserID")

// session cookie name
const sessionCookie = "session-id"

// Standart cookie lifespan
const cookieLifespan = 3600 * 12 // 12h

/* ---- session and cookie funcionality communicating with client request --- */
/* -------------------------------------------------------------------------- */
/*                                   session                                  */
/* -------------------------------------------------------------------------- */

// Creates creates sid, http cookie, sends it to client and returns session
func SessionStart(w http.ResponseWriter, r *http.Request, userID string) models.Session {
	sessionID := UniqueId()
	// create cookie
	cookie := CreateCookie(sessionID, cookieLifespan)
	// create session
	session := models.Session{
		ID:             sessionID,
		UserID:         userID,
		ExpirationTime: time.Now().Add(30 * time.Minute),
	}
	// Send cookie to client
	http.SetCookie(w, &cookie)
	// return session
	return session
}

// Returns true if session time is  not expired
func CheckSessionExpiration(session models.Session) bool {
	return session.ExpirationTime.After(time.Now())
}

/* -------------------------------------------------------------------------- */
/*                                   cookie                                   */
/* -------------------------------------------------------------------------- */

// session cookie blueprint
func CreateCookie(sessionID string, lifespan int) http.Cookie {
	return http.Cookie{
		Name:     sessionCookie,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: false,
		MaxAge:   lifespan,
		SameSite: http.SameSiteLaxMode, // Pour le dev local, change en SameSiteNoneMode + Secure: true en prod
		Secure:   false,                // Passe à true en prod avec HTTPS
	}
}

// get cookie from web
func GetCookie(r *http.Request) (string, error) {
	cookieFromWeb, err := r.Cookie(sessionCookie)
	if err != nil {
		return "", err
	}
	cookieValue := cookieFromWeb.Value
	if len(cookieValue) == 0 {
		return "", err
	}
	return cookieValue, nil
}

// Delete cookie
func DeleteCookie(w http.ResponseWriter) {
	cookie := CreateCookie(" ", -1)
	http.SetCookie(w, &cookie)
}
