package generator

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"testing"
)

func TestServer_Serve(t *testing.T) {
	type fields struct {
		app       *iris.Application
		Bindata   Binary
		Extension string
		FsRoute   string
		HttpPort  int
		PublicDir string
		Reload    bool
		Routes    []Routes
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "test",
			fields: fields{
				Bindata:   Binary{},
				Extension: ".html",
				FsRoute:   "/",
				PublicDir: "./public",
				Reload:    true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				app:       tt.fields.app,
				Bindata:   tt.fields.Bindata,
				Extension: tt.fields.Extension,
				FsRoute:   tt.fields.FsRoute,
				HttpPort:  tt.fields.HttpPort,
				PublicDir: tt.fields.PublicDir,
				Reload:    tt.fields.Reload,
				Routes:    tt.fields.Routes,
				Testing:   true,
			}
			s.App()

			// register routes
			s.Register(Routes{
				Fn: func(ctx Context) {
					ctx.ViewData("code", "LoremIpsum")
					_ = ctx.View("index.html")
				},
				Method: "GET",
				Path:   "/",
			})
			s.Register(Routes{
				Fn: func(ctx Context) {
					ctx.JSON(Map{"status": true, "content": "LoremIpsum"})
				},
				Method: "GET",
				Path:   "/api",
			})
			s.Serve(false)

			// test routes
			app := httptest.New(t, s.app)
			app.GET("/api").Expect().Status(200).
				JSON().
				Object().Keys().Contains("status", "content")
		})
	}
}
