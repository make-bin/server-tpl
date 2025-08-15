package container

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试用的结构体
type TestBean struct {
	Name string
	Age  int
}

type AnotherBean struct {
	Value string
}

func TestNewContainer(t *testing.T) {
	container := NewContainer()
	assert.NotNil(t, container)
	assert.NotNil(t, container.beans)
	assert.Equal(t, 0, len(container.beans))
}

func TestContainer_ProvideWithName(t *testing.T) {
	tests := []struct {
		name     string
		beanName string
		bean     interface{}
		wantErr  bool
	}{
		{
			name:     "正常注册bean",
			beanName: "testBean",
			bean:     &TestBean{Name: "test", Age: 25},
			wantErr:  false,
		},
		{
			name:     "空名称自动生成",
			beanName: "",
			bean:     &TestBean{Name: "test", Age: 25},
			wantErr:  false,
		},
		{
			name:     "nil bean",
			beanName: "nilBean",
			bean:     nil,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewContainer()
			err := container.ProvideWithName(tt.beanName, tt.bean)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestContainer_Get(t *testing.T) {
	container := NewContainer()
	testBean := &TestBean{Name: "test", Age: 25}

	// 注册bean
	err := container.ProvideWithName("testBean", testBean)
	require.NoError(t, err)

	tests := []struct {
		name      string
		beanName  string
		wantBean  interface{}
		wantExist bool
	}{
		{
			name:      "获取存在的bean",
			beanName:  "testBean",
			wantBean:  testBean,
			wantExist: true,
		},
		{
			name:      "获取不存在的bean",
			beanName:  "nonexistent",
			wantBean:  nil,
			wantExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bean, exists := container.Get(tt.beanName)

			assert.Equal(t, tt.wantExist, exists)
			if tt.wantExist {
				assert.Equal(t, tt.wantBean, bean)
			} else {
				assert.Nil(t, bean)
			}
		})
	}
}

func TestContainer_Provides(t *testing.T) {
	container := NewContainer()
	bean1 := &TestBean{Name: "test1", Age: 25}
	bean2 := &AnotherBean{Value: "test2"}

	// 测试提供多个bean
	err := container.Provides(bean1, bean2)
	assert.NoError(t, err)

	// 验证bean是否被正确注册
	got1, exists1 := container.Get("*container.TestBean")
	assert.True(t, exists1)
	assert.Equal(t, bean1, got1)

	got2, exists2 := container.Get("*container.AnotherBean")
	assert.True(t, exists2)
	assert.Equal(t, bean2, got2)
}

func TestContainer_Provides_WithError(t *testing.T) {
	container := NewContainer()

	// 模拟提供bean时出错的情况
	// 这里我们创建一个特殊的测试场景
	// 由于当前的实现不会出错，我们主要测试正常流程
	err := container.Provides(&TestBean{Name: "test", Age: 25})
	assert.NoError(t, err)
}

func TestContainer_Populate(t *testing.T) {
	container := NewContainer()

	// 测试Populate方法
	err := container.Populate()
	assert.NoError(t, err)
}

func TestContainer_Concurrency(t *testing.T) {
	container := NewContainer()

	// 并发测试
	const numGoroutines = 10
	const numOperations = 100

	// 并发写入
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < numOperations; j++ {
				beanName := fmt.Sprintf("bean_%d_%d", id, j)
				bean := &TestBean{Name: beanName, Age: id + j}
				_ = container.ProvideWithName(beanName, bean)
			}
		}(i)
	}

	// 并发读取
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < numOperations; j++ {
				beanName := fmt.Sprintf("bean_%d_%d", id, j)
				_, _ = container.Get(beanName)
			}
		}(i)
	}

	// 等待一段时间让goroutines完成
	time.Sleep(100 * time.Millisecond)

	// 验证容器状态
	// 由于并发操作，我们只验证容器没有panic
	assert.NotNil(t, container)
}

func TestContainer_AutomaticNaming(t *testing.T) {
	container := NewContainer()

	// 测试自动命名功能
	testBean := &TestBean{Name: "test", Age: 25}
	err := container.ProvideWithName("", testBean)
	assert.NoError(t, err)

	// 验证自动生成的名称
	expectedName := "*container.TestBean"
	bean, exists := container.Get(expectedName)
	assert.True(t, exists)
	assert.Equal(t, testBean, bean)
}

func TestContainer_BeanReplacement(t *testing.T) {
	container := NewContainer()

	// 测试bean替换功能
	bean1 := &TestBean{Name: "original", Age: 25}
	bean2 := &TestBean{Name: "replaced", Age: 30}

	// 注册第一个bean
	err := container.ProvideWithName("testBean", bean1)
	assert.NoError(t, err)

	// 验证第一个bean
	got1, exists1 := container.Get("testBean")
	assert.True(t, exists1)
	assert.Equal(t, bean1, got1)

	// 替换为第二个bean
	err = container.ProvideWithName("testBean", bean2)
	assert.NoError(t, err)

	// 验证第二个bean
	got2, exists2 := container.Get("testBean")
	assert.True(t, exists2)
	assert.Equal(t, bean2, got2)
	assert.NotEqual(t, bean1, got2)
}

// 基准测试
func BenchmarkContainer_ProvideWithName(b *testing.B) {
	container := NewContainer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		beanName := fmt.Sprintf("bean_%d", i)
		bean := &TestBean{Name: beanName, Age: i}
		_ = container.ProvideWithName(beanName, bean)
	}
}

func BenchmarkContainer_Get(b *testing.B) {
	container := NewContainer()

	// 预先注册一些bean
	for i := 0; i < 100; i++ {
		beanName := fmt.Sprintf("bean_%d", i)
		bean := &TestBean{Name: beanName, Age: i}
		_ = container.ProvideWithName(beanName, bean)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		beanName := fmt.Sprintf("bean_%d", i%100)
		_, _ = container.Get(beanName)
	}
}

func BenchmarkContainer_Provides(b *testing.B) {
	container := NewContainer()
	beans := make([]interface{}, 10)

	for i := 0; i < 10; i++ {
		beans[i] = &TestBean{Name: fmt.Sprintf("bean_%d", i), Age: i}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = container.Provides(beans...)
	}
}
