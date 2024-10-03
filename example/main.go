package main

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/go-advanced-admin/admin"
	"github.com/go-advanced-admin/orm-gorm"
	"github.com/go-advanced-admin/web-echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"log"
)

type TestModel1 struct {
	ID   uint `gorm:"primarykey" admin:"editForm:exclude"`
	Name string
}

type TestModel2 struct {
	ID   uint `gorm:"primarykey"`
	Name string
}

type TestModel3 struct {
	ID       uint `gorm:"primarykey"`
	Username string
}

type TestModel4 struct {
	ID    uint `gorm:"primarykey"`
	Title string
}

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} error=${error}\n",
	}))

	web := adminecho.NewIntegrator(e.Group(""))

	permissionFunc := func(req admin.PermissionRequest, ctx interface{}) (bool, error) {
		return true, nil
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&TestModel1{}, &TestModel2{}, &TestModel3{}, &TestModel4{})
	if err != nil {
		log.Fatal(err)
	}

	populateTestModels(db)

	orm := admingorm.NewIntegrator(db)

	panel, err := admin.NewPanel(orm, web, permissionFunc, nil)
	if err != nil {
		log.Fatal(err)
	}

	testApp1, err := panel.RegisterApp("TestApp1", "Test App 1")
	if err != nil {
		log.Fatal(err)
	}

	_, err = testApp1.RegisterModel(&TestModel1{})
	if err != nil {
		log.Fatal(err)
	}
	_, err = testApp1.RegisterModel(&TestModel2{})
	if err != nil {
		log.Fatal(err)
	}

	testApp2, err := panel.RegisterApp("TestApp2", "Test App 2")
	if err != nil {
		log.Fatal(err)
	}

	_, err = testApp2.RegisterModel(&TestModel3{})
	if err != nil {
		log.Fatal(err)
	}

	_, err = testApp2.RegisterModel(&TestModel4{})
	if err != nil {
		log.Fatal(err)
	}

	for _, route := range e.Routes() {
		log.Println(route.Method, route.Path)
	}

	e.Logger.Fatal(e.Start(":8080"))
}

func populateTestModels(db *gorm.DB) {
	for i := 0; i < 100; i++ {
		db.Create(&TestModel1{
			Name: fmt.Sprintf("Name%d", i),
		})
	}

	for i := 0; i < 100; i++ {
		db.Create(&TestModel2{
			Name: fmt.Sprintf("Product%d", i),
		})
	}

	for i := 0; i < 100; i++ {
		db.Create(&TestModel3{
			Username: fmt.Sprintf("user%d", i),
		})
	}

	for i := 0; i < 100; i++ {
		db.Create(&TestModel4{
			Title: fmt.Sprintf("Post Title %d", i),
		})
	}
}
