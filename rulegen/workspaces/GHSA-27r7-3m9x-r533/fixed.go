package main

	"github.com/traPtitech/traQ/utils/optional"
	"github.com/traPtitech/traQ/utils/random"
	"github.com/traPtitech/traQ/utils/twemoji"
	gormlogger "gorm.io/gorm/logger"
)

// serveCommand サーバー起動コマンド
			if err != nil {
				logger.Fatal("failed to connect database", zap.Error(err))
			}
			engine.Logger = gormzap.New(logger.Named("gorm")).LogMode(gormlogger.Silent)
			db, err := engine.DB()
			if err != nil {
				logger.Fatal("failed to get *sql.DB", zap.Error(err))
cmd/serve.go | 2 +-
1 file changed, 1 insertion(+), 1 deletion(-)
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	gormlogger "gorm.io/gorm/logger"

	"github.com/traPtitech/traQ/event"
	"github.com/traPtitech/traQ/repository"
	"github.com/traPtitech/traQ/utils/optional"
	"github.com/traPtitech/traQ/utils/random"
	"github.com/traPtitech/traQ/utils/twemoji"
)

// serveCommand サーバー起動コマンド
