package auth

// // @Id 001
// // @Tags 用户权限
// // @Summary 用户注册
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param signup_info body SignupRequest true "登录类型"
// // @Success 200 object response.SuccessRes{data=int} 注册成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /signup [POST]
// func Signup(c *gin.Context) {
// 	var signupInfo SignupRequest
// 	err := c.ShouldBindJSON(&signupInfo)
// 	if err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	authID, err := authService.CreateAuth(signupInfo)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, authID)
// }

// // @Summary 登录
// // @Id 002
// // @Tags 用户权限
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param signin_info body SigninRequest true "登录类型"
// // @Success 200 object response.SuccessRes{data=SigninResponse} 登录成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Failure 401 object response.ErrorRes 登录失败
// // @Router /signin [POST]
// func Signin(c *gin.Context) {
// 	var signinInfo SigninRequest
// 	var userInfo *UserResponse
// 	err := c.ShouldBindJSON(&signinInfo)
// 	if err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	userInfo, err = authService.VerifyCredential(signinInfo)
// 	if err != nil {
// 		response.ResponseUnauthorized(c, "AuthError", err)
// 		return
// 	}
// 	claims := service.CustomClaims{
// 		UserID:           userInfo.ID,
// 		UserType:         userInfo.Type,
// 		Username:         userInfo.Name,
// 		RoleID:           userInfo.RoleID,
// 		OrganizationID:   userInfo.OrganizationID,
// 		OrganizationName: userInfo.OrganizationName,
// 		PositionID:       userInfo.PositionID,
// 		StandardClaims: jwt.StandardClaims{
// 			NotBefore: time.Now().Unix() - 1000,
// 			ExpiresAt: time.Now().Unix() + 72000,
// 			Issuer:    "wms",
// 		},
// 	}
// 	jwtServices := service.JWTAuthService()
// 	generatedToken := jwtServices.GenerateToken(claims)
// 	var res SigninResponse
// 	res.Token = generatedToken
// 	res.User = *userInfo
// 	response.Response(c, res)
// }

// // @Summary 角色列表
// // @Id 18
// // @Tags 角色管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param page_id query int true "页码"
// // @Param page_size query int true "每页行数（5/10/15/20）"
// // @Param name query string false "角色名称"
// // @Success 200 object response.ListRes{data=[]Role} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /roles [GET]
// func GetRoleList(c *gin.Context) {
// 	var filter RoleFilter
// 	err := c.ShouldBindQuery(&filter)
// 	if err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	count, list, err := authService.GetRoleList(filter)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
// }

// // @Summary 新建角色
// // @Id 19
// // @Tags 角色管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param role_info body RoleNew true "角色信息"
// // @Success 200 object response.SuccessRes{data=Role} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /roles [POST]
// func NewRole(c *gin.Context) {
// 	var role RoleNew
// 	if err := c.ShouldBindJSON(&role); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	role.User = claims.Username
// 	authService := NewAuthService()
// 	new, err := authService.NewRole(role)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, new)
// }

// // @Summary 根据ID获取角色
// // @Id 20
// // @Tags 角色管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "角色ID"
// // @Success 200 object response.SuccessRes{data=Role} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /roles/:id [GET]
// func GetRoleByID(c *gin.Context) {
// 	var uri RoleID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	role, err := authService.GetRoleByID(uri.ID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, role)

// }

// // @Summary 根据ID更新角色
// // @Id 21
// // @Tags 角色管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "角色ID"
// // @Param role_info body RoleNew true "角色信息"
// // @Success 200 object response.SuccessRes{data=Role} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /roles/:id [PUT]
// func UpdateRole(c *gin.Context) {
// 	var uri RoleID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	var role RoleNew
// 	if err := c.ShouldBindJSON(&role); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	role.User = claims.Username
// 	authService := NewAuthService()
// 	new, err := authService.UpdateRole(uri.ID, role)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, new)
// }

// // @Summary 根据ID更新用户
// // @Id 23
// // @Tags 用户管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "用户ID"
// // @Param menu_info body UserUpdate true "用户信息"
// // @Success 200 object response.SuccessRes{data=User} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /users/:id [PUT]
// func UpdateUser(c *gin.Context) {
// 	var uri UserID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	var user UserUpdate
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	user.User = claims.Username
// 	authService := NewAuthService()
// 	new, err := authService.UpdateUser(uri.ID, user, claims.UserID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, new)
// }

// // @Summary 用户列表
// // @Id 32
// // @Tags 用户管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param page_id query int true "页码"
// // @Param page_size query int true "每页行数（5/10/15/20）"
// // @Param name query string false "用户名称"
// // @Param type query string false "用户类型wx/admin"
// // @Param organization_id query int false "用户组织"
// // @Success 200 object response.ListRes{data=[]UserResponse} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /users [GET]
// func GetUserList(c *gin.Context) {
// 	var filter UserFilter
// 	err := c.ShouldBindQuery(&filter)
// 	if err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	organizationID := claims.OrganizationID
// 	authService := NewAuthService()
// 	count, list, err := authService.GetUserList(filter, organizationID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
// }

// // @Summary 根据ID获取用户
// // @Id 33
// // @Tags 用户管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "用户ID"
// // @Success 200 object response.SuccessRes{data=User} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /users/:id [GET]
// func GetUserByID(c *gin.Context) {
// 	var uri UserID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	organizationID := claims.OrganizationID
// 	authService := NewAuthService()
// 	user, err := authService.GetUserByID(uri.ID, organizationID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, user)

// }

// // @Summary API列表
// // @Id 36
// // @Tags API管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param page_id query int true "页码"
// // @Param page_size query int true "每页行数"
// // @Param name query string false "API名称"
// // @Param route query string false "API路由"
// // @Success 200 object response.ListRes{data=[]API} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /apis [GET]
// func GetAPIList(c *gin.Context) {
// 	var filter APIFilter
// 	err := c.ShouldBindQuery(&filter)
// 	if err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	count, list, err := authService.GetAPIList(filter)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
// }

// // @Summary 新建API
// // @Id 37
// // @Tags API管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param api_info body APINew true "API信息"
// // @Success 200 object response.SuccessRes{data=API} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /apis [POST]
// func NewAPI(c *gin.Context) {
// 	var api APINew
// 	if err := c.ShouldBindJSON(&api); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	api.User = claims.Username
// 	authService := NewAuthService()
// 	new, err := authService.NewAPI(api)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, new)
// }

// // @Summary 根据ID获取API
// // @Id 38
// // @Tags API管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "API ID"
// // @Success 200 object response.SuccessRes{data=API} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /apis/:id [GET]
// func GetAPIByID(c *gin.Context) {
// 	var uri APIID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	api, err := authService.GetAPIByID(uri.ID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, api)

// }

// // @Summary 根据ID更新API
// // @Id 39
// // @Tags API管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "API ID"
// // @Param api_info body APINew true "API信息"
// // @Success 200 object response.SuccessRes{data=API} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /apis/:id [PUT]
// func UpdateAPI(c *gin.Context) {
// 	var uri APIID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	var api APINew
// 	if err := c.ShouldBindJSON(&api); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	api.User = claims.Username
// 	authService := NewAuthService()
// 	new, err := authService.UpdateAPI(uri.ID, api)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, new)
// }

// // @Summary 菜单列表
// // @Id 40
// // @Tags 菜单管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param page_id query int true "页码"
// // @Param page_size query int true "每页行数（5/10/15/20）"
// // @Param name query string false "菜单名称"
// // @Param only_top query bool false "只显示顶级菜单"
// // @Success 200 object response.ListRes{data=[]Menu} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /menus [GET]
// func GetMenuList(c *gin.Context) {
// 	var filter MenuFilter
// 	err := c.ShouldBindQuery(&filter)
// 	if err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	count, list, err := authService.GetMenuList(filter)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
// }

// // @Summary 新建菜单
// // @Id 41
// // @Tags 菜单管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param menu_info body MenuNew true "菜单信息"
// // @Success 200 object response.SuccessRes{data=Menu} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /menus [POST]
// func NewMenu(c *gin.Context) {
// 	var menu MenuNew
// 	if err := c.ShouldBindJSON(&menu); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	menu.User = claims.Username
// 	authService := NewAuthService()
// 	new, err := authService.NewMenu(menu)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, new)
// }

// // @Summary 根据ID获取菜单
// // @Id 42
// // @Tags 菜单管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "菜单ID"
// // @Success 200 object response.SuccessRes{data=Menu} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /menus/:id [GET]
// func GetMenuByID(c *gin.Context) {
// 	var uri MenuID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	menu, err := authService.GetMenuByID(uri.ID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, menu)

// }

// // @Summary 根据ID更新菜单
// // @Id 43
// // @Tags 菜单管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "菜单ID"
// // @Param menu_info body MenuNew true "菜单信息"
// // @Success 200 object response.SuccessRes{data=Menu} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /menus/:id [PUT]
// func UpdateMenu(c *gin.Context) {
// 	var uri MenuID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	var menu MenuUpdate
// 	if err := c.ShouldBindJSON(&menu); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	menu.User = claims.Username
// 	authService := NewAuthService()
// 	new, err := authService.UpdateMenu(uri.ID, menu)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, new)
// }

// // @Summary 根据ID更新菜单
// // @Id 50
// // @Tags 菜单管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "菜单ID"
// // @Param menu_info body MenuNew true "菜单信息"
// // @Success 200 object response.SuccessRes{data=string} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /menus/:id [DELETE]
// func DeleteMenu(c *gin.Context) {
// 	var uri MenuID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	authService := NewAuthService()
// 	err := authService.DeleteMenu(uri.ID, claims.Username)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, "OK")
// }

// // @Summary 根据角色ID获取菜单权限
// // @Id 44
// // @Tags 权限管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "角色ID"
// // @Success 200 object response.SuccessRes{data=[]int64} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /rolemenus/:id [GET]
// func GetRoleMenu(c *gin.Context) {
// 	var uri RoleID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	menu, err := authService.GetRoleMenuByID(uri.ID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, menu)

// }

// // @Summary 根据角色ID更新菜单权限
// // @Id 45
// // @Tags 权限管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "角色ID"
// // @Param menu_info body RoleMenu true "菜单信息"
// // @Success 200 object response.SuccessRes{data=string} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /rolemenus/:id [POST]
// func NewRoleMenu(c *gin.Context) {
// 	var uri RoleID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	var menu RoleMenuNew
// 	if err := c.ShouldBindJSON(&menu); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	menu.User = claims.Username
// 	authService := NewAuthService()
// 	err := authService.NewRoleMenu(uri.ID, menu)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, "OK")
// }

// // @Summary 根据菜单ID获取API权限
// // @Id 46
// // @Tags 权限管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "菜单ID"
// // @Success 200 object response.SuccessRes{data=[]int64} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /menuapis/:id [GET]
// func GetMenuApi(c *gin.Context) {
// 	var uri MenuID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	menu, err := authService.GetMenuAPIByID(uri.ID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, menu)

// }

// // @Summary 根据菜单ID更新API权限
// // @Id 47
// // @Tags 权限管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "菜单ID"
// // @Param menu_info body MenuAPINew true "菜单信息"
// // @Success 200 object response.SuccessRes{data=string} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /menuapis/:id [POST]
// func NewMenuApi(c *gin.Context) {
// 	var uri MenuID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	var menu MenuAPINew
// 	if err := c.ShouldBindJSON(&menu); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	menu.User = claims.Username
// 	authService := NewAuthService()
// 	err := authService.NewMenuAPI(uri.ID, menu)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, "OK")
// }

// // @Summary 获取当前用户的前端路由
// // @Id 48
// // @Tags 权限管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Success 200 object response.SuccessRes{data=interface{}} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /mymenu [GET]
// func GetMyMenu(c *gin.Context) {
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	role_id := claims.RoleID
// 	authService := NewAuthService()
// 	new, err := authService.GetMyMenu(role_id)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	res := make(map[int64]*MyMenuDetail)
// 	for i := 0; i < len(new); i++ {
// 		if new[i].ParentID == -1 {
// 			var m MyMenuDetail
// 			m.Action = new[i].Action
// 			m.Component = new[i].Component
// 			m.Name = new[i].Name
// 			m.Title = new[i].Title
// 			m.Path = new[i].Path
// 			m.IsHidden = new[i].IsHidden
// 			m.Status = new[i].Status
// 			res[new[i].ID] = &m
// 		} else {
// 			var m MyMenuDetail
// 			m.Action = new[i].Action
// 			m.Component = new[i].Component
// 			m.Name = new[i].Name
// 			m.Title = new[i].Title
// 			m.Path = new[i].Path
// 			m.IsHidden = new[i].IsHidden
// 			m.Status = new[i].Status
// 			res[new[i].ParentID].Items = append(res[new[i].ParentID].Items, m)
// 		}
// 	}
// 	response.Response(c, res)
// }

// // @Summary 根据ID删除角色
// // @Id 52
// // @Tags 菜单管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "菜单ID"
// // @Success 200 object response.SuccessRes{data=string} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /roles/:id [DELETE]
// func DeleteRole(c *gin.Context) {
// 	var uri MenuID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	authService := NewAuthService()
// 	err := authService.DeleteRole(uri.ID, claims.Username)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, "OK")
// }

// // @Summary 用户列表
// // @Id 81
// // @Tags 小程序接口
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param page_id query int true "页码"
// // @Param page_size query int true "每页行数（5/10/15/20）"
// // @Param name query string false "用户名称"
// // @Success 200 object response.ListRes{data=[]Role} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /wx/users [GET]
// func WxGetUserList(c *gin.Context) {
// 	GetUserList(c)
// }

// // @Summary 根据ID更新用户
// // @Id 88
// // @Tags 小程序接口
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "用户ID"
// // @Param menu_info body UserUpdate true "用户信息"
// // @Success 200 object response.SuccessRes{data=User} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /wx/users/:id [PUT]
// func WxUpdateUser(c *gin.Context) {
// 	UpdateUser(c)
// }

// // @Summary 更新密码
// // @Id 102
// // @Tags 用户管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param menu_info body UserUpdate true "用户信息"
// // @Success 200 object response.SuccessRes{data=string} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /password [POST]
// func UpdatePassword(c *gin.Context) {
// 	var info PasswordUpdate
// 	if err := c.ShouldBindJSON(&info); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	info.User = claims.Username
// 	info.UserID = claims.UserID
// 	authService := NewAuthService()
// 	err := authService.UpdatePassword(info)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, "ok")
// }
