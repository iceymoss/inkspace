package models

import (
	"errors"
	"time"
)

const (
	DefaultUITheme     = "magazine"
	DefaultColorScheme = "system"
)

type UserAppearance struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	UITheme     string    `gorm:"size:20;default:'magazine';not null" json:"ui_theme"`
	ColorScheme string    `gorm:"size:10;default:'system';not null" json:"color_scheme"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserAppearanceRequest struct {
	UITheme     string `json:"ui_theme" binding:"required"`
	ColorScheme string `json:"color_scheme" binding:"required"`
}

type UserAppearanceResponse struct {
	UITheme     string `json:"ui_theme"`
	ColorScheme string `json:"color_scheme"`
}

func DefaultUserAppearanceResponse() *UserAppearanceResponse {
	return &UserAppearanceResponse{
		UITheme:     DefaultUITheme,
		ColorScheme: DefaultColorScheme,
	}
}

func (r *UserAppearanceRequest) Validate() error {
	if r.UITheme != DefaultUITheme {
		return errors.New("无效或尚未开放的界面主题")
	}

	switch r.ColorScheme {
	case "system", "light", "dark":
		return nil
	default:
		return errors.New("无效的明暗模式")
	}
}

func (a *UserAppearance) ToResponse() *UserAppearanceResponse {
	return &UserAppearanceResponse{
		UITheme:     a.UITheme,
		ColorScheme: a.ColorScheme,
	}
}
