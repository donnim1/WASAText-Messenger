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

	rt.router.GET("/conversation/myconversations", rt.getMyConversations)
	rt.router.GET("/conversations/:conversationId", rt.getConversation)

	rt.router.POST("/message", rt.sendMessage)
    rt.router.POST("/message/:messageId/forward", rt.forwardMessage)
    rt.router.POST("/message/:messageId/comment", rt.commentMessage)
    rt.router.DELETE("/message/:messageId/comment", rt.uncommentMessage)
    rt.router.DELETE("/message/:messageId", rt.deleteMessage)



	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
