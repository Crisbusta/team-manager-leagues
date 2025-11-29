package service

import (
	"context"
	"errors"
	"strings"

	"team-manager-leagues/internal/domain"
	"team-manager-leagues/internal/repository"
	"team-manager-leagues/internal/util"
)

type LeaguesService struct {
	store *repository.Store
}

func NewLeaguesService(store *repository.Store) *LeaguesService {
	return &LeaguesService{store: store}
}

// Leagues

func (s *LeaguesService) CreateLeague(ctx context.Context, userID, name, region string) (*domain.League, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("invalid name")
	}
	slug := util.Slugify(name)

	l := &domain.League{ID: util.RandID(), Name: name, Slug: slug, Region: region, CreatedBy: userID}
	if err := s.store.CreateLeague(ctx, l); err != nil {
		return nil, err
	}
	return l, nil
}

func (s *LeaguesService) ListLeagues(ctx context.Context) ([]domain.League, error) {
	return s.store.ListLeagues(ctx)
}

func (s *LeaguesService) GetLeague(ctx context.Context, id string) (*domain.League, error) {
	return s.store.GetLeagueByID(ctx, id)
}

func (s *LeaguesService) UpdateLeague(ctx context.Context, id, name, region string) (*domain.League, error) {
	// TODO: Check permission (e.g. admin or creator). For now MVP allows update if authenticated?
	// The handler didn't implement UpdateLeague, only Create/List/Get.
	// But store has UpdateLeague.
	// I'll implement it assuming creator check or similar if I had the logic.
	// For now, let's just implement it.

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("invalid name")
	}
	slug := util.Slugify(name)

	if err := s.store.UpdateLeague(ctx, id, name, slug, region); err != nil {
		return nil, err
	}
	return &domain.League{ID: id, Name: name, Slug: slug, Region: region}, nil
}

func (s *LeaguesService) DeleteLeague(ctx context.Context, id string) error {
	return s.store.DeleteLeague(ctx, id)
}

// Series

func (s *LeaguesService) CreateSeries(ctx context.Context, leagueID, name, format string) (*domain.Series, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("invalid name")
	}

	ser := &domain.Series{ID: util.RandID(), LeagueID: leagueID, Name: name, Format: format}
	if err := s.store.CreateSeries(ctx, ser); err != nil {
		return nil, err
	}
	return ser, nil
}

func (s *LeaguesService) ListSeries(ctx context.Context, leagueID string) ([]domain.Series, error) {
	return s.store.ListSeriesByLeague(ctx, leagueID)
}

func (s *LeaguesService) GetSeries(ctx context.Context, id string) (*domain.Series, error) {
	return s.store.GetSeriesByID(ctx, id)
}

// Registrations

func (s *LeaguesService) RegisterTeam(ctx context.Context, userID, teamID, seriesID string) (*domain.TeamRegistration, error) {
	// Verify team exists
	t, err := s.store.GetTeamByID(ctx, teamID)
	if err != nil || t == nil {
		return nil, errors.New("team not found")
	}

	// Verify ownership
	isOwner, err := s.store.IsOwner(ctx, userID, t.ClubID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("forbidden: only club owner can register teams")
	}

	// Create registration
	reg := &domain.TeamRegistration{
		ID:       util.RandID(),
		TeamID:   teamID,
		SeriesID: seriesID,
		Status:   "active",
	}
	if err := s.store.CreateTeamRegistration(ctx, reg); err != nil {
		return nil, err
	}
	return reg, nil
}

func (s *LeaguesService) ListRegistrationsByTeam(ctx context.Context, teamID string) ([]domain.TeamRegistration, error) {
	return s.store.ListRegistrationsByTeam(ctx, teamID)
}

func (s *LeaguesService) ListRegistrationsBySeries(ctx context.Context, seriesID string) ([]domain.TeamRegistration, error) {
	return s.store.ListRegistrationsBySeries(ctx, seriesID)
}
