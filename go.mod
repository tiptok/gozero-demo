module zero-demo

go 1.16

require (
	github.com/Masterminds/squirrel v1.5.3
	github.com/go-pg/pg/v10 v10.10.6
	github.com/golang-jwt/jwt/v4 v4.4.1
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/copier v0.3.5
	github.com/pkg/errors v0.9.1
	github.com/shopspring/decimal v1.3.1
	github.com/silenceper/wechat/v2 v2.1.3
	github.com/sony/sonyflake v1.0.0
	github.com/tiptok/gocomm v1.0.12
	github.com/zeromicro/go-zero v1.3.4
	google.golang.org/grpc v1.47.0
	gorm.io/gorm v1.23.4
)

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/googleapis/gnostic v0.5.5 // indirect
	gorm.io/driver/mysql v1.3.4
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
)

replace github.com/tiptok/gocomm v1.0.12 => D:\\Go\src\learn_project\gocomm

replace (
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
)
