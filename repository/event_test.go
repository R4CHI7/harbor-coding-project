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

type EventTestSuite struct {
	suite.Suite
	repo Event
	mock sqlmock.Sqlmock
}

func (suite *EventTestSuite) SetupTest() {
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

	suite.repo = Event{db: db}
	suite.mock = mock
}

func (suite *EventTestSuite) TestCreateHappyFlow() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "events" ("user_id","slot_id","invitee_email","invitee_name","invitee_notes","start_time","end_time","created_at","updated_at","deleted_at")
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`)).
		WithArgs(1, 1, "test@example.xyz", "test", "test", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1))
	suite.mock.ExpectCommit()

	resp, err := suite.repo.Create(context.Background(), model.Event{UserID: 1, SlotID: 1, InviteeEmail: "test@example.xyz", InviteeName: "test", InviteeNotes: "test"})

	suite.Equal(1, int(resp.ID))
	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *EventTestSuite) TestCreateReturnsErrorWhenDBFails() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "events" ("user_id","slot_id","invitee_email","invitee_name","invitee_notes","start_time","end_time","created_at","updated_at","deleted_at")
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`)).
		WithArgs(1, 1, "test@example.xyz", "test", "test", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("some error"))
	suite.mock.ExpectRollback()

	resp, err := suite.repo.Create(context.Background(), model.Event{UserID: 1, SlotID: 1, InviteeEmail: "test@example.xyz", InviteeName: "test", InviteeNotes: "test"})

	suite.Empty(resp)
	suite.Error(err, "some error")
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *EventTestSuite) TestGetReturnsDataIfExists() {
	now := time.Now()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "events" WHERE user_id = $1`)).
		WithArgs(1).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "user_id", "slot_id", "invitee_email", "invitee_name", "invitee_notes", "start_time", "end_time", "created_at", "updated_at", "deleted_at"},
	).AddRow(1, 1, 1, "test@example.xyz", "test", "test", now, now, now, now, now).
		AddRow(2, 1, 2, "test1@example.xyz", "test", "test", now, now, now, now, now))

	resp, err := suite.repo.Get(context.Background(), 1)
	suite.NoError(err)
	suite.Equal(2, len(resp))
	suite.Equal(1, int(resp[0].ID))
}

func (suite *EventTestSuite) TestGetReturnsErrorIfDBReturnsError() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "events" WHERE user_id = $1`)).
		WithArgs(1).WillReturnError(errors.New("some error"))

	resp, err := suite.repo.Get(context.Background(), 1)
	suite.Equal("some error", err.Error())
	suite.Nil(resp)
}

func TestEventTestSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
