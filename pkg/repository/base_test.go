package repository_test

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/iagopicos/game-awards-api/pkg/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestModel struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(100)"`
	Email     string `gorm:"type:varchar(100)"`
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

	repo := repository.NewRepository[TestModel](db)

	testModel := &TestModel{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "test_models"`)).
		WithArgs(testModel.Name, testModel.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(testModel)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByID(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repository.NewRepository[TestModel](db)

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

	result, err := repo.FindByID(uint(1))

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedModel.ID, result.ID)
	assert.Equal(t, expectedModel.Name, result.Name)
	assert.Equal(t, expectedModel.Email, result.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByIDNotFound(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repository.NewRepository[TestModel](db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_models" WHERE "test_models"."id" = $1 ORDER BY "test_models"."id" LIMIT $2`)).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.FindByID(uint(999))

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAll(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repository.NewRepository[TestModel](db)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
		AddRow(1, "John Doe", "john@example.com", time.Now(), time.Now()).
		AddRow(2, "Jane Smith", "jane@example.com", time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_models"`)).
		WillReturnRows(rows)

	results, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "John Doe", results[0].Name)
	assert.Equal(t, "Jane Smith", results[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAllEmpty(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repository.NewRepository[TestModel](db)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "test_models"`)).
		WillReturnRows(rows)

	results, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, results, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repository.NewRepository[TestModel](db)

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

	err := repo.Update(testModel)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repository.NewRepository[TestModel](db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "test_models" WHERE "test_models"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(uint(1))

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTx(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repository.NewRepository[TestModel](db)

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
		return repo.CreateTx(tx, testModel)
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTx(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repository.NewRepository[TestModel](db)

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
		return repo.UpdateTx(tx, testModel)
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTx(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repository.NewRepository[TestModel](db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "test_models" WHERE "test_models"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := db.Transaction(func(tx *gorm.DB) error {
		return repo.DeleteTx(tx, uint(1))
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRollback(t *testing.T) {
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repository.NewRepository[TestModel](db)

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
		return repo.CreateTx(tx, testModel)
	})

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
