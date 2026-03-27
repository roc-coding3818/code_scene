module github.com/roc-coding3818/code_scene/modules/user

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/roc-coding3818/code_scene/shared v0.0.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/mysql v1.5.2
	gorm.io/gorm v1.25.5
)

replace github.com/roc-coding3818/code_scene/shared => ../../shared
