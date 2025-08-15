package e2e

import (
	"net/http"
	"testing"
	"time"
)

func TestHealthEndpoint(t *testing.T) {
	// 这里可以添加端到端测试
	// 目前只是一个示例
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 测试健康检查端点
	resp, err := client.Get("http://localhost:8080/health")
	if err != nil {
		t.Skipf("Server not running, skipping test: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
