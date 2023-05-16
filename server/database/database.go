package database

import (
	"ariga.io/atlas/sql/sqltool"
	"bufio"
	"context"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	"fmt"
	goMigrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/chain"
	"github.com/loomi-labs/star-scope/ent/migrate"
	"github.com/pkg/errors"
	"github.com/shifty11/go-logger/log"
	"os"
	"path/filepath"
)

var dbClient *ent.Client

var (
	dbType     = "postgres"
	dbHost     = "localhost"
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "star-scope-db"
	dbPort     = "5432"
	dbSSLMode  = "disable"
	dbTimezone = "Europe/Zurich"
)

func DbCon() string {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=%v&TimeZone=%v", dbType, dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode, dbTimezone)
	}
	return dsn
}

func connect() *ent.Client {
	if dbClient == nil {
		newClient, err := ent.Open("postgres", DbCon())
		if err != nil {
			log.Sugar.Panic("failed to connect to server ", err)
		}
		dbClient = newClient
	}
	return dbClient
}

func Close() {
	if dbClient != nil {
		err := dbClient.Close()
		if err != nil {
			log.Sugar.Error(err)
		}
	}
}

type TxContextValue struct {
	Tx           *ent.Tx
	IsCommited   bool
	IsRolledBack bool
}

func getClient(ctx context.Context, client *ent.Client) *ent.Client {
	if ctx.Value(common.ContextKeyTx) != nil {
		return ctx.Value(common.ContextKeyTx).(TxContextValue).Tx.Client()
	}
	return client
}

func startTx(ctx context.Context, client *ent.Client) (context.Context, error) {
	if ctx.Value(common.ContextKeyTx) != nil {
		return nil, errors.New("transaction already started")
	}
	tx, err := client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	val := TxContextValue{
		Tx:         tx,
		IsCommited: false,
	}
	ctx = context.WithValue(ctx, common.ContextKeyTx, val)
	return ctx, nil
}

func RollbackTxIfUncommitted(ctx context.Context) (context.Context, error) {
	if ctx.Value(common.ContextKeyTx) == nil {
		return ctx, errors.New("transaction not started")
	}
	val := ctx.Value(common.ContextKeyTx).(TxContextValue)
	if val.IsCommited {
		return ctx, nil
	}
	if val.IsRolledBack {
		return ctx, nil
	}
	err := val.Tx.Rollback()
	if err != nil {
		log.Sugar.Error(err)
	}
	val.IsRolledBack = true
	ctx = context.WithValue(ctx, common.ContextKeyTx, val)
	return ctx, nil
}

func CommitTx(ctx context.Context) (context.Context, error) {
	if ctx.Value(common.ContextKeyTx) == nil {
		return ctx, errors.New("transaction not started")
	}
	val := ctx.Value(common.ContextKeyTx).(TxContextValue)
	if val.IsCommited {
		return ctx, nil
	}
	err := val.Tx.Commit()
	if err != nil {
		return ctx, errors.Wrap(err, "committing transaction")
	}
	val.IsCommited = true
	ctx = context.WithValue(ctx, common.ContextKeyTx, val)
	return ctx, nil
}

func withTx(client *ent.Client, ctx context.Context, fn func(tx *ent.Tx) error) error {
	_, err := withTxResult(client, ctx, func(tx *ent.Tx) (*interface{}, error) {
		return nil, fn(tx)
	})
	return err
}

func withTxResult[T any](client *ent.Client, ctx context.Context, fn func(tx *ent.Tx) (*T, error)) (*T, error) {
	tx, err := client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if v := recover(); v != nil {
			//goland:noinspection GoUnhandledErrorResult
			tx.Rollback()
			panic(v)
		}
	}()
	result, err := fn(tx)
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = errors.Wrapf(err, "rolling back transaction: %v", rerr)
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrapf(err, "committing transaction: %v", err)
	}
	return result, nil
}

func getMigrationChecksum(migrationPath string) string {
	file, err := os.Open(filepath.Join(migrationPath, "atlas.sum"))
	if err != nil {
		return ""
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	scanner := bufio.NewScanner(file)

	const maxCapacity int = 100 // required line length
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func CreateMigrations(dbCon string) {
	ctx := context.Background()
	// Create a local migration directory able to understand golang-migrate migration files for replay.
	migrationPath := "database/migrations"
	dir, err := sqltool.NewGolangMigrateDir(migrationPath)
	if err != nil {
		log.Sugar.Fatalf("failed creating atlas migration directory: %v", err)
	}

	checksum := getMigrationChecksum(migrationPath)

	// Write migration diff.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                         // provide migration directory
		schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
		schema.WithDialect(dialect.Postgres),        // Ent dialect to use
		schema.WithDropIndex(true),                  // Drop index if exists
		schema.WithDropColumn(true),                 // Drop column if exists
	}

	err = migrate.NamedDiff(ctx, dbCon, "migration", opts...)
	if err != nil {
		log.Sugar.Fatalf("failed generating migration file: %v", err)
	}
	if checksum == getMigrationChecksum(migrationPath) {
		log.Sugar.Info("no changes detected")
	} else {
		log.Sugar.Info("migrations created successfully")
	}
}

func MigrateDb() {
	m, err := goMigrate.New("file://database/migrations/", DbCon())
	if err != nil {
		log.Sugar.Panicf("failed to migrate database: %v", err)
	}
	err = m.Up()
	if err != nil {
		if err == goMigrate.ErrNoChange {
			log.Sugar.Info("no migration needed")
		} else {
			log.Sugar.Panicf("failed to migrate database: %v", err)
		}
	} else {
		log.Sugar.Info("database migrated successfully")
	}
}

func InitDb() {
	client := connect()
	ctx := context.Background()
	doesOsmosisExist := client.Chain.
		Query().
		Where(chain.NameEQ("Osmosis")).
		ExistX(ctx)
	if !doesOsmosisExist {
		_, err := client.Chain.
			Create().
			SetName("Osmosis").
			SetPath("osmosis").
			SetIndexingHeight(0).
			SetHasCustomIndexer(true).
			SetImage("").
			Save(ctx)
		if err != nil {
			log.Sugar.Panicf("failed to init database: %v", err)
		}
	}
	doesCosmosExist := client.Chain.
		Query().
		Where(chain.NameEQ("Cosmos")).
		ExistX(ctx)
	if !doesCosmosExist {
		_, err := client.Chain.
			Create().
			SetName("Cosmos").
			SetPath("cosmoshub").
			SetIndexingHeight(0).
			SetHasCustomIndexer(false).
			SetImage("").
			Save(ctx)
		if err != nil {
			log.Sugar.Panicf("failed to init database: %v", err)
		}
	}
	log.Sugar.Info("database initialized successfully")
}

type DbManagers struct {
	UserManager          *UserManager
	ProjectManager       *ProjectManager
	ChainManager         *ChainManager
	EventListenerManager *EventListenerManager
}

func NewDefaultDbManagers() *DbManagers {
	client := connect()
	return NewCustomDbManagers(client)
}

func NewCustomDbManagers(client *ent.Client) *DbManagers {
	userManager := NewUserManager(client)
	projectManager := NewProjectManager(client)
	chainManager := NewChainManager(client)
	eventListenerManager := NewEventListenerManager(client)
	return &DbManagers{
		UserManager:          userManager,
		ProjectManager:       projectManager,
		ChainManager:         chainManager,
		EventListenerManager: eventListenerManager,
	}
}
