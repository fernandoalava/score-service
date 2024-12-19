package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/fernandoalava/softwareengineer-test-task/domain"
)

type TicketRepository struct {
	Conn *sql.DB
}

func NewTicketRepository(conn *sql.DB) *TicketRepository {
	return &TicketRepository{conn}
}

func (repository *TicketRepository) FetchAll(ctx context.Context) (result []domain.Ticket, err error) {

	query := `SELECT id, subject, created_at FROM tickets`

	rows, err := repository.Conn.QueryContext(ctx, query)
	if err != nil {
		log.Println("error while querying tickets table", err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Println("error trying to close rows", err)
		}
	}()

	result = make([]domain.Ticket, 0)
	for rows.Next() {
		ticket := domain.Ticket{}
		err = rows.Scan(
			&ticket.ID,
			&ticket.Subject,
			&ticket.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		result = append(result, ticket)
	}

	return result, nil

}
