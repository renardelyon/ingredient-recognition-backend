package service

import (
	"context"
	"fmt"
	"ingredient-recognition-backend/internal/domain"
	repointerface "ingredient-recognition-backend/internal/repository/repo_interface"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines authentication-related methods
type AuthService interface {
	Register(ctx context.Context, req *domain.UserRegistrationRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req *domain.UserLoginRequest) (*domain.AuthResponse, error)
	GetUserFromToken(ctx context.Context, token string) (*domain.User, error)
}

// authService is a concrete implementation of AuthService
type authService struct {
	userRepo   repointerface.UserRepository
	jwtSecret  string
	expiryTime time.Duration
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo repointerface.UserRepository, jwtSecret string, expiryTime time.Duration) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		expiryTime: expiryTime,
	}
}

// Register registers a new user
func (a *authService) Register(ctx context.Context, req *domain.UserRegistrationRequest) (*domain.AuthResponse, error) {
	// Check if user already exists
	_, err := a.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	user := &domain.User{
		Id:        uuid.New().String(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user to repository
	if err := a.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := a.generateToken(user.Id, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return response without password
	user.Password = ""

	return &domain.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// Login authenticates a user
func (a *authService) Login(ctx context.Context, req *domain.UserLoginRequest) (*domain.AuthResponse, error) {
	// Get user by email
	user, err := a.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, domain.ErrInvalidPassword
	}

	// Generate JWT token
	token, err := a.generateToken(user.Id, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return response without password
	user.Password = ""

	return &domain.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// ValidateToken validates a JWT token
func (a *authService) validateToken(tokenString string) (jwt.Claims, error) {
	claims := &jwt.RegisteredClaims{}
	// TODO: check how this function works
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, domain.ErrUnauthorized
	}

	return claims, nil
}

// GetUserFromToken retrieves user from JWT token
func (a *authService) GetUserFromToken(ctx context.Context, token string) (*domain.User, error) {
	claims, err := a.validateToken(token)
	if err != nil {
		return nil, err
	}

	claim, err := claims.GetSubject()
	if err != nil {
		return nil, err
	}

	user, err := a.userRepo.GetByID(ctx, claim)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// generateToken generates a JWT token
func (a *authService) generateToken(userID, email string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.expiryTime)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
