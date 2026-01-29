package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server       ServerConfig       `mapstructure:"server"`
    Database     DatabaseConfig     `mapstructure:"database"`
    Redis        RedisConfig        `mapstructure:"redis"`
    JWT          JWTConfig          `mapstructure:"jwt"`
    Shouqianba   ShouqianbaConfig  `mapstructure:"shouqianba"`
    Wechat       WechatConfig       `mapstructure:"wechat"`
    Alipay       AlipayConfig       `mapstructure:"alipay"`
    Log          LogConfig          `mapstructure:"log"`
    OSS          OSSConfig          `mapstructure:"oss"`
    SMS          SMSConfig          `mapstructure:"sms"`
    Email        EmailConfig        `mapstructure:"email"`
    Security     SecurityConfig     `mapstructure:"security"`
    RiskControl  RiskControlConfig  `mapstructure:"risk_control"`
    Monitoring   MonitoringConfig   `mapstructure:"monitoring"`
}

type ServerConfig struct {
    Port         string `mapstructure:"port"`
    Mode         string `mapstructure:"mode"`
    ReadTimeout  int    `mapstructure:"read_timeout"`
    WriteTimeout int    `mapstructure:"write_timeout"`
    Domain       string `mapstructure:"domain"`
    IsHTTPS      bool   `mapstructure:"is_https"`
}

type DatabaseConfig struct {
    Host       string `mapstructure:"host"`
    Port       int    `mapstructure:"port"`
    Username   string `mapstructure:"username"`
    Password   string `mapstructure:"password"`
    Database   string `mapstructure:"database"`
    MaxIdleConn int   `mapstructure:"max_idle_conn"`
    MaxOpenConn int   `mapstructure:"max_open_conn"`
    LogLevel   string `mapstructure:"log_level"`
}

type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
    PoolSize int    `mapstructure:"pool_size"`
}

type JWTConfig struct {
    Secret             string `mapstructure:"secret"`
    ExpireTime         int    `mapstructure:"expire_time"`
    RefreshExpireTime  int    `mapstructure:"refresh_expire_time"`
}

type ShouqianbaConfig struct {
    AppID        string `mapstructure:"app_id"`
    MerchantNo   string `mapstructure:"merchant_no"`
    SecretKey    string `mapstructure:"secret_key"`
    APIURL       string `mapstructure:"api_url"`
    SandboxURL   string `mapstructure:"sandbox_url"`
    Sandbox      bool   `mapstructure:"sandbox"`
    NotifyURL    string `mapstructure:"notify_url"`
}

type WechatConfig struct {
    AppID     string `mapstructure:"app_id"`
    AppSecret string `mapstructure:"app_secret"`
    MchID     string `mapstructure:"mch_id"`
    APIKey    string `mapstructure:"api_key"`
}

type AlipayConfig struct {
    AppID       string `mapstructure:"app_id"`
    PrivateKey  string `mapstructure:"private_key"`
    PublicKey   string `mapstructure:"public_key"`
    GatewayURL  string `mapstructure:"gateway_url"`
}

type LogConfig struct {
    Level    string `mapstructure:"level"`
    Format   string `mapstructure:"format"`
    Output   string `mapstructure:"output"`
    FilePath string `mapstructure:"file_path"`
    MaxSize  int    `mapstructure:"max_size"`
    MaxAge   int    `mapstructure:"max_age"`
    MaxBackups int  `mapstructure:"max_backups"`
}

type OSSConfig struct {
    Endpoint        string `mapstructure:"endpoint"`
    AccessKeyID     string `mapstructure:"access_key_id"`
    AccessKeySecret string `mapstructure:"access_key_secret"`
    Bucket          string `mapstructure:"bucket"`
}

type SMSConfig struct {
    Provider  string `mapstructure:"provider"`
    AccessKey string `mapstructure:"access_key"`
    SecretKey string `mapstructure:"secret_key"`
    SignName  string `mapstructure:"sign_name"`
}

type EmailConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    From     string `mapstructure:"from"`
}

type SecurityConfig struct {
    EnableCORS   bool     `mapstructure:"enable_cors"`
    AllowOrigins []string `mapstructure:"allow_origins"`
    AllowMethods []string `mapstructure:"allow_methods"`
    AllowHeaders []string `mapstructure:"allow_headers"`
    RateLimit    int      `mapstructure:"rate_limit"`
}

type RiskControlConfig struct {
    EnableDeviceFingerprint bool `mapstructure:"enable_device_fingerprint"`
    MaxRegisterPerIP        int  `mapstructure:"max_register_per_ip"`
    MaxTaskPerUser          int  `mapstructure:"max_task_per_user"`
    MaxWithdrawPerDay       int  `mapstructure:"max_withdraw_per_day"`
}

type MonitoringConfig struct {
    EnablePrometheus bool   `mapstructure:"enable_prometheus"`
    PrometheusPort   string `mapstructure:"prometheus_port"`
    EnableTrace      bool   `mapstructure:"enable_trace"`
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
    v := viper.New()
    
    // 设置配置文件路径和类型
    v.SetConfigFile(configPath)
    
    // 设置环境变量前缀
    v.SetEnvPrefix("TASK_PLATFORM")
    v.AutomaticEnv()
    
    // 读取配置文件
    if err := v.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var config Config
    if err := v.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}