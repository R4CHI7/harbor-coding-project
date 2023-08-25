package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/harbor-xyz/coding-project/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SlotTestSuite struct {
	suite.Suite
	repo Slot
	mock sqlmock.Sqlmock
}

func (suite *SlotTestSuite) SetupTest() {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		suite.NoError(err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	suite.NoError(err)

	suite.repo = Slot{db: db}
	suite.mock = mock
}

func (suite *SlotTestSuite) TestCreateHappyFlow() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "slots" ("user_id","start_time","end_time","status","created_at","updated_at","deleted_at") 
	VALUES ($1,$2,$3,$4,$5,$6,$7),($8,$9,$10,$11,$12,$13,$14)`)).
		WithArgs(1, sqlmock.AnyArg(), sqlmock.AnyArg(), 0, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			1, sqlmock.AnyArg(), sqlmock.AnyArg(), 0, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1))
	suite.mock.ExpectCommit()

	now := time.Now()
	err := suite.repo.Create(context.Background(), []model.Slot{
		{
			UserID:    1,
			StartTime: now,
			EndTime:   now.Add(30 * time.Minute),
			Status:    model.StatusCreated,
		},
		{
			UserID:    1,
			StartTime: now.Add(30 * time.Minute),
			EndTime:   now.Add(60 * time.Minute),
			Status:    model.StatusCreated,
		},
	})

	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *SlotTestSuite) TestCreateReturnsErrorWhenDBFails() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "slots" ("user_id","start_time","end_time","status","created_at","updated_at","deleted_at") 
	VALUES ($1,$2,$3,$4,$5,$6,$7),($8,$9,$10,$11,$12,$13,$14)`)).
		WithArgs(1, sqlmock.AnyArg(), sqlmock.AnyArg(), 0, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			1, sqlmock.AnyArg(), sqlmock.AnyArg(), 0, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("some error"))
	suite.mock.ExpectRollback()

	now := time.Now()
	err := suite.repo.Create(context.Background(), []model.Slot{
		{
			UserID:    1,
			StartTime: now,
			EndTime:   now.Add(30 * time.Minute),
			Status:    model.StatusCreated,
		},
		{
			UserID:    1,
			StartTime: now.Add(30 * time.Minute),
			EndTime:   now.Add(60 * time.Minute),
			Status:    model.StatusCreated,
		},
	})

	suite.Error(err, "some error")
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *SlotTestSuite) TestGetReturnsDataIfExists() {
	now := time.Now()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "slots" WHERE user_id = $1 AND start_time BETWEEN $2 AND $3`)).
		WithArgs(1, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "user_id", "start_time", "end_time", "status", "created_at", "updated_at", "deleted_at"},
		).AddRow(1, 1, now, now.Add(30*time.Minute), 0, now, now, now).
			AddRow(2, 1, now.Add(30*time.Minute), now.Add(60*time.Minute), 0, now, now, now))

	resp, err := suite.repo.Get(context.Background(), 1, now, now.AddDate(0, 0, 7))
	suite.NoError(err)
	suite.Equal(2, len(resp))
	suite.Equal(1, int(resp[0].ID))
}

func (suite *SlotTestSuite) TestGetReturnsErrorIfDBReturnsError() {
	now := time.Now()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "slots" WHERE user_id = $1 AND start_time BETWEEN $2 AND $3`)).
		WithArgs(1, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("some error"))

	resp, err := suite.repo.Get(context.Background(), 1, now, now.AddDate(0, 0, 7))
	suite.Equal("some error", err.Error())
	suite.Nil(resp)
}

func TestSlotTestSuite(t *testing.T) {
	suite.Run(t, new(SlotTestSuite))
}
