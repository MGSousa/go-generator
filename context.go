package generator

import "github.com/kataras/iris/v12"

type (
	// Context is the middle-man server's "object" for the clients.
	//
	// A New context is being acquired from a sync.Pool on each connection.
	// The Context is the most important thing on the iris's http flow.
	//
	// Developers send responses to the client's request through a Context.
	// Developers get request information from the client's request by a Context.
	Context = iris.Context

	// Map An alias for iris.Map
	// which internally is an alias of map[string]interface{}.
	Map = iris.Map
)
