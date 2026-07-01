package service

import (
	"context"
	"errors"
	"time"

	"project/internal/models"
	"project/internal/repository"
	"project/internal/utils"
	"project/pkg/cache"
	"project/pkg/jwt"
	"project/pkg/password"
	smtp2 "project/pkg/smtp"
)

type UserService struct {
	repo *repository.UserRepository
	cache cache.MemoryCache
	smtp  *smtp2.SMTP
}

func NewUserService(repo *repository.UserRepository, cache cache.MemoryCache, smtp *smtp2.SMTP) *UserService {
	return &UserService{repo: repo,
	cache: cache,
	smtp: smtp,
	}
}
func (s *UserService) Register(ctx context.Context, request models.RegisterRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	exists, err := s.repo.ExistsByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("user with this email already exists")
	}

	_, ok := s.cache.Get(request.Email)
	if ok {
		return errors.New("user with this email already exists")
	}

	passwordHash, err := password.Hash(request.Password)
	if err != nil {
		return err
	}

	otp := utils.GenerateOTP()
	err = s.smtp.SendOTP(ctx, request.Email, otp)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: passwordHash,
		Role:     models.UserRole,
		OtpCode:  otp,
	}

	s.cache.Set(request.Email, user, time.Minute*5)
	return nil
}

func (s *UserService) Login(ctx context.Context, request models.LoginRequest) (response models.LoginResponse, err error) {
	err = request.Validate()
	if err != nil {
		return models.LoginResponse{}, err
	}

	user, err := s.repo.GetByEmail(ctx, request.Email)
	if err != nil {
		return models.LoginResponse{}, err
	}

	err = password.Compare(user.Password, request.Password)
	if err != nil {
		return models.LoginResponse{}, err
	}

	accessToken, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return models.LoginResponse{}, err
	}

	refreshToken := utils.GenerateRand()
	refreshTokenHash := utils.HashRefreshToken(refreshToken)

	refreshTokenEntity := models.RefreshToken{
		UserID:    int64(user.ID),
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	err = s.repo.AddRefreshToken(ctx, refreshTokenEntity)
	if err != nil {
		return models.LoginResponse{}, err
	}

	response = models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}

func (s *UserService) GetMe(ctx context.Context, userID int64) (user models.User, err error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *UserService) UpdateMe(ctx context.Context, userID int, request models.UpdateRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	user := models.User{
		ID:    userID,
		Name:  request.Name,
		Phone: request.Phone,
	}

	err = s.repo.Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteMe(ctx context.Context, userID int64) error {
	err := s.repo.Delete(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) CreateOrder(ctx context.Context, userID int, request models.CreateOrderRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	err := s.repo.CreateOrder(ctx, userID, request.Description, request.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetMyOrders(ctx context.Context, userID int) ([]models.Order, error) {
	return s.repo.GetOrdersByUserID(ctx, userID)
}

func (s *UserService) GetOrderByID(ctx context.Context, orderID int, userID int) (models.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return models.Order{}, err
	}

	if order.UserID != userID {
		return models.Order{}, errors.New("Error")
	}

	return order, nil
}

func (s *UserService) UpdateOrder(ctx context.Context, orderID int, userID int, request models.UpdateOrderRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return errors.New("Error")
	}

	order.Description = request.Description
	order.Amount = request.Amount
	order.Status = request.Status

	err = s.repo.UpdateOrder(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteOrder(ctx context.Context, orderID int, userID int) error {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return errors.New("Error")
	}

	err = s.repo.DeleteOrder(ctx, orderID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) AdminGetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.AdminGetAllUsers(ctx)
}

func (s *UserService) AdminGetAllOrders(ctx context.Context) ([]models.Order, error) {
	return s.repo.AdminGetAllOrders(ctx)
}

func (s *UserService) AdminUpdateRole(ctx context.Context, targetUserID int, role string) error {
	return s.repo.AdminUpdateRole(ctx, targetUserID, role)
}

func (s *UserService) UpdateUserPassword(ctx context.Context, userID int64, request models.UpdateUserPassword) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("Internal error")
	}
	err = password.Compare(user.Password, request.OldPassword)
	if err != nil {
		return err
	}

	newPassword,err := password.Hash(request.NewPassword)
	if err != nil {
		return err
	}

	err = s.repo.UpdateUserPassword(ctx, userID,newPassword)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) CancelOrder(ctx context.Context, orderID int, userID int) error {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return errors.New("access denied")
	}

	if order.Status == "paid" || order.Status == "completed" {
		return errors.New("cannot cancel paid or completed order")
	}

	return s.repo.UpdateOrderStatus(ctx, int64(orderID), "cancelled")
}

func (s *UserService) GetUserAndOrders(ctx context.Context, userID int) (models.UserAndOrders, error) {
	var response models.UserAndOrders
	var err error
	response.User, err = s.repo.GetByID(ctx, int64(userID))
	if err != nil {
		return models.UserAndOrders{},err
	}

	response.Orders, err = s.repo.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return models.UserAndOrders{}, err
	}
	return response, nil
}


func (s *UserService) Verify(ctx context.Context, request models.VerifyRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	cacheInfo, ok := s.cache.Get(request.Email)
	if !ok {
		return errors.New("user not found")
	}

	user, ok := cacheInfo.(models.User)
	if !ok {
		return errors.New("user not found")
	}

	if user.AttemptOTP >= 3 {
		return errors.New("user is too many attempts")
	}

	if request.Otp != user.OtpCode {
		user.AttemptOTP++
		s.cache.Set(request.Email, user, time.Minute*5)
		return errors.New("invalid otp")
	}

	err = s.repo.Add(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Refresh(ctx context.Context, request models.RefreshRequest) (models.LoginResponse, error) {
	if err := request.Validate(); err != nil {
		return models.LoginResponse{}, err
	}

	oldHash := utils.HashRefreshToken(request.RefreshToken)

	storedToken, err := s.repo.GetRefreshToken(ctx, oldHash)
	if err != nil {
		return models.LoginResponse{}, errors.New("unauthorized: invalid refresh token")
	}

	if time.Now().After(storedToken.ExpiresAt) {
		_ = s.repo.DeleteRefreshToken(ctx, oldHash)
		return models.LoginResponse{}, errors.New("unauthorized: refresh token expired")
	}

	err = s.repo.DeleteRefreshToken(ctx, oldHash)
	if err != nil {
		return models.LoginResponse{}, err
	}

	user, err := s.repo.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return models.LoginResponse{}, err
	}

	newAccessToken, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return models.LoginResponse{}, err
	}

	newRawRefreshToken := utils.GenerateRand()
	newHash := utils.HashRefreshToken(newRawRefreshToken)

	newTokenModel := models.RefreshToken{
		UserID:    storedToken.UserID,
		TokenHash: newHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	err = s.repo.AddRefreshToken(ctx, newTokenModel)
	if err != nil {
		return models.LoginResponse{}, err
	}

	response := models.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRawRefreshToken,
	}
	
	return response, nil
}

func (s *UserService) DeleteToken(ctx context.Context, hash string) error {
    return s.repo.DeleteRefreshToken(ctx, hash)
}