package common

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCacheClient 模拟缓存客户端
type MockCacheClient struct {
	mock.Mock
}

func (m *MockCacheClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheClient) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheClient) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheClient) Flush(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestCacheHelper_GetJSON(t *testing.T) {
	// 创建模拟缓存客户端
	mockClient := new(MockCacheClient)
	cacheHelper := &CacheHelper{cache: mockClient}

	ctx := context.Background()
	testKey := "test_key"
	testValue := map[string]interface{}{"name": "test", "value": 123}

	// 模拟缓存命中
	mockClient.On("Get", ctx, testKey).Return(`{"name":"test","value":123}`, nil)

	var result map[string]interface{}
	err := cacheHelper.GetJSON(ctx, testKey, &result)

	assert.NoError(t, err)
	assert.Equal(t, testValue, result)
	mockClient.AssertExpectations(t)
}

func TestCacheHelper_GetJSON_NotFound(t *testing.T) {
	mockClient := new(MockCacheClient)
	cacheHelper := &CacheHelper{cache: mockClient}

	ctx := context.Background()
	testKey := "non_existent_key"

	// 模拟缓存未找到
	mockClient.On("Get", ctx, testKey).Return("", ErrCacheNotFound)

	var result map[string]interface{}
	err := cacheHelper.GetJSON(ctx, testKey, &result)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "缓存不存在")
	mockClient.AssertExpectations(t)
}

func TestCacheHelper_SetJSON(t *testing.T) {
	mockClient := new(MockCacheClient)
	cacheHelper := &CacheHelper{cache: mockClient}

	ctx := context.Background()
	testKey := "test_key"
	testValue := map[string]interface{}{"name": "test", "value": 123}
	expiration := 30 * time.Minute

	// 模拟成功设置缓存
	mockClient.On("Set", ctx, testKey, `{"name":"test","value":123}`, expiration).Return(nil)

	err := cacheHelper.SetJSON(ctx, testKey, testValue, expiration)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestCacheHelper_GetOrSet_CacheHit(t *testing.T) {
	mockClient := new(MockCacheClient)
	cacheHelper := &CacheHelper{cache: mockClient}

	ctx := context.Background()
	testKey := "test_key"
	testValue := map[string]interface{}{"name": "test", "value": 123}

	// 模拟缓存命中
	mockClient.On("Get", ctx, testKey).Return(`{"name":"test","value":123}`, nil)

	fn := func() (interface{}, error) {
		return map[string]interface{}{"should": "not", "be": "called"}, nil
	}

	result, err := cacheHelper.GetOrSet(ctx, testKey, 30*time.Minute, fn)

	assert.NoError(t, err)
	assert.Equal(t, testValue, result)
	mockClient.AssertExpectations(t)
}

func TestCacheHelper_GetOrSet_CacheMiss(t *testing.T) {
	mockClient := new(MockCacheClient)
	cacheHelper := &CacheHelper{cache: mockClient}

	ctx := context.Background()
	testKey := "test_key"
	testValue := map[string]interface{}{"name": "test", "value": 123}

	// 模拟缓存未命中，然后设置缓存
	mockClient.On("Get", ctx, testKey).Return("", ErrCacheNotFound)
	mockClient.On("Set", ctx, testKey, `{"name":"test","value":123}`, mock.AnythingOfType("time.Duration")).Return(nil)

	fn := func() (interface{}, error) {
		return testValue, nil
	}

	result, err := cacheHelper.GetOrSet(ctx, testKey, 30*time.Minute, fn)

	assert.NoError(t, err)
	assert.Equal(t, testValue, result)
	mockClient.AssertExpectations(t)
}

func TestCacheHelper_GetOrSet_FnError(t *testing.T) {
	mockClient := new(MockCacheClient)
	cacheHelper := &CacheHelper{cache: mockClient}

	ctx := context.Background()
	testKey := "test_key"
	expectedErr := ErrCacheNotFound

	// 模拟缓存未命中，函数返回错误
	mockClient.On("Get", ctx, testKey).Return("", expectedErr)

	fn := func() (interface{}, error) {
		return nil, expectedErr
	}

	result, err := cacheHelper.GetOrSet(ctx, testKey, 30*time.Minute, fn)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
	mockClient.AssertExpectations(t)
}

func TestGetUserCacheKey(t *testing.T) {
	userID := "123"
	suffix := "profile"

	result := GetUserCacheKey(userID, suffix)

	expected := "user:123:profile"
	assert.Equal(t, expected, result)
}

func TestGetOptionCacheKey(t *testing.T) {
	optionType := "industry"
	suffix := "list"

	result := GetOptionCacheKey(optionType, suffix)

	expected := "options:industry:list"
	assert.Equal(t, expected, result)
}

// 基准测试
func BenchmarkCacheHelper_SetJSON(b *testing.B) {
	// 注意：这个基准测试需要真实的 Redis 连接才能运行
	// 这里仅作为示例，实际测试时需要配置测试环境
	b.Skip("需要真实的 Redis 连接")
}

func BenchmarkCacheHelper_GetJSON(b *testing.B) {
	// 注意：这个基准测试需要真实的 Redis 连接才能运行
	b.Skip("需要真实的 Redis 连接")
}

// 集成测试标记
func TestCacheIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 这里可以添加集成测试逻辑
	// 需要配置测试用的 Redis 实例
	t.Skip("需要配置测试 Redis 实例")
}
