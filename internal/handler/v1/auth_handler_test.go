package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-eagle/eagle/pkg/errcode"
	"go.uber.org/mock/gomock"

	"github.com/go-eagle/eagle-layout/internal/mocks"
	"github.com/go-eagle/eagle-layout/internal/types"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestAuthHandlers(t *testing.T) {
	errDown := errors.New("service down")

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		setup      func(svc *mocks.MockUserService)
		wantStatus int
		wantCode   int
		wantUser   string
	}{
		{
			name:   "login success",
			method: http.MethodPost,
			path:   "/v2/auth/login",
			body:   `{"email":"alice@example.com","username":"alice","password":"password123"}`,
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().Login(gomock.Any(), types.LoginInput{
					Email:    "alice@example.com",
					Username: "alice",
					Password: "password123",
				}).Return(&types.LoginOutput{ID: 7, AccessToken: "tok"}, nil)
			},
			wantStatus: http.StatusOK,
			wantCode:   0,
			wantUser:   "tok",
		},
		{
			name:       "login rejects invalid json",
			method:     http.MethodPost,
			path:       "/v2/auth/login",
			body:       `{`,
			setup:      func(svc *mocks.MockUserService) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   errcode.ErrInvalidParam.Code(),
		},
		{
			name:   "login maps service errors",
			method: http.MethodPost,
			path:   "/v2/auth/login",
			body:   `{"email":"alice@example.com","password":"nope"}`,
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, errDown)
			},
			wantStatus: http.StatusBadRequest,
			wantCode:   errcode.ErrUnauthorized.Code(),
		},
		{
			name:   "register success",
			method: http.MethodPost,
			path:   "/v2/auth/register",
			body:   `{"username":"alice","password":"password123"}`,
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().Register(gomock.Any(), types.RegisterInput{
					Username: "alice",
					Password: "password123",
				}).Return(&types.RegisterOutput{ID: 9, Username: "alice"}, nil)
			},
			wantStatus: http.StatusOK,
			wantCode:   0,
			wantUser:   "alice",
		},
		{
			name:       "register rejects invalid json",
			method:     http.MethodPost,
			path:       "/v2/auth/register",
			body:       ``,
			setup:      func(svc *mocks.MockUserService) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   errcode.ErrInvalidParam.Code(),
		},
		{
			name:   "register maps service errors",
			method: http.MethodPost,
			path:   "/v2/auth/register",
			body:   `{"username":"alice","password":"password123"}`,
			setup: func(svc *mocks.MockUserService) {
				svc.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil, errDown)
			},
			wantStatus: http.StatusBadRequest,
			wantCode:   errcode.ErrUnauthorized.Code(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockUserService(ctrl)
			tt.setup(svc)

			router := gin.New()
			router.POST("/v2/auth/login", NewLoginHandler(svc).LoginHandler)
			router.POST("/v2/auth/register", NewRegisterHandler(svc).RegisterHandler)

			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d, body %s", rec.Code, tt.wantStatus, rec.Body.String())
			}

			var resp struct {
				Code int `json:"code"`
				Data struct {
					User json.RawMessage `json:"user"`
				} `json:"data"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("decode response: %v body %s", err, rec.Body.String())
			}
			if resp.Code != tt.wantCode {
				t.Fatalf("code = %d, want %d", resp.Code, tt.wantCode)
			}
			if tt.wantUser != "" && !strings.Contains(string(resp.Data.User), tt.wantUser) {
				t.Fatalf("user payload = %s, want it to contain %s", resp.Data.User, tt.wantUser)
			}
		})
	}
}
