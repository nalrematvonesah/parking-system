package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"

	"gateway-service/internal/auth"
	mw "gateway-service/internal/middleware"
)

func noopLogger() *zap.Logger { return zap.NewNop() }
func newAuth() *auth.Manager  { return auth.New("test-secret", time.Hour) }

func TestJWTAuth_MissingHeader(t *testing.T) {
	mgr := newAuth()
	h := mw.JWTAuth(mgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	mgr := newAuth()
	h := mw.JWTAuth(mgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-token")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestJWTAuth_ValidToken_PassesUserID(t *testing.T) {
	mgr := newAuth()
	token, _ := mgr.Issue(55)
	var gotUID int64
	h := mw.JWTAuth(mgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUID = mw.UserIDFrom(r.Context())
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if gotUID != 55 {
		t.Errorf("expected uid=55, got %d", gotUID)
	}
}

func TestJWTAuth_NoBearerPrefix(t *testing.T) {
	mgr := newAuth()
	token, _ := mgr.Issue(1)
	h := mw.JWTAuth(mgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestJWTAuth_ExpiredToken(t *testing.T) {
	expiredMgr := auth.New("test-secret", -time.Second)
	validMgr := auth.New("test-secret", time.Hour)
	token, _ := expiredMgr.Issue(1)
	h := mw.JWTAuth(validMgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for expired token, got %d", rec.Code)
	}
}

func TestUserIDFrom_NoAuth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	uid := mw.UserIDFrom(req.Context())
	if uid != 0 {
		t.Errorf("expected 0 without auth, got %d", uid)
	}
}

func TestRecovery_HandlesPanic(t *testing.T) {
	h := mw.Recovery(noopLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 on panic, got %d", rec.Code)
	}
}

func TestAccessLog_DoesNotBlockRequest(t *testing.T) {
	h := mw.AccessLog(noopLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/test", nil))
	if rec.Code != http.StatusAccepted {
		t.Errorf("expected 202, got %d", rec.Code)
	}
}
