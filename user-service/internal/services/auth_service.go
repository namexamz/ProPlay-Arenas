package service

import (
	"errors"
	"log/slog"
	"time"
	"user-service/internal/models"
	"user-service/internal/repository"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	logger    *slog.Logger
	jwtSecret []byte
	userRepo  repository.UserRepository
}

func NewAuthService(
	logger *slog.Logger,
	secret string,
	userRepo repository.UserRepository,
) *AuthService {
	return &AuthService{
		logger:    logger.With(slog.String("layer", "service"), slog.String("service", "auth")),
		jwtSecret: []byte(secret),
		userRepo:  userRepo,
	}
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	s.logger.Debug("generate token started",
		slog.Uint64("user_id", uint64(user.ID)),
		slog.String("role", string(user.Role)),
	)

	claims := models.Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		s.logger.Error("generate token failed",
			slog.Uint64("user_id", uint64(user.ID)),
			slog.Any("error", err),
		)
		return "", err
	}

	s.logger.Debug("generate token completed",
		slog.Uint64("user_id", uint64(user.ID)),
	)
	return signedToken, nil
}

func (s *AuthService) ParseToken(tokenStr string) (*models.Claims, error) {
	s.logger.Debug("parse token started")

	token, err := jwt.ParseWithClaims(tokenStr, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		s.logger.Warn("parse token failed", slog.Any("error", err))
		return nil, err
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		s.logger.Warn("invalid token")
		return nil, errors.New("invalid token")
	}

	s.logger.Debug("parse token completed",
		slog.Uint64("user_id", uint64(claims.UserID)),
		slog.String("role", string(claims.Role)),
	)
	return claims, nil
}

func (s *AuthService) RegisterUser(req models.RegisterRequest) (string, error) {
	s.logger.Info("user registration started",
		slog.String("email", req.Email),
	)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("password hashing failed",
			slog.String("email", req.Email),
			slog.Any("error", err),
		)
		return "", err
	}

	user := &models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     models.RoleClient,
	}

	if err := s.userRepo.Create(user); err != nil {
		s.logger.Error("user creation failed",
			slog.String("email", req.Email),
			slog.Any("error", err),
		)
		return "", err
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return "", err
	}

	s.logger.Info("user registration completed",
		slog.Uint64("user_id", uint64(user.ID)),
		slog.String("email", user.Email),
	)
	return token, nil
}

func (s *AuthService) LoginUser(req models.LoginRequest) (string, error) {
	s.logger.Info("user login started",
		slog.String("email", req.Email),
	)

	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil || user == nil {
		s.logger.Warn("login failed: user not found",
			slog.String("email", req.Email),
		)
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logger.Warn("login failed: invalid password",
			slog.Uint64("user_id", uint64(user.ID)),
		)
		return "", errors.New("invalid email or password")
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return "", err
	}

	s.logger.Info("user login completed",
		slog.Uint64("user_id", uint64(user.ID)),
		slog.String("email", user.Email),
	)
	return token, nil
}
