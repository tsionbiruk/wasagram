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
	rt.router.DELETE("/users/:user/photos/:photoid", rt.wrap(rt.DeletePost))
	rt.router.PUT("/users/:user/likes/:photoid", rt.wrap(rt.Photolike))
	rt.router.DELETE("/users/:user/unlikes/:photoid", rt.wrap(rt.Photounlike))
	rt.router.POST("/users/:user/comments/:photoid", rt.wrap(rt.PostComments))
	rt.router.DELETE("/users/:user/comments/:photoid/:commentid", rt.wrap(rt.DeleteComments))
	rt.router.PUT("/users/:user/follow/:targetuser", rt.wrap(rt.Followuser))
	rt.router.DELETE("/users/:user/unfollow/:targetuser", rt.wrap(rt.Unfollowuser))
	rt.router.PUT("/users/:user/ban/:targetuser", rt.wrap(rt.Ban))
	rt.router.DELETE("/users/:user/unban/:targetuser", rt.wrap(rt.Unban))
	rt.router.GET("/users/:user/banned/:targetuser", rt.wrap(rt.GetBanned))
	rt.router.GET("/users/:user/followers/:targetuser", rt.wrap(rt.GetFollowers))
	rt.router.GET("/users/:user/following/:targetuser", rt.wrap(rt.Getfollowing))

	//rt.router.GET("/photos/:photoid", rt.wrap(rt.getPhoto))
	rt.router.POST("/users/:user/username", rt.wrap(rt.Rename))

	rt.router.GET("/users/:user/profile/:targetuser", rt.wrap(rt.GetProfile))
	rt.router.GET("/users/:user/stream", rt.wrap(rt.GetStream))
	rt.router.GET("/users", rt.wrap(rt.GetUsers))

	return rt.router
}
