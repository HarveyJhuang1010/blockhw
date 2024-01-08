package block

import (
	"context"
	"testing"

	"github.com/HarveyJhuang1010/blockhw/internal/appcontext"
	"github.com/HarveyJhuang1010/blockhw/internal/config"
	"github.com/HarveyJhuang1010/blockhw/internal/logging"
	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/HarveyJhuang1010/blockhw/internal/models/po"
	"github.com/HarveyJhuang1010/blockhw/internal/repository/block"
	"github.com/HarveyJhuang1010/blockhw/internal/wrapper/database"
	"github.com/HarveyJhuang1010/blockhw/internal/wrapper/evmcli"
	"github.com/stretchr/testify/suite"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type blockTestApp struct {
	dig.In

	UseCase bo.BlockUseCase
}

type blockTestSuite struct {
	suite.Suite
	app    *blockTestApp
	logger *zap.Logger
	cfg    *config.Config
	rdb    *database.DB
	dbName string
	ctx    *appcontext.AppContext
}

func (s *blockTestSuite) SetupSuite() {
	// Init
	cfg := config.NewConfig()
	logger := logging.NewZapLogger("log-test")
	appCtx := appcontext.New(context.Background())
	appcontext.SetLogger(logger)
	s.logger = logger
	s.ctx = &appCtx
	s.cfg = cfg

	evmcli.Initialize(cfg.Ethereum)
	//database.Initialize(cfg.Database)
	//db := database.GetDB()

	// Init test database
	var testDBName = "testBlock"

	db, err := database.InitTestPostgresSQL(s.T(), cfg.Database, testDBName)
	s.Require().Nil(err)
	s.rdb = db
	s.dbName = testDBName

	// Migration
	s.Require().Nil(db.AutoMigrate(
		&po.Block{},
		&po.Transaction{},
		&po.TransactionLog{},
	))

	// Init Your Provider
	binder := dig.New()
	s.Require().Nil(binder.Provide(database.NewDatabaseClient))
	s.Require().Nil(binder.Provide(evmcli.NewEVMClient))
	s.Require().Nil(binder.Provide(block.NewBlockRepo))
	s.Require().Nil(binder.Provide(newBlockUseCase))
	s.Require().Nil(binder.Invoke(func(app blockTestApp) {
		s.app = &app
	}))
}

func (s *blockTestSuite) SetupTest() {
}

func (s *blockTestSuite) TearDownTest() {
}

func (s *blockTestSuite) TearDownSuite() {
	database.RemoveTestDB(s.T(), config.NewConfig().Database, s.dbName)
	database.Finalize()
	evmcli.Finalize()
}

func (s *blockTestSuite) TestSyncBlock() {
	s.Require().Nil(s.app.UseCase.SyncBlockByNum(s.ctx, 18962369))
}

func TestBlock(t *testing.T) {
	suite.Run(t, &blockTestSuite{})
}
