package middlewares

import (
	"github.com/gin-gonic/gin"
)

// CheckPermission 检查权限
func CheckPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		////对url进行处理去掉请求
		//url := c.Request.URL.Path
		//for _, param := range c.Params {
		//	re, _ := regexp.Compile(param.Value)
		//	url = re.ReplaceAllString(url, ":"+param.Key)
		//}
		//method := c.Request.Method
		//// 得到当前的权限
		//permission, err := rbac.NewPermissionService().Get(url, method)
		//if err == nil {
		//	// 得到当前用户信息
		//	u, _ := c.Get("admin")
		//	staff := u.(*models.Staff)
		//	userID := staff.ID
		//	// 得到用户信息
		//	user, _ := staff2.NewService().Get(userID)
		//	if !user.IsAdmin && !staff2.NewService().HasPermission(userID, permission.ID) {
		//		response.New().APIError(c, http.StatusUnauthorized, "permission deny")
		//		c.Abort()
		//		return
		//	}
		//}
		c.Next()
	}
}

// ShowMenu 是否显示菜单
func ShowMenu(userID string, path interface{}) bool {
	return true
	//p, ok := path.(string)
	//if !ok {
	//	p = ""
	//}
	//// 1. 得到权限信息
	//permission, err := rbac.NewPermissionService().Get(p, "GET")
	//if err != nil {
	//	return false
	//}
	//userInfo, err := staff2.NewService().Get(userID)
	//if err != nil {
	//	return false
	//}
	////  2. 判断用户是否有权限
	//if userInfo.IsAdmin || staff2.NewService().HasPermission(userID, permission.ID) {
	//	return true
	//}
	//return falsee
}

// ShowStaffMenu ShowStaffMenu
// func ShowStaffMenu(userID string, path string) bool {
// 	// 1. 得到权限信息
// 	permission, err := services.GetPermissionByPathAndMethod(path, "GET")
// 	if err != nil {
// 		return false
// 	}
// 	//  2. 判断用户是否有权限
// 	if services.HasPermission(userID, permission.(*rbac.Permission).ID) {
// 		return true
// 	}
// 	return false
// }
