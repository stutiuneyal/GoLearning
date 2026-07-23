package repository

import (
	"context"
	"database/sql"

	"example.com/learning/gin/db/query"
	"example.com/learning/gin/models"
)

var _ EventRepository = (*EventRepositoryImpl)(nil)

type EventRepository interface {
	Save(event *models.Event) error
	GetAllEvents(userId int) ([]models.Event, error)
	GetEventById(id, userId int) (models.Event, error)
	Update(event models.Event) error
	Delete(id, userId int) error
}

type EventRepositoryImpl struct {
	db *sql.DB
}

func NewEventRepositoryImpl(db *sql.DB) *EventRepositoryImpl {
	return &EventRepositoryImpl{
		db: db,
	}
}

func (er *EventRepositoryImpl) Save(event *models.Event) error {

	ctx := context.Background()

	tx, err := er.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query.InsertEventQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// insert the row and update the event.Id from the return value
	if err := stmt.QueryRow(event.Name, event.Description, event.Location, event.Datetime, event.UserId).Scan(&event.Id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (er *EventRepositoryImpl) GetAllEvents(userId int) ([]models.Event, error) {

	var events []models.Event

	// We don't require a statement here as it will not be re-used
	rows, err := er.db.Query(query.GetAllEventsQuery, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var event models.Event

		if err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Datetime, &event.UserId); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil

}

func (er *EventRepositoryImpl) GetEventById(id, userId int) (models.Event, error) {

	var event models.Event

	// since it is a get query we dont need to prepare a statement
	if err := er.db.QueryRow(query.GetEventByIdQuery, id, userId).Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Datetime, &event.UserId); err != nil {
		return models.Event{}, err
	}

	return event, nil
}

func (er *EventRepositoryImpl) Update(event models.Event) error {

	ctx := context.Background()

	tx, err := er.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query.UpdateEventQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(event.Id, event.UserId, event.Name, event.Description, event.Location, event.Datetime); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (er *EventRepositoryImpl) Delete(id, userId int) error {

	ctx := context.Background()

	tx, err := er.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query.DeleteEventQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(id, userId); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
