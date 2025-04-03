package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes

	rt.router.POST("/session", rt.wrap(rt.doLogin))

	rt.router.PUT("/user/username", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/user/photo", rt.wrap(rt.setMyPhoto))

	rt.router.GET("/users", rt.wrap(rt.listUsers))
	rt.router.GET("/conversationsfor/:receiverId", rt.wrap(rt.GetConversationByReceiver))
	rt.router.GET("/conversation/myconversations", rt.wrap(rt.getMyConversations))
	rt.router.GET("/conversations/:conversationId", rt.wrap(rt.getConversation))

	rt.router.POST("/messages", rt.wrap(rt.sendMessage))
	rt.router.POST("/messages/:messageId/forward", rt.wrap(rt.forwardMessage))

	rt.router.POST("/messages/:messageId/comments", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/messages/:messageId/uncomment", rt.wrap(rt.uncommentMessage))
	rt.router.DELETE("/messages/:messageId", rt.wrap(rt.deleteMessage))
	rt.router.POST("/messages/:messageId/status/:status", rt.wrap(rt.updateMessageStatus))

	// Group endpoints
	rt.router.GET("/groups", rt.wrap(rt.listGroups))
	rt.router.POST("/group", rt.wrap(rt.createGroup))

	// Then register dynamic routes.
	rt.router.POST("/groups/:groupId/members", rt.wrap(rt.addToGroup))
	rt.router.PUT("/groups/:groupId/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/groups/:groupId/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.DELETE("/groups/:groupId/leave", rt.wrap(rt.leaveGroup))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
