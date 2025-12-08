package repository

import (
	"context"
	"database/sql"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestModel struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100)"`
	Email     string    `gorm:"type:varchar(100)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (TestModel) TableName() string {
	return "test_models"
}

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	return db, mock, sqlDB
}

func TestCreate(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)
	ctx := context.Background()

	testModel := &TestModel{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "test_models"`)).
		WithArgs(testModel.Name, testModel.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(ctx, testModel)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateWithCanceledContext(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)

	// Context cancelado
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancelar inmediatamente

	testModel := &TestModel{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// No esperamos ninguna query porque el context está cancelado
	err := repo.Create(ctx, testModel)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateWithTimeout(t *testing.T) {
	db, _, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)

	// Context con timeout muy corto
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond) // Asegurar que expire

	testModel := &TestModel{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	err := repo.Create(ctx, testModel)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestFindByID(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)
	ctx := context.Background()

	expectedModel := TestModel{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
		AddRow(expectedModel.ID, expectedModel.Name, expectedModel.Email, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_models" WHERE "test_models"."id" = $1 ORDER BY "test_models"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	result, err := repo.FindByID(ctx, uint(1))

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedModel.ID, result.ID)
	assert.Equal(t, expectedModel.Name, result.Name)
	assert.Equal(t, expectedModel.Email, result.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByIDWithCanceledContext(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := repo.FindByID(ctx, uint(1))

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByIDNotFound(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_models" WHERE "test_models"."id" = $1 ORDER BY "test_models"."id" LIMIT $2`)).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.FindByID(ctx, uint(999))

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAll(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
		AddRow(1, "John Doe", "john@example.com", time.Now(), time.Now()).
		AddRow(2, "Jane Smith", "jane@example.com", time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_models"`)).
		WillReturnRows(rows)

	results, err := repo.FindAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "John Doe", results[0].Name)
	assert.Equal(t, "Jane Smith", results[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAllEmpty(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_models"`)).
		WillReturnRows(rows)

	results, err := repo.FindAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, results, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAllWithCanceledContext(t *testing.T) {
	db, _, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	results, err := repo.FindAll(ctx)

	assert.Error(t, err)
	assert.Nil(t, results)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestUpdate(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)
	ctx := context.Background()

	testModel := &TestModel{
		ID:    1,
		Name:  "John Updated",
		Email: "john.updated@example.com",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "test_models"`)).
		WithArgs(testModel.Name, testModel.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), testModel.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, testModel)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateWithCanceledContext(t *testing.T) {
	db, _, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	testModel := &TestModel{
		ID:    1,
		Name:  "John Updated",
		Email: "john.updated@example.com",
	}

	err := repo.Update(ctx, testModel)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestDelete(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "test_models" WHERE "test_models"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, uint(1))

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteWithCanceledContext(t *testing.T) {
	db, _, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.Delete(ctx, uint(1))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestCreateTx(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)
	ctx := context.Background()

	testModel := &TestModel{
		Name:  "Transaction Test",
		Email: "tx@example.com",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "test_models"`)).
		WithArgs(testModel.Name, testModel.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := db.Transaction(func(tx *gorm.DB) error {
		return repo.CreateTx(ctx, tx, testModel)
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTxWithCanceledContext(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	testModel := &TestModel{
		Name:  "Transaction Test",
		Email: "tx@example.com",
	}

	mock.ExpectBegin()
	mock.ExpectRollback()

	err := db.Transaction(func(tx *gorm.DB) error {
		return repo.CreateTx(ctx, tx, testModel)
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTx(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)
	ctx := context.Background()

	testModel := &TestModel{
		ID:    1,
		Name:  "TX Updated",
		Email: "tx.updated@example.com",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "test_models"`)).
		WithArgs(testModel.Name, testModel.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), testModel.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := db.Transaction(func(tx *gorm.DB) error {
		return repo.UpdateTx(ctx, tx, testModel)
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTx(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "test_models" WHERE "test_models"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := db.Transaction(func(tx *gorm.DB) error {
		return repo.DeleteTx(ctx, tx, uint(1))
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRollback(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)
	ctx := context.Background()

	testModel := &TestModel{
		Name:  "Rollback Test",
		Email: "rollback@example.com",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "test_models"`)).
		WithArgs(testModel.Name, testModel.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	err := db.Transaction(func(tx *gorm.DB) error {
		return repo.CreateTx(ctx, tx, testModel)
	})

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestContextWithTimeout(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)

	// Context con timeout de 100ms
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
		AddRow(1, "John", "john@example.com", time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_models" WHERE "test_models"."id" = $1 ORDER BY "test_models"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillDelayFor(200 * time.Millisecond). // Simula query lenta
		WillReturnRows(rows)

	result, err := repo.FindByID(ctx, uint(1))

	assert.Error(t, err)
	assert.Nil(t, result)
	// The error can be either "context deadline exceeded" or "canceling query due to user request"
	// depending on timing and database implementation
	errStr := err.Error()
	assert.True(t,
		strings.Contains(errStr, "context deadline exceeded") || strings.Contains(errStr, "cancel"),
		"expected error to contain 'context deadline exceeded' or 'cancel', got: %s", errStr)
}

func TestContextWithValue(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := NewRepository[TestModel](db)

	// Context con valor
	ctx := context.WithValue(context.Background(), "request_id", "test-123")

	rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
		AddRow(1, "John", "john@example.com", time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_models" WHERE "test_models"."id" = $1 ORDER BY "test_models"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	result, err := repo.FindByID(ctx, uint(1))

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verificar que el context mantiene el valor
	requestID := ctx.Value("request_id").(string)
	assert.Equal(t, "test-123", requestID)

	assert.NoError(t, mock.ExpectationsWereMet())
}
