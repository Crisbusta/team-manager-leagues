package repository

import (
	"context"
	"errors"

	"team-manager-leagues/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store { return &Store{Pool: pool} }

// Leagues
func (s *Store) CreateLeague(ctx context.Context, l *domain.League) error {
	_, err := s.Pool.Exec(ctx, QInsertLeague, l.ID, l.Name, l.Slug, l.Region, l.CreatedBy)
	return err
}
func (s *Store) GetLeagueByID(ctx context.Context, id string) (*domain.League, error) {
	row := s.Pool.QueryRow(ctx, QSelectLeagueByID, id)
	var l domain.League
	if err := row.Scan(&l.ID, &l.Name, &l.Slug, &l.Region, &l.CreatedBy, &l.CreatedAt, &l.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &l, nil
}
func (s *Store) ListLeagues(ctx context.Context) ([]domain.League, error) {
	rows, err := s.Pool.Query(ctx, QSelectLeagues)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.League{}
	for rows.Next() {
		var l domain.League
		if err := rows.Scan(&l.ID, &l.Name, &l.Slug, &l.Region, &l.CreatedBy, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}
func (s *Store) UpdateLeague(ctx context.Context, id, name, slug, region string) error {
	_, err := s.Pool.Exec(ctx, QUpdateLeague, id, name, slug, region)
	return err
}
func (s *Store) DeleteLeague(ctx context.Context, id string) error {
	_, err := s.Pool.Exec(ctx, QDeleteLeague, id)
	return err
}

// Series
func (s *Store) CreateSeries(ctx context.Context, ser *domain.Series) error {
	_, err := s.Pool.Exec(ctx, QInsertSeries, ser.ID, ser.LeagueID, ser.Name, ser.Format)
	return err
}
func (s *Store) ListSeriesByLeague(ctx context.Context, leagueID string) ([]domain.Series, error) {
	rows, err := s.Pool.Query(ctx, QSelectSeriesByLeague, leagueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Series{}
	for rows.Next() {
		var ser domain.Series
		if err := rows.Scan(&ser.ID, &ser.LeagueID, &ser.Name, &ser.Format, &ser.CreatedAt, &ser.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, ser)
	}
	return out, rows.Err()
}
func (s *Store) GetSeriesByID(ctx context.Context, id string) (*domain.Series, error) {
	row := s.Pool.QueryRow(ctx, QSelectSeriesByID, id)
	var ser domain.Series
	if err := row.Scan(&ser.ID, &ser.LeagueID, &ser.Name, &ser.Format, &ser.CreatedAt, &ser.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ser, nil
}
func (s *Store) UpdateSeries(ctx context.Context, id, name, format string) error {
	_, err := s.Pool.Exec(ctx, QUpdateSeries, id, name, format)
	return err
}
func (s *Store) DeleteSeries(ctx context.Context, id string) error {
	_, err := s.Pool.Exec(ctx, QDeleteSeries, id)
	return err
}

// Team Registrations
func (s *Store) CreateTeamRegistration(ctx context.Context, tr *domain.TeamRegistration) error {
	_, err := s.Pool.Exec(ctx, QInsertTeamRegistration, tr.ID, tr.TeamID, tr.SeriesID, tr.Status)
	return err
}
func (s *Store) ListRegistrationsByTeam(ctx context.Context, teamID string) ([]domain.TeamRegistration, error) {
	rows, err := s.Pool.Query(ctx, QSelectRegistrationsByTeam, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.TeamRegistration{}
	for rows.Next() {
		var tr domain.TeamRegistration
		if err := rows.Scan(&tr.ID, &tr.TeamID, &tr.SeriesID, &tr.Status, &tr.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, tr)
	}
	return out, rows.Err()
}
func (s *Store) ListRegistrationsBySeries(ctx context.Context, seriesID string) ([]domain.TeamRegistration, error) {
	rows, err := s.Pool.Query(ctx, QSelectRegistrationsBySeries, seriesID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.TeamRegistration{}
	for rows.Next() {
		var tr domain.TeamRegistration
		if err := rows.Scan(&tr.ID, &tr.TeamID, &tr.SeriesID, &tr.Status, &tr.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, tr)
	}
	return out, rows.Err()
}

func (s *Store) UpdateRegistrationStatus(ctx context.Context, id, status string) error {
	_, err := s.Pool.Exec(ctx, QUpdateRegistrationStatus, id, status)
	return err
}

func (s *Store) DeleteRegistration(ctx context.Context, id string) error {
	_, err := s.Pool.Exec(ctx, QDeleteRegistration, id)
	return err
}

// Read-only helpers
func (s *Store) GetTeamByID(ctx context.Context, id string) (*domain.Team, error) {
	row := s.Pool.QueryRow(ctx, QSelectTeamByID, id)
	var t domain.Team
	if err := row.Scan(&t.ID, &t.ClubID, &t.Name, &t.Format, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (s *Store) IsOwner(ctx context.Context, userID, clubID string) (bool, error) {
	row := s.Pool.QueryRow(ctx, QOwnerMembershipExists, userID, clubID)
	var one int
	if err := row.Scan(&one); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
