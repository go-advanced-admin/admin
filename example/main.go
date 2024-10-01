package main

import (
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
	ID   uint `gorm:"primarykey"`
	Name string
}

type TestModel2 struct {
	ID   uint `gorm:"primarykey"`
	Name string
}

type TestModel3 struct {
	ID   uint `gorm:"primarykey"`
	Name string
}

type TestModel4 struct {
	ID   uint `gorm:"primarykey"`
	Name string
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	web := adminecho.NewIntegrator(e.Group(""))

	permissionFunc := func(req admin.PermissionRequest, ctx interface{}) (bool, error) {
		return true, nil
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&TestModel1{}, &TestModel2{}, &TestModel3{}, &TestModel4{})
	if err != nil {
		log.Fatal(err)
	}

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
