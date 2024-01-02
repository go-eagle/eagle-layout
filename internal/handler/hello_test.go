package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHello(t *testing.T) {
	tests := []struct {
		name         string
		queryParams  string
		expectedCode int
		expectedBody string
	}{
		{"ValidParam", "name=eagle", http.StatusOK, `{"code":0,"message":"Ok","data":{"result":"hello eagle"}}`},
		{"InvalidParam", "email=abc", http.StatusOK, `{"code":0,"message":"Ok","data":{"result":"hello "}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建一个虚拟的 HTTP 请求
			req, err := http.NewRequest("GET", "/hello?"+tt.queryParams, nil)
			if err != nil {
				t.Fatal(err)
			}

			// 创建一个 ResponseRecorder 来记录响应
			recorder := httptest.NewRecorder()

			// 创建一个 Gin 上下文
			context, _ := gin.CreateTestContext(recorder)
			context.Request = req

			// 调用被测试的处理函数
			Hello(context)

			// 检查响应状态码
			assert.Equal(t, tt.expectedCode, recorder.Code)

			// 检查响应体
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(recorder.Body.String()))
		})
	}
}
