package main

import (
	"github.com/go-advanced-admin/admin"
	"github.com/go-advanced-admin/web-echo"
	"log"
)

type TestModel1 struct {
	Name string
}

type TestModel2 struct {
	Name string
}

type TestModel3 struct {
	Name string
}

type TestModel4 struct {
	Name string
}

func main() {
	e := echo.New()

	web := adminecho.NewIntegrator(e.Group(""))

	permissionFunc := func(req admin.PermissionRequest, ctx interface{}) (bool, error) {
		return true, nil
	}

	panel := admin.NewPanel(nil, web, permissionFunc, nil)

	testApp1, err := panel.RegisterApp("Test App 1")
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

	testApp2, err := panel.RegisterApp("Test App 2")
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
