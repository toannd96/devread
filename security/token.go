package security

import "net/http"

func ExtractTokenMail(r *http.Request) string {
	token := r.URL.Query().Get("token")
	return token
}
