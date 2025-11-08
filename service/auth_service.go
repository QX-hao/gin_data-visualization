package service

import (
	"errors"
	"fmt"
	"time"

	"gin_data-visualization/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService 认证服务接口
type AuthService interface {
	Login(username, password string) (*model.User, string, error)
	Register(req *model.RegisterRequest) (*model.User, error)
	GenerateToken(user *model.User) (string, error)
}

// authService 认证服务实现
type authService struct {
	db *gorm.DB
}

// NewAuthService 创建认证服务实例
func NewAuthService(db *gorm.DB) AuthService {
	return &authService{db: db}
}

// Login 用户登录
func (s *authService) Login(username, password string) (*model.User, string, error) {
	var user model.User
	
	// 查询用户
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("用户名或密码错误")
		}
		return nil, "", err
	}
	
	// 检查用户状态
	if user.Status != 1 {
		return nil, "", errors.New("用户已被禁用")
	}
	
	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("用户名或密码错误")
	}
	
	// 生成token
	token, err := s.GenerateToken(&user)
	if err != nil {
		return nil, "", err
	}
	
	return &user, token, nil
}

// Register 用户注册
func (s *authService) Register(req *model.RegisterRequest) (*model.User, error) {
	// 检查用户名是否已存在
	var existingUser model.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}
	
	// 检查邮箱是否已存在
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("邮箱已存在")
	}
	
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Nickname: req.Nickname,
		Status:   1,
	}
	
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	
	return &user, nil
}

// GenerateToken 生成JWT token（简化版，实际项目中应使用jwt-go等库）
func (s *authService) GenerateToken(user *model.User) (string, error) {
	// 这里使用简单的token生成，实际项目中应该使用JWT
	// 返回一个包含用户ID和时间戳的简单token
	token := fmt.Sprintf("%d-%d", user.ID, time.Now().Unix())
	return token, nil
}

// HashPassword 密码加密（辅助函数）
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword 验证密码（辅助函数）
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}