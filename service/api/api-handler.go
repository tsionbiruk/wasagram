package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.POST("/users/:user/username", rt.wrap(rt.postUsername))
	rt.router.POST("/users/:user/photos", rt.wrap(rt.postPhoto))
	rt.router.DELETE("/users/:user/photos/:photoid", rt.wrap(rt.deletePhoto))
	rt.router.GET("/photos/:photoid", rt.wrap(rt.getPhoto))
	rt.router.PUT("/users/:user/likes/:photoid", rt.wrap(rt.putLike))
	rt.router.DELETE("/users/:user/likes/:photoid", rt.wrap(rt.deleteLike))
	rt.router.POST("/users/:user/comments/:photoid", rt.wrap(rt.postComments))
	rt.router.DELETE("/users/:user/photos/:photoid/comments/:commentid", rt.wrap(rt.deleteComments))

	rt.router.PUT("/users/:user/followed/:targetuser", rt.wrap(rt.putFollowed))
	rt.router.DELETE("/users/:user/followed/:targetuser", rt.wrap(rt.deleteFollowed))
	rt.router.GET("/users/:user/followed", rt.wrap(rt.getFollowed))
	rt.router.PUT("/users/:user/banned/:targetuser", rt.wrap(rt.putBanned))
	rt.router.DELETE("/users/:user/banned/:targetuser", rt.wrap(rt.deleteBanned))
	rt.router.GET("/users/:user/banned", rt.wrap(rt.getBanned))
	rt.router.GET("/users/:user/profile", rt.wrap(rt.getProfile))
	rt.router.GET("/users/:user/stream", rt.wrap(rt.getStream))
	rt.router.GET("/users", rt.wrap(rt.getUsers))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
