package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/session", rt.postSession)

	rt.router.PUT("/user/username", rt.setMyUserName)
	rt.router.PUT("/user/photo", rt.setMyPhoto)

	//rt.router.GET()

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
