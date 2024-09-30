module github.com/go-advanced-admin/admin/example

go 1.23.1

replace github.com/go-advanced-admin/admin => ./../

replace github.com/go-advanced-admin/web-echo => ./../web-integrations/echo/

require (
	github.com/go-advanced-admin/admin v0.0.0
	github.com/go-advanced-admin/web-echo v0.0.0-00010101000000-000000000000
)

require (
	github.com/labstack/echo/v4 v4.12.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)
