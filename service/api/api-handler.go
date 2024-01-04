package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes

	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.POST("/users/:id/photo", rt.wrap(rt.uploadPhoto))
	rt.router.POST("/users/:id/newlogo", rt.wrap(rt.uploadLogo))
	rt.router.POST("/users/:id/comment/:photoId", rt.wrap(rt.commentPhoto))
	rt.router.PUT("/users/:id/like/:photoId", rt.wrap(rt.likePhoto))
	rt.router.PUT("/users/:id/username", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/users/:id/follower/:id2", rt.wrap(rt.followUser))
	rt.router.PUT("/users/:id/ban/:id2", rt.wrap(rt.banUser))
	rt.router.GET("/users/:id/profile", rt.wrap(rt.getUserProfile))
	rt.router.GET("/users/:id/following", rt.wrap(rt.getFollowingUsers))
	rt.router.GET("/users/:id", rt.wrap(rt.getUserByName))
	rt.router.GET("/username/:id", rt.wrap(rt.getUserName))
	rt.router.GET("/users/:id/stream", rt.wrap(rt.getMyStream))
	rt.router.GET("/users/:id/logo", rt.wrap(rt.getLogo))
	rt.router.GET("/images/:id", rt.wrap(rt.getImage))
	rt.router.DELETE("/users/:id/unfollower/:id2", rt.wrap(rt.unfollowUser))
	rt.router.DELETE("/users/:id/unban/:id2", rt.wrap(rt.unbanUser))
	rt.router.DELETE("/users/:id/dislike/:photoId", rt.wrap(rt.unlikePhoto))
	rt.router.DELETE("/users/:id/uncomment/:commentId", rt.wrap(rt.uncommentPhoto))
	rt.router.DELETE("/users/:id/photo/:photoId", rt.wrap(rt.deletePhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
