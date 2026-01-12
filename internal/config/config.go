package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Admin      AdminConfig      `mapstructure:"admin"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Upload     UploadConfig     `mapstructure:"upload"`
	Pagination PaginationConfig `mapstructure:"pagination"`
	Cache      CacheConfig      `mapstructure:"cache"`
}

type AdminConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	Charset      string `mapstructure:"charset"`
	ParseTime    bool   `mapstructure:"parseTime"`
	Loc          string `mapstructure:"loc"`
	MaxIdleConns int    `mapstructure:"maxIdleConns"`
	MaxOpenConns int    `mapstructure:"maxOpenConns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"poolSize"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expireHours"`
	AdminSecret string `mapstructure:"adminSecret"` // 管理员专用secret
}

type UploadConfig struct {
	StorageType string           `mapstructure:"storageType"` // local or cos
	MaxSize     int64            `mapstructure:"maxSize"`
	AllowTypes  []string         `mapstructure:"allowTypes"`
	SavePath    string           `mapstructure:"savePath"`
	TencentCOS  TencentCOSConfig `mapstructure:"tencentCOS"`
}

type TencentCOSConfig struct {
	BucketURL string `mapstructure:"bucketURL"`
	SecretID  string `mapstructure:"secretID"`
	SecretKey string `mapstructure:"secretKey"`
	Domain    string `mapstructure:"domain"` // CDN域名
}

type PaginationConfig struct {
	PageSize    int `mapstructure:"pageSize"`
	MaxPageSize int `mapstructure:"maxPageSize"`
}

type CacheConfig struct {
	ArticleExpire int `mapstructure:"articleExpire"`
	UserExpire    int `mapstructure:"userExpire"`
}

var AppConfig *Config

func Init() error {
	return InitWithFile("config")
}

func InitWithFile(configName string) error {
	// 1. 首先加载 .env 文件（如果存在）
	loadEnvFile()

	// 2. 设置环境变量支持
	viper.AutomaticEnv()
	// 将环境变量中的下划线转换为点，例如 DATABASE_HOST -> database.host
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 3. 加载 YAML 配置文件
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，尝试使用环境变量
		log.Printf("警告: 无法读取配置文件 %s.yaml，将使用环境变量和默认值: %v\n", configName, err)
	} else {
		log.Printf("已加载配置文件: %s\n", viper.ConfigFileUsed())
	}

	// 4. 绑定环境变量到配置键
	bindEnvVars()

	// 5. 解析配置到结构体
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return err
	}

	return nil
}

// loadEnvFile 加载 .env 文件
func loadEnvFile() {
	envFile := ".env"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		// .env 文件不存在，跳过
		return
	}

	// 读取 .env 文件并设置到环境变量
	file, err := os.Open(envFile)
	if err != nil {
		log.Printf("警告: 无法打开 .env 文件: %v\n", err)
		return
	}
	defer file.Close()

	// 读取文件内容
	content, err := os.ReadFile(envFile)
	if err != nil {
		log.Printf("警告: 无法读取 .env 文件: %v\n", err)
		return
	}

	// 解析 .env 文件内容
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 KEY=VALUE 格式
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// 移除引号（如果有）
			value = strings.Trim(value, `"'`)
			// 设置环境变量
			os.Setenv(key, value)
		}
	}

	log.Printf("已加载 .env 文件: %s\n", envFile)
}

// bindEnvVars 绑定环境变量到配置键
func bindEnvVars() {
	// Server 配置
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.mode", "SERVER_MODE")

	// Admin 配置
	viper.BindEnv("admin.port", "ADMIN_PORT")
	viper.BindEnv("admin.mode", "ADMIN_MODE")

	// Database 配置
	viper.BindEnv("database.host", "DATABASE_HOST")
	viper.BindEnv("database.port", "DATABASE_PORT")
	viper.BindEnv("database.username", "DATABASE_USERNAME")
	viper.BindEnv("database.password", "DATABASE_PASSWORD")
	viper.BindEnv("database.database", "DATABASE_NAME")
	viper.BindEnv("database.charset", "DATABASE_CHARSET")
	viper.BindEnv("database.parseTime", "DATABASE_PARSE_TIME")
	viper.BindEnv("database.loc", "DATABASE_LOC")
	viper.BindEnv("database.maxIdleConns", "DATABASE_MAX_IDLE_CONNS")
	viper.BindEnv("database.maxOpenConns", "DATABASE_MAX_OPEN_CONNS")

	// Redis 配置
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")
	viper.BindEnv("redis.poolSize", "REDIS_POOL_SIZE")

	// JWT 配置
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.adminSecret", "JWT_ADMIN_SECRET")
	viper.BindEnv("jwt.expireHours", "JWT_EXPIRE_HOURS")

	// Upload 配置
	viper.BindEnv("upload.storageType", "UPLOAD_STORAGE_TYPE")
	viper.BindEnv("upload.maxSize", "UPLOAD_MAX_SIZE")
	viper.BindEnv("upload.savePath", "UPLOAD_SAVE_PATH")

	// COS 配置
	viper.BindEnv("upload.tencentCOS.bucketURL", "COS_BUCKET_URL")
	viper.BindEnv("upload.tencentCOS.secretID", "COS_SECRET_ID")
	viper.BindEnv("upload.tencentCOS.secretKey", "COS_SECRET_KEY")
	viper.BindEnv("upload.tencentCOS.domain", "COS_DOMAIN")

	// Cache 配置
	viper.BindEnv("cache.articleExpire", "CACHE_ARTICLE_EXPIRE")
	viper.BindEnv("cache.userExpire", "CACHE_USER_EXPIRE")
}
