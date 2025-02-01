package middleware

import (
	"context"
	"net/http"
)

func (m *middleware) AddUserSubFromHeaderToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sub := r.Header.Get("X-Sub")
		ctx := context.WithValue(r.Context(), SubKey, sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *middleware) AddUserNameFromHeaderToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.Header.Get("X-Name")
		ctx := context.WithValue(r.Context(), NameKey, name)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *middleware) AddUserNicknameFromHeaderToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nickname := r.Header.Get("X-Nickname")
		ctx := context.WithValue(r.Context(), NicknameKey, nickname)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *middleware) AddUserEmailFromHeaderToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.Header.Get("X-Email")
		ctx := context.WithValue(r.Context(), EmailKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *middleware) AddUserPictureFromHeaderToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		picture := r.Header.Get("X-Picture")
		ctx := context.WithValue(r.Context(), PictureKey, picture)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *middleware) GetUserSubFromContext(r *http.Request) string {
	sub := r.Context().Value(SubKey).(string)
	if sub == "" {
		sub = "no-x-sub"
	}
	return sub
}

func (m *middleware) GetUserNameFromContext(r *http.Request) string {
	name := r.Context().Value(NameKey).(string)
	if name == "" {
		name = "no-name"
	}
	return name
}

func (m *middleware) GetUserNicknameFromContext(r *http.Request) string {
	name := r.Context().Value(NicknameKey).(string)
	if name == "" {
		name = "no-x-nickname"
	}
	return name
}

func (m *middleware) GetUserEmailFromContext(r *http.Request) string {
	email := r.Context().Value(EmailKey).(string)
	if email == "" {
		email = "no-x-email"
	}
	return email
}

func (m *middleware) GetUserPictureFromContext(r *http.Request) string {
	picture := r.Context().Value(PictureKey).(string)
	if picture == "" {
		picture = "no-x-picture"
	}
	return picture
}
