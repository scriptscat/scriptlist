module github.com/scriptscat/scriptlist

go 1.18

require (
	github.com/codfrm/cago v0.0.0-20221210114855-bac5b1062a49
	github.com/gin-gonic/gin v1.8.1
	github.com/go-gormigrate/gormigrate/v2 v2.0.2
	github.com/gorilla/csrf v1.7.1
	github.com/swaggo/swag v1.8.6
	go.opentelemetry.io/otel v1.11.1
	go.uber.org/zap v1.24.0
	gorm.io/datatypes v1.0.7
	gorm.io/gorm v1.24.1
)

replace github.com/codfrm/cago => ../../cago

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/asaskevich/EventBus v0.0.0-20200907212545-49d423059eef // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gin-contrib/sessions v0.0.5 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.11.1 // indirect
	github.com/go-redis/redis/extra/rediscmd/v9 v9.0.0-rc.2 // indirect
	github.com/go-redis/redis/extra/redisotel/v9 v9.0.0-rc.2 // indirect
	github.com/go-redis/redis/v9 v9.0.0-rc.2 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/goccy/go-json v0.10.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.14.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nsqio/go-nsq v1.1.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/scriptscat/cloudcat v0.1.0 // indirect
	github.com/swaggo/gin-swagger v1.3.1 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	github.com/uptrace/uptrace-go v1.11.6 // indirect
	go.etcd.io/etcd/api/v3 v3.5.6 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.6 // indirect
	go.etcd.io/etcd/client/v3 v3.5.6 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.36.4 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.9.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.11.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.33.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.33.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.11.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.11.1 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.11.1 // indirect
	go.opentelemetry.io/otel/metric v0.33.0 // indirect
	go.opentelemetry.io/otel/sdk v1.11.1 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.33.0 // indirect
	go.opentelemetry.io/otel/trace v1.11.1 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.4.0 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	golang.org/x/tools v0.2.0 // indirect
	google.golang.org/genproto v0.0.0-20221207170731-23e4bf6bdc37 // indirect
	google.golang.org/grpc v1.51.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.4.3 // indirect
	gorm.io/plugin/opentelemetry v0.1.0 // indirect
)
