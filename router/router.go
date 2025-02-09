package router

import (
	"net/http"

	"github.com/yoshinori0811/chat_app_backend/controller"
	"github.com/yoshinori0811/chat_app_backend/middleware"
)

func NewRouter(m middleware.MiddlewareInterface, uc controller.UserControllerInterface, fc controller.FriendControllerInterface, rc controller.RoomControllerInterface) {
	http.HandleFunc("/signup", m.CorsMiddleware(&middleware.MethodHandler{
		Post: uc.SignUp,
	}))
	http.HandleFunc("/login", m.CorsMiddleware(&middleware.MethodHandler{
		Post: uc.Login,
	}))
	http.HandleFunc("/logout", m.CorsMiddleware(&middleware.MethodHandler{
		Post: uc.Logout,
	}))
	http.HandleFunc("/user", m.CorsMiddleware(m.AuthMiddleware(&middleware.MethodHandler{
		Get: uc.GetUser,
	})))
	http.HandleFunc("/users", m.CorsMiddleware(m.AuthMiddleware(&middleware.MethodHandler{
		Get: uc.SearchUsers,
	})))
	http.HandleFunc("/dmlist", m.CorsMiddleware(m.AuthMiddleware(&middleware.MethodHandler{
		Get: fc.GetDmList,
	})))
	http.HandleFunc("/friends", m.CorsMiddleware(m.AuthMiddleware(&middleware.MethodHandler{
		Get: fc.GetFriends,
	})))
	http.HandleFunc("/friends/requests", m.CorsMiddleware(m.AuthMiddleware(&middleware.MethodHandler{
		Get:  fc.GetFriendRequestList,
		Post: fc.SendFriendRequest,
	})))
	http.HandleFunc("/friends/requests/accept", m.CorsMiddleware(m.AuthMiddleware(&middleware.MethodHandler{
		Put: fc.AcceptFriendRequest,
	})))
	http.HandleFunc("/friends/requests/reject", m.CorsMiddleware(m.AuthMiddleware(&middleware.MethodHandler{
		Put: fc.RejectFriendRequest,
	})))

	http.HandleFunc("/rooms/create", m.CorsMiddleware(m.AuthMiddleware(&middleware.MethodHandler{
		Post: rc.CreateRoom,
	})))

	http.HandleFunc("/rooms/", m.CorsMiddleware(m.AuthMiddleware(&middleware.MethodHandler{
		Get: rc.GetRooms,
	})))

	http.HandleFunc("/rooms/{roomUUID}/", m.CorsMiddleware(m.AuthMiddleware(&middleware.MethodHandler{
		Get:  rc.GetRoomChat,
		Post: rc.CreateMessage,
	})))

	http.HandleFunc("/rooms/{roomUUID}/messages/{messageUUID}/", m.CorsMiddleware(m.AuthMiddleware(&middleware.MethodHandler{
		Patch:  rc.UpdateMessage,
		Delete: rc.DeleteMessage,
	})))
	http.HandleFunc("/rooms/{roomUUID}/invite", m.CorsMiddleware(m.AuthMiddleware(&middleware.MethodHandler{
		Post: rc.InviteRoom,
	})))
	http.HandleFunc("/rooms/{roomUUID}/delete", m.CorsMiddleware(m.AuthMiddleware(&middleware.MethodHandler{
		Delete: rc.DeleteRoom,
	})))
	http.HandleFunc("/rooms/{roomUUID}/leave", m.CorsMiddleware(m.AuthMiddleware(&middleware.MethodHandler{
		Delete: rc.LeaveRoom,
	})))
}
