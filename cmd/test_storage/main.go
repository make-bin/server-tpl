package main

import (
	"context"
	"fmt"
	"log"

	"github.com/make-bin/server-tpl/pkg/domain/model"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore/factory"
)

func main() {
	// 创建内存存储配置
	config := &datastore.Config{
		Type:     "memory",
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		Database: "testdb",
		SSLMode:  "disable",
		MaxIdle:  10,
		MaxOpen:  100,
		Timeout:  30,
	}

	// 创建数据存储实例
	store, err := factory.NewDataStore(config)
	if err != nil {
		log.Fatalf("Failed to create data store: %v", err)
	}

	// 连接数据库
	ctx := context.Background()
	if err := store.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer store.Disconnect(ctx)

	fmt.Println("=== 测试用户存储 ===")
	testUserStorage(ctx, store)

	fmt.Println("\n=== 测试应用存储 ===")
	testApplicationStorage(ctx, store)

	fmt.Println("\n=== 测试变量存储 ===")
	testVariableStorage(ctx, store)

	fmt.Println("\n=== 所有测试完成 ===")
}

func testUserStorage(ctx context.Context, store datastore.DataStore) {
	// 创建用户
	user := &model.User{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
		Role:     "user",
		Status:   "active",
	}

	// 添加用户
	if err := store.Add(ctx, user); err != nil {
		log.Printf("Failed to add user: %v", err)
		return
	}
	fmt.Printf("User created with ID: %d\n", user.BaseEntity.ID)

	// 查询用户
	queryUser := &model.User{}
	queryUser.BaseEntity.ID = user.BaseEntity.ID
	if err := store.Get(ctx, queryUser); err != nil {
		log.Printf("Failed to get user: %v", err)
		return
	}
	fmt.Printf("Retrieved user: %s (%s)\n", queryUser.Name, queryUser.Email)

	// 更新用户
	queryUser.Name = "Updated Test User"
	if err := store.Put(ctx, queryUser); err != nil {
		log.Printf("Failed to update user: %v", err)
		return
	}
	fmt.Println("User updated successfully")

	// 列出用户
	listOptions := &datastore.ListOptions{
		Page:     1,
		PageSize: 10,
		SortBy: []datastore.SortOption{
			{Key: "id", Order: datastore.SortOrderAscending},
		},
	}

	users, err := store.List(ctx, &model.User{}, listOptions)
	if err != nil {
		log.Printf("Failed to list users: %v", err)
		return
	}
	fmt.Printf("Found %d users\n", len(users))

	// 统计用户数量
	count, err := store.Count(ctx, &model.User{}, nil)
	if err != nil {
		log.Printf("Failed to count users: %v", err)
		return
	}
	fmt.Printf("Total users: %d\n", count)

	// 删除用户
	if err := store.Delete(ctx, queryUser); err != nil {
		log.Printf("Failed to delete user: %v", err)
		return
	}
	fmt.Println("User deleted successfully")

	// 验证删除
	exists, err := store.IsExist(ctx, queryUser)
	if err != nil {
		log.Printf("Failed to check user existence: %v", err)
		return
	}
	fmt.Printf("User exists: %t\n", exists)
}

func testApplicationStorage(ctx context.Context, store datastore.DataStore) {
	// 创建应用
	app := &model.Application{
		Name:        "Test App",
		Description: "A test application",
		Version:     "1.0.0",
		Status:      "active",
		CreatedBy:   1,
	}

	// 添加应用
	if err := store.Add(ctx, app); err != nil {
		log.Printf("Failed to add application: %v", err)
		return
	}
	fmt.Printf("Application created with ID: %d\n", app.BaseEntity.ID)

	// 查询应用
	queryApp := &model.Application{}
	queryApp.BaseEntity.ID = app.BaseEntity.ID
	if err := store.Get(ctx, queryApp); err != nil {
		log.Printf("Failed to get application: %v", err)
		return
	}
	fmt.Printf("Retrieved application: %s (%s)\n", queryApp.Name, queryApp.Version)

	// 更新应用
	queryApp.Description = "Updated test application"
	if err := store.Put(ctx, queryApp); err != nil {
		log.Printf("Failed to update application: %v", err)
		return
	}
	fmt.Println("Application updated successfully")

	// 列出应用
	listOptions := &datastore.ListOptions{
		Page:     1,
		PageSize: 10,
		SortBy: []datastore.SortOption{
			{Key: "id", Order: datastore.SortOrderAscending},
		},
	}

	apps, err := store.List(ctx, &model.Application{}, listOptions)
	if err != nil {
		log.Printf("Failed to list applications: %v", err)
		return
	}
	fmt.Printf("Found %d applications\n", len(apps))

	// 删除应用
	if err := store.Delete(ctx, queryApp); err != nil {
		log.Printf("Failed to delete application: %v", err)
		return
	}
	fmt.Println("Application deleted successfully")
}

func testVariableStorage(ctx context.Context, store datastore.DataStore) {
	// 创建变量
	variable := &model.Variable{
		ApplicationID: 1,
		Key:           "TEST_KEY",
		Value:         "test_value",
		Description:   "A test variable",
		Type:          "string",
		IsSecret:      false,
	}

	// 添加变量
	if err := store.Add(ctx, variable); err != nil {
		log.Printf("Failed to add variable: %v", err)
		return
	}
	fmt.Printf("Variable created with ID: %d\n", variable.BaseEntity.ID)

	// 查询变量
	queryVariable := &model.Variable{}
	queryVariable.BaseEntity.ID = variable.BaseEntity.ID
	if err := store.Get(ctx, queryVariable); err != nil {
		log.Printf("Failed to get variable: %v", err)
		return
	}
	fmt.Printf("Retrieved variable: %s = %s\n", queryVariable.Key, queryVariable.Value)

	// 更新变量
	queryVariable.Value = "updated_test_value"
	if err := store.Put(ctx, queryVariable); err != nil {
		log.Printf("Failed to update variable: %v", err)
		return
	}
	fmt.Println("Variable updated successfully")

	// 列出变量
	listOptions := &datastore.ListOptions{
		Page:     1,
		PageSize: 10,
		SortBy: []datastore.SortOption{
			{Key: "id", Order: datastore.SortOrderAscending},
		},
	}

	variables, err := store.List(ctx, &model.Variable{}, listOptions)
	if err != nil {
		log.Printf("Failed to list variables: %v", err)
		return
	}
	fmt.Printf("Found %d variables\n", len(variables))

	// 删除变量
	if err := store.Delete(ctx, queryVariable); err != nil {
		log.Printf("Failed to delete variable: %v", err)
		return
	}
	fmt.Println("Variable deleted successfully")
}
