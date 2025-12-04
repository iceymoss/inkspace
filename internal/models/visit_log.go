package models

import (
	"time"
)

// VisitLog 访问日志表（用于统计分析）
type VisitLog struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `gorm:"index:idx_created_at" json:"created_at"`
	UserID      uint      `gorm:"index:idx_user_id" json:"user_id"` // 用户ID（0表示游客）
	IP          string    `gorm:"size:50;index:idx_ip" json:"ip"` // IP地址
	UserAgent   string    `gorm:"size:500" json:"user_agent"` // User Agent
	Path        string    `gorm:"size:255;index:idx_path" json:"path"` // 访问路径
	Method      string    `gorm:"size:10" json:"method"` // 请求方法
	Referer     string    `gorm:"size:500" json:"referer"` // 来源URL
	Duration    int       `json:"duration"` // 请求耗时(毫秒)
	StatusCode  int       `gorm:"index:idx_status_code" json:"status_code"` // HTTP状态码
	Country     string    `gorm:"size:50" json:"country"` // 国家
	Province    string    `gorm:"size:50" json:"province"` // 省份
	City        string    `gorm:"size:50" json:"city"` // 城市
	Browser     string    `gorm:"size:50" json:"browser"` // 浏览器
	OS          string    `gorm:"size:50" json:"os"` // 操作系统
	Device      string    `gorm:"size:50" json:"device"` // 设备类型: desktop, mobile, tablet
}

// VisitLogSummary 访问统计汇总（按天）
type VisitLogSummary struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Date        time.Time `gorm:"uniqueIndex:idx_date;type:date;not null" json:"date"` // 日期
	PV          int       `gorm:"default:0;not null" json:"pv"` // 页面浏览量
	UV          int       `gorm:"default:0;not null" json:"uv"` // 独立访客数
	IP          int       `gorm:"default:0;not null" json:"ip"` // 独立IP数
	NewUsers    int       `gorm:"default:0;not null" json:"new_users"` // 新增用户数
	ArticleView int       `gorm:"default:0;not null" json:"article_view"` // 文章浏览量
	AvgDuration int       `gorm:"default:0;not null" json:"avg_duration"` // 平均访问时长(秒)
	BounceRate  float64   `gorm:"default:0;type:decimal(5,2)" json:"bounce_rate"` // 跳出率
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type VisitLogQuery struct {
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	UserID     uint   `form:"user_id"`
	IP         string `form:"ip"`
	Path       string `form:"path"`
	StatusCode int    `form:"status_code"`
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=20"`
}

type VisitStats struct {
	TodayPV       int     `json:"today_pv"`
	TodayUV       int     `json:"today_uv"`
	YesterdayPV   int     `json:"yesterday_pv"`
	YesterdayUV   int     `json:"yesterday_uv"`
	TotalPV       int     `json:"total_pv"`
	TotalUV       int     `json:"total_uv"`
	AvgDuration   int     `json:"avg_duration"`
	BounceRate    float64 `json:"bounce_rate"`
}

