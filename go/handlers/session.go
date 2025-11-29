package handlers

import (
	"net/http"
)

func GetUserFromSession(r *http.Request) (map[string]string, bool) {
	c, err := r.Cookie("session")
	if err != nil || c.Value == "" {
		return nil, false
	}

	value := map[string]string{}
	if err := s.Decode("session", c.Value, &value); err != nil {
		return nil, false
	}

	return value, true
}
