# Go Advanced Admin

A Highly Customizable Advanced Admin Panel for Go.

[![Go Report Card](https://goreportcard.com/badge/github.com/go-advanced-admin/admin)](https://goreportcard.com/report/github.com/go-advanced-admin/admin)
[![Go](https://github.com/go-advanced-admin/admin/actions/workflows/tests.yml/badge.svg)](https://github.com/go-advanced-admin/admin/actions/workflows/tests.yml)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/go-advanced-admin/admin?tab=doc)

Go Advanced Admin is a powerful and highly customizable admin panel for Go applications. It allows developers to 
quickly create admin interfaces with minimal configuration, supporting multiple web frameworks and ORMs.

## Features

- **Framework Agnostic**: Compatible with popular Go web frameworks like Gin, Echo, Chi, Fiber, and more.
- **ORM Support**: Integrates seamlessly with ORMs such as GORM, XORM, SQLX, Bun, etc.
- **Customizable Templates**: Override default templates or create your own for complete control over the admin UI.
- **Fine-Grained Permissions**: Implement custom permission schemes (role-based, attribute-based) for robust access 
control.
- **Extensible**: Easily extend functionality with custom modules, themes, and widgets.
- **Logging**: Configurable logging system with support for custom log stores.

## Installation

Add the module to your project by running:

```sh
go get github.com/go-advanced-admin/admin
```

## Documentation

For detailed documentation, quick start guides, and advanced topics, please visit the 
[official documentation website](https://goadmin.dev).

## Quick Start

Here's a minimal example to get you started:

```go
package main

import (
    "github.com/go-advanced-admin/admin"
    "github.com/go-advanced-admin/web-echo"
    "github.com/go-advanced-admin/orm-gorm"
    "github.com/labstack/echo/v4"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func main() {
    // Initialize Echo
    e := echo.New()

    // Initialize GORM
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Initialize the admin panel
    webIntegrator := adminecho.NewIntegrator(e.Group("/admin"))
    ormIntegrator := admingorm.NewIntegrator(db)
    permissionFunc := func(req admin.PermissionRequest, ctx interface{}) (bool, error) {
        return true, nil // Implement your permission logic here
    }

    panel, err := admin.NewPanel(ormIntegrator, webIntegrator, permissionFunc, nil)
    if err != nil {
        panic(err)
    }

    // Register your models
    app, err := panel.RegisterApp("MainApp", "Main Application", nil)
    if err != nil {
        panic(err)
    }

    _, err = app.RegisterModel(&YourModel{}, nil)
    if err != nil {
        panic(err)
    }

    // Start the server
    e.Logger.Fatal(e.Start(":8080"))
}
```

For more detailed examples and configuration options, please refer to the 
[official documentation](https://goadmin.dev/quickstart).

## Contributing

Contributions are always welcome! If you're interested in contributing to the project, please take a look at our 
[Contributing Guidelines](CONTRIBUTING.md) for guidelines on how to get started. We appreciate your help in improving 
the library!

Special thank you to the current maintainers:

- [Yidi Sprei](https://github.com/YidiDev)
- [Coal Rock](https://github.com/coal-rock)

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
