package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/LCA/internal/adapter/validator"
	"github.com/wang900115/LCA/internal/application/usecase"
	"github.com/wang900115/LCA/pkg/common"
	iresponse "github.com/wang900115/LCA/pkg/common/response"
	"github.com/wang900115/LCA/pkg/domain"
)

type UserController struct {
	user usecase.UserUsecase
	resp iresponse.IResponse
}

func NewUserController(user *usecase.UserUsecase, resp iresponse.IResponse) *UserController {
	return &UserController{user: *user, resp: resp}
}

func (u *UserController) Login(c *gin.Context) {
	var req validator.LoginAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}

	accessToken, refreshToken, user, err := u.user.Login(c, req.Username, req.Password)
	if err != nil {
		u.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	u.resp.SuccessWithData(c, common.LOGIN_SUCESS, map[string]interface{}{
		"aaccessToken": accessToken,
		"refreshToken": refreshToken,
		"user":         user})
}

func (u *UserController) Logout(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := u.user.Logout(c, userID); err != nil {
		u.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	u.resp.Success(c, common.LOGOUT_SUCCESS)
}

func (u *UserController) Register(c *gin.Context) {
	var req validator.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}

	user := domain.User{
		Username: req.Username,
		Password: &req.Password,

		FirstEmail:  req.FirstEmail,
		SecondEmail: req.SecondEmail,
		Phone:       req.Phone,
		NickName:    req.NickName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Birth:       req.Birth,
		Country:     req.Country,
		City:        req.City,
	}
	created, err := u.user.Register(c, user)
	if err != nil {
		u.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	u.resp.SuccessWithData(c, common.CREATE_SUCCESS, map[string]interface{}{
		"user": created,
	})
}

func (u *UserController) Query(c *gin.Context) {
	users, err := u.user.QueryUser(c)
	if err != nil {
		u.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	u.resp.SuccessWithData(c, common.QUERY_SUCCESS, map[string]interface{}{
		"users": users,
	})
}

func (u *UserController) Update(c *gin.Context) {
	var req validator.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}

	userID := c.GetUint("user_id")
	user := domain.User{
		ID:       userID,
		Username: req.Username,
		Password: &req.Password,

		FirstEmail:  req.FirstEmail,
		SecondEmail: req.SecondEmail,
		Phone:       req.Phone,
		NickName:    req.NickName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Birth:       req.Birth,
		Country:     req.Country,
		City:        req.City,
	}

	updated, err := u.user.UpdateUser(c, user)
	if err != nil {
		u.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	u.resp.SuccessWithData(c, common.UPDATE_SUCCESS, map[string]interface{}{
		"user": updated,
	})
}

func (u *UserController) Delete(c *gin.Context) {
	var req validator.DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}
	userID := c.GetUint("user_id")
	if err := u.user.DeleteUser(c, userID); err != nil {
		u.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	u.resp.Success(c, common.DELETE_SUCCESS)
}
