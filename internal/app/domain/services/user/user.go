package user

import (
	"github.com/nade-harlow/ecom-api/internal/adapter/api/http/dto"
	"github.com/nade-harlow/ecom-api/internal/app/domain/models"
	"github.com/nade-harlow/ecom-api/internal/app/domain/repositories/user"
	"github.com/nade-harlow/ecom-api/internal/app/utils/apperrors"
	utils "github.com/nade-harlow/ecom-api/internal/app/utils/auth"
	"log"
)

type UserService interface {
	RegisterUser(request dto.RegisterUserRequest) (models.User, apperrors.AppError)
	Login(request dto.LoginRequest) (*models.User, apperrors.AppError)
}

type userService struct {
	repository user.UserRepository
}

func NewUserService() UserService {
	return &userService{
		repository: user.NewUserRepository(),
	}
}

func (u *userService) RegisterUser(request dto.RegisterUserRequest) (models.User, apperrors.AppError) {
	var user models.User

	existingUser, err := u.repository.GetUserByEmail(request.Email)
	if err != nil {
		log.Println("failed to get user by email", request.Email, " err: ", err)
		return user, apperrors.InternalServerError("something went wrong")
	}

	if existingUser != nil {
		return user, apperrors.BadRequestError("user with this email already exists")
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		log.Println("failed to hash password", err)
		return user, apperrors.InternalServerError("something went wrong")
	}

	user.NewUser(request.Email, hashedPassword)

	err = u.repository.CreateUser(&user)
	if err != nil {
		log.Println("failed to create user", err)
		return user, apperrors.InternalServerError("something went wrong")
	}

	return user, nil
}

func (u *userService) Login(request dto.LoginRequest) (*models.User, apperrors.AppError) {
	user, err := u.repository.GetUserByEmail(request.Email)
	if err != nil {
		log.Println("failed to get user by email", request.Email, " err: ", err)
		return nil, apperrors.InternalServerError("something went wrong")
	}

	if user == nil {
		return nil, apperrors.BadRequestError("invalid email or password")
	}

	if ok := utils.VerifyPassword(user.PasswordHash, request.Password); !ok {
		return nil, apperrors.BadRequestError("invalid email or password")
	}

	return user, nil
}
