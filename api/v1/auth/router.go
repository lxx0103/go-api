package auth

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.POST("/signin", Signin)
	// g.POST("/signup", Signup)
}

func AuthRouter(g *gin.RouterGroup) {
	g.GET("/roles", GetRoleList)
	g.POST("/roles", NewRole)
	g.PUT("/roles/:id", UpdateRole)
	g.GET("/roles/:id", GetRoleByID)
	g.DELETE("/roles/:id", DeleteRole)

	// g.PUT("/users/:id", UpdateUser)
	// g.GET("/users", GetUserList)
	// g.GET("/users/:id", GetUserByID)
	// g.POST("/password", UpdatePassword)

	// g.GET("/menus", GetMenuList)

	// g.GET("/rolemenus/:id", GetRoleMenu)
	// g.POST("/rolemenus/:id", NewRoleMenu)

	// g.GET("/mymenu", GetMyMenu)

}
