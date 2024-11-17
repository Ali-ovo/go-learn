package main

import (
	"go-learn/ch19/telemetry/ch05_gorm/model"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/tracing"
)

const (
	traceName = "mxshop-otel"
)

var tp *trace.TracerProvider

func tracerProvider() error {
	// 创建 Jaeger 导出器
	url := "http://192.168.189.128:14268/api/traces"
	jexp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tp = trace.NewTracerProvider(
		trace.WithBatcher(jexp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL, // 必填参数 可以指定是多少版本
			semconv.ServiceNameKey.String("shop-user"),
			attribute.String("environment", "dev"),
			attribute.Int("ID", 1),
		)),
	)
	// 设置全局的提取器
	otel.SetTracerProvider(tp)
	// 设置全局的传播器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}

func Server(c *gin.Context) {
	dsn := "root:123456@tcp(192.168.189.128:3306)/user_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			Colorful: true,        // 彩色打印
			LogLevel: logger.Info, // log level
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}

	if err := db.WithContext(c.Request.Context()).Model(model.User{}).Where("id = ?", 1).First(&model.User{}).Error; err != nil {
		panic(err)
	}

	if err := db.WithContext(c.Request.Context()).Model(model.User{}).Where("id = ?", 2).First(&model.User{}).Error; err != nil {
		panic(err)
	}

	time.Sleep(time.Duration(500) * time.Millisecond)

	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func main() {
	_ = tracerProvider()
	r := gin.Default()
	r.Use(otelgin.Middleware("my-server"))
	r.GET("/", func(c *gin.Context) {})
	r.GET("/server", Server)
	r.Run(":8090")
}
