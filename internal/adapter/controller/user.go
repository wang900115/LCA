package controller

import (
	iresponse "github.com/wang900115/LCA/internal/adapter/controller/response"
	"github.com/wang900115/LCA/internal/adapter/validator"
	"github.com/wang900115/LCA/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	response iresponse.IResponse
	user     usecase.UserUsecase
	channel  usecase.ChannelUsecase
	message  usecase.MessageUsecase
}

func NewUserController(response iresponse.IResponse, user *usecase.UserUsecase, channel *usecase.ChannelUsecase, message *usecase.MessageUsecase) *UserController {
	return &UserController{response: response, user: *user, channel: *channel}
}

// todo email and otp
func (uc *UserController) Register(c *gin.Context) {
	var request validator.UserCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		uc.response.ValidatorFail(c, INVALID_PARAM_ERROR)
		return
	}
	err := uc.user.CreateUser(c, request)
	if err != nil {
		uc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	uc.response.Success(c, CREATED_SUCCESS)
}

func (uc *UserController) Delete(c *gin.Context) {
	id := c.GetUint("user_id")
	if err := uc.user.DeleteUser(c, id); err != nil {
		uc.response.FailWithError(c, INVALID_PARAM_ERROR, err)
		return
	}
	uc.response.Success(c, DELETED_SUCCESS)
	return
}

func (uc *UserController) Login(c *gin.Context) {
	var request validator.UserLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		uc.response.ValidatorFail(c, INVALID_PARAM_ERROR)
		return
	}
	token, userLogin, err := uc.user.Login(c, request)
	if token == "" || userLogin == nil || err != nil {
		uc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	uc.response.SuccessWithData(c, ACCEPTED_SUCCESS, map[string]interface{}{
		"token": token,
		"info":  userLogin,
	})
}

func (uc *UserController) Logout(c *gin.Context) {
	ip := c.GetString("ip_address")
	if ip != c.ClientIP() {
		uc.response.Fail(c, VALIDATION_ERROR)
		return
	}
	err := uc.user.Logout(c, c.ClientIP())
	if err != nil {
		uc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	uc.response.Success(c, SUCCESS)
}

func (uc *UserController) Join(c *gin.Context) {
	id := c.GetUint("user")
	var request validator.UserJoinRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		uc.response.ValidatorFail(c, INVALID_PARAM_ERROR)
		return
	}
	token, userJoin, err := uc.user.JoinChannel(c, id, request)
	if token == "" || userJoin == nil || err != nil {
		uc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	user, err := uc.user.ReadUser(c, id)
	if err != nil {
		uc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	err = uc.channel.UserJoin(c, request.ChannelID, *user)
	if err != nil {
		uc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	uc.response.SuccessWithData(c, ACCEPTED_SUCCESS, map[string]interface{}{
		"token": token,
		"info":  userJoin,
	})
}

func (uc *UserController) Leave(c *gin.Context) {
	id := c.GetUint("user")
	channelId := c.GetUint("channel")
	if err := uc.user.LeaveChannel(c, id, channelId); err != nil {
		uc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	if err := uc.channel.UserLeave(c, channelId, id); err != nil {
		uc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	uc.response.Success(c, SUCCESS)
}

func (uc *UserController) Comment(c *gin.Context) {
	var request validator.UserCommentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		uc.response.ValidatorFail(c, INVALID_PARAM_ERROR)
		return
	}
	id := c.GetUint("user_id")
	channelId := c.GetUint("channel_id")
	if err := uc.channel.CommentMessage(c, id, channelId, request); err != nil {
		uc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	if err := uc.message.CreateMessage(c, id, request.Content); err != nil {
		uc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	uc.response.Success(c, SUCCESS)
}

func (uc *UserController) Edite(c *gin.Context) {
	var request validator.UserEditeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		uc.response.ValidatorFail(c, INVALID_PARAM_ERROR)
		return
	}
	id := c.GetUint("user_id")
	channelId := c.GetUint("channel_id")
	if err := uc.channel.EditeMessage(c, id, channelId, request.NewContent); err != nil {
		uc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	if err := uc.message.UpdateMessage(c, id, "content", request.NewContent); err != nil {
		uc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	uc.response.Success(c, SUCCESS)
}

func (uc *UserController) Regain(c *gin.Context) {
	var request validator.UserRegainRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		uc.response.ValidatorFail(c, INVALID_PARAM_ERROR)
		return
	}
	id := c.GetUint("channel_id")
	if err := uc.channel.RegainMessage(c, id, request.MessageID); err != nil {
		uc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	if err := uc.message.DeleteMessage(c, request.MessageID); err != nil {
		uc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	uc.response.Success(c, SUCCESS)
}
