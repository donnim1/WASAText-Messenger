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

	rt.router.GET("/users", rt.listUsers)

	rt.router.GET("/conversation/myconversations", rt.getMyConversations)
	rt.router.GET("/conversations/:conversationId", rt.getConversation)

	rt.router.POST("/messages", rt.sendMessage)
	rt.router.POST("/messages/:messageId/forward", rt.forwardMessage)

	rt.router.POST("/messages/:messageId/comments", rt.commentMessage)
	rt.router.DELETE("/messages/:messageId/uncomment", rt.uncommentMessage)
	rt.router.DELETE("/messages/:messageId/delete", rt.deleteMessage)

	// Group endpoints

	// For example, in your API handler registration:
	rt.router.POST("/groups/create", rt.createGroup)
	// Register the GET /groups endpoint.
	rt.router.GET("/groups", rt.listGroups)

	rt.router.POST("/groups/{groupId}/members", rt.addToGroup)
	rt.router.DELETE("/groups/{groupId}/leave", rt.leaveGroup)
	rt.router.PUT("/groups/{groupId}/name", rt.setGroupName)
	rt.router.PUT("/groups/{groupId}/photo", rt.setGroupPhoto)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
