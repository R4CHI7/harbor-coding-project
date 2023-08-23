package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/harbor-xyz/coding-project/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserAvailabilityTestSuite struct {
	suite.Suite
	repo UserAvailability
	mock sqlmock.Sqlmock
}

func (suite *UserAvailabilityTestSuite) SetupTest() {
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

	suite.repo = UserAvailability{db: db}
	suite.mock = mock
}

func (suite *UserAvailabilityTestSuite) TestCreateHappyFlow() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_availabilities" ("user_id","availability","meeting_duration_mins","created_at","updated_at")
	VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(1, `[{"day":"monday","start_time":"10:00:00","end_time":"17:00:00"}]`, 30, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1))
	suite.mock.ExpectCommit()

	resp, err := suite.repo.Set(context.Background(), model.UserAvailability{
		UserID: 1,
		Availability: datatypes.JSONSlice[model.DayAvailability]{
			model.DayAvailability{
				Day:       "monday",
				StartTime: datatypes.NewTime(10, 0, 0, 0),
				EndTime:   datatypes.NewTime(17, 0, 0, 0),
			},
		},
		MeetingDurationMins: 30,
	})

	suite.Equal(int(resp.UserID), 1)
	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *UserAvailabilityTestSuite) TestCreateReturnsErrorWhenDBFails() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_availabilities" ("user_id","availability","meeting_duration_mins","created_at","updated_at")
	VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(1, `[{"day":"monday","start_time":"10:00:00","end_time":"17:00:00"}]`, 30, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("some error"))
	suite.mock.ExpectRollback()

	resp, err := suite.repo.Set(context.Background(), model.UserAvailability{
		UserID: 1,
		Availability: datatypes.JSONSlice[model.DayAvailability]{
			model.DayAvailability{
				Day:       "monday",
				StartTime: datatypes.NewTime(10, 0, 0, 0),
				EndTime:   datatypes.NewTime(17, 0, 0, 0),
			},
		},
		MeetingDurationMins: 30,
	})

	suite.Empty(resp)
	suite.Error(err, "some error")
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func TestUserAvailabilityTestSuite(t *testing.T) {
	suite.Run(t, new(UserAvailabilityTestSuite))
}
