package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nade-harlow/ecom-api/internal/adapter/api/http/dto"
	"github.com/nade-harlow/ecom-api/internal/adapter/api/http/response"
	"github.com/nade-harlow/ecom-api/internal/app/domain/services/user"
	"github.com/nade-harlow/ecom-api/internal/app/utils/apperrors"
	utils "github.com/nade-harlow/ecom-api/internal/app/utils/auth"
	"github.com/nade-harlow/ecom-api/internal/app/utils/helper"
)

type UserHandler struct {
	userService user.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: user.NewUserService(),
	}
}

func (u *UserHandler) RegisterUser(ctx *gin.Context) {
	var registerUserDto dto.RegisterUserRequest

	if err := ctx.ShouldBind(&registerUserDto); err != nil {
		response.JsonError(ctx, apperrors.BadRequestError(helper.ValidatorFormatErrors(err).Error()))
		return
	}

	if err := helper.ValidateRequestBody(registerUserDto); err != nil {
		response.JsonError(ctx, apperrors.BadRequestError(err.Error()))
		return
	}

	user, err := u.userService.RegisterUser(registerUserDto)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	claims := utils.GetJWTClaims(user.ID, string(user.Role))
	token := utils.GenerateJwt(claims)

	data := struct {
		ID    uuid.UUID `json:"id"`
		Email string    `json:"email"`
		Token string    `json:"token"`
	}{ID: user.ID, Email: user.Email, Token: token}

	response.JsonCreated(ctx, data, "user")
}

func (u *UserHandler) Login(ctx *gin.Context) {
	var loginRequest dto.LoginRequest

	if err := ctx.ShouldBind(&loginRequest); err != nil {
		response.JsonError(ctx, apperrors.BadRequestError(helper.ValidatorFormatErrors(err).Error()))
		return
	}

	if err := helper.ValidateRequestBody(loginRequest); err != nil {
		response.JsonError(ctx, apperrors.BadRequestError(err.Error()))
		return
	}

	user, err := u.userService.Login(loginRequest)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	claims := utils.GetJWTClaims(user.ID, string(user.Role))
	token := utils.GenerateJwt(claims)

	data := struct {
		ID    uuid.UUID `json:"id"`
		Email string    `json:"email"`
		Token string    `json:"token"`
	}{ID: user.ID, Email: user.Email, Token: token}

	response.JsonOk(ctx, data)
}
