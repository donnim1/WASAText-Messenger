package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes

	rt.router.POST("/session", rt.doLogin)

	rt.router.PUT("/user/username", rt.setMyUserName)
	rt.router.PUT("/user/photo", rt.setMyPhoto)

	rt.router.GET("/users", rt.listUsers)
	rt.router.GET("/conversationsfor/:receiverId", rt.GetConversationByReceiver)
	rt.router.GET("/conversation/myconversations", rt.getMyConversations)
	rt.router.GET("/conversations/:conversationId", rt.getConversation)

	rt.router.POST("/messages", rt.sendMessage)
	rt.router.POST("/messages/:messageId/forward", rt.forwardMessage)

	rt.router.POST("/messages/:messageId/comments", rt.commentMessage)
	rt.router.DELETE("/messages/:messageId/uncomment", rt.uncommentMessage)
	rt.router.DELETE("/messages/:messageId", rt.deleteMessage)
	rt.router.POST("/messages/:messageId/status/:status", rt.updateMessageStatus)

	// Group endpoints
	rt.router.GET("/groups", rt.listGroups)
	rt.router.POST("/group", rt.createGroup)

	// Then register dynamic routes.
	rt.router.POST("/groups/:groupId/members", rt.addToGroup)
	rt.router.PUT("/groups/:groupId/name", rt.setGroupName)
	rt.router.PUT("/groups/:groupId/photo", rt.setGroupPhoto)
	rt.router.DELETE("/groups/:groupId/leave", rt.leaveGroup)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
