package api

import (
	"net/http"
)

// redirect all the requests here. ALL off them
// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	rt.router.POST("/session", rt.wrap(rt.DoLogin))
	rt.router.POST("/users/:user/photos", rt.wrap(rt.UploadPhoto))
	rt.router.GET("/photos/:photoid", rt.wrap(rt.getPhoto))
	rt.router.DELETE("/users/:user/photos/:photoid", rt.wrap(rt.DeletePost))
	rt.router.PUT("/users/:user/photos/:photoid/like", rt.wrap(rt.Photolike))
	rt.router.DELETE("/users/:user/photos/:photoid/unlike", rt.wrap(rt.Photounlike))
	rt.router.POST("/users/:user/photos/:photoid/comment", rt.wrap(rt.PostComments))
	rt.router.DELETE("/users/:user/photos/:photoid/comment/:commentid", rt.wrap(rt.DeleteComments))
	rt.router.PUT("/users/:user/follows", rt.wrap(rt.Followuser))
	rt.router.DELETE("/users/:user/follows/:targetuser", rt.wrap(rt.Unfollowuser))
	rt.router.PUT("/users/:user/bans/:targetuser", rt.wrap(rt.Ban))
	rt.router.DELETE("/users/:user/bans/:targetuser", rt.wrap(rt.Unban))
	rt.router.GET("/users/:user/bans/:targetuser", rt.wrap(rt.GetBanned))
	rt.router.GET("/users/:user/follows", rt.wrap(rt.GetFollowers))
	rt.router.GET("/users/:user/follows/:targetuser", rt.wrap(rt.Getfollowing))
	rt.router.GET("/users", rt.wrap(rt.GetUsers))
	rt.router.GET("/users/:user/user/:targetuser/profile", rt.wrap(rt.GetProfile))

	rt.router.GET("/users/:user/stream/:targetuser", rt.wrap(rt.GetStream))

	rt.router.POST("/users/:user/username", rt.wrap(rt.Rename))

	return rt.router
}
