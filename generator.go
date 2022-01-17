package generator

import (
	"errors"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
)

type (
	Server struct {
		app     *iris.Application
		event   *socketio.Server
		channel string

		// specify generated assets from go-bindata
		Bindata Binary

		// ".html"
		Extension string

		// "/"
		FsRoute string

		// HttpPort by default 5000
		HttpPort int

		// specify templates directory ./public
		PublicDir string

		// auto-reload template files
		Reload bool

		Testing bool

		// specify Server routes
		Routes []Routes

		// enable Ws events
		Ws bool
	}
)

var err error

// App init application
// if defined also init WS events
func (s *Server) App() {
	s.app = iris.New()

	if s.Ws {
		s.event, err = socketio.NewServer(nil)
		if err != nil {
			log.Fatal(err)
		}
		s.rtHandle()
		go s.event.Serve()
	}
}

// Serve application
// compatible only with Handlebars
func (s *Server) Serve(productionMode bool) {
	var (
		engine  = iris.Handlebars(s.PublicDir, s.Extension)
		options iris.DirOptions
		devMode = true
	)
	if productionMode {
		devMode = false
		// to use Assets you need to compress template files using go-bindata
		options = iris.DirOptions{
			Asset:      s.Bindata.Asset,
			AssetInfo:  s.Bindata.AssetInfo,
			AssetNames: s.Bindata.AssetNames,
			Gzip:       s.Bindata.Gzip,
		}
		// embed assets
		engine.Binary(s.Bindata.Asset, s.Bindata.AssetNames)
	}

	// Reload template files on changes
	if s.Reload {
		engine.Reload(devMode)
	}

	// register the view engine to load the templates
	// and create default handler for Assets
	s.app.RegisterView(engine)
	if s.FsRoute == "" {
		s.FsRoute = "/"
	}
	s.app.HandleDir(s.FsRoute, s.PublicDir, options)

	// register multiple routes
	for i := range s.Routes {
		s.app.Handle(s.Routes[i].Method, s.Routes[i].Path, s.Routes[i].Fn)
	}

	if s.Testing {
		return
	}
	if s.HttpPort == 0 {
		s.HttpPort = 5000
	}

	if s.Ws {
		defer s.event.Close()

		// registers socket io routes
		s.app.HandleMany("GET POST", "/socket.io/{any:path}", iris.FromStd(s.event))
		_ = s.app.Listen(fmt.Sprintf(":%d", s.HttpPort), iris.WithoutPathCorrection)
	} else {
		_ = s.app.Listen(fmt.Sprintf(":%d", s.HttpPort))
	}
}

// WsEvents returns event object to be called externally
func (s *Server) WsEvents(channel string) (*socketio.Server, error) {
	if s.Ws {
		if channel == "" {
			return nil, errors.New("error: WebSocket channel must not be empty")
		}
		s.channel = channel
		return s.event, nil
	}
	return nil, errors.New("error: Ws is disabled! To use WebSocket events, please enable it")
}

// rtHandle handles real-time connections
func (s *Server) rtHandle() {
	s.event.OnConnect("/", func(ws socketio.Conn) error {
		ws.SetContext("")
		log.Infoln("connected:", ws.RemoteAddr())
		ws.Join(s.channel)
		return nil
	})

	s.event.OnError("/", func(ws socketio.Conn, err error) {
		log.Infoln(err, ws.RemoteAddr())
	})
	s.event.OnDisconnect("/", func(ws socketio.Conn, reason string) {
		log.Println("closed", reason)
		ws.Leave(s.channel)
	})
}
