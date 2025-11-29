package domain

import "time"

type League struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Region    string    `json:"region"` // e.g., "Santiago", "North"
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Series struct {
	ID        string    `json:"id"`
	LeagueID  string    `json:"leagueId"`
	Name      string    `json:"name"`   // e.g., "Series A", "Golden"
	Format    string    `json:"format"` // e.g., "baby", "7", "11" - should match Team format ideally
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TeamRegistration struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamId"`
	SeriesID  string    `json:"seriesId"`
	Status    string    `json:"status"` // e.g., "active", "pending", "archived"
	CreatedAt time.Time `json:"createdAt"`
}

// Read-only models for validation
type Team struct {
	ID        string    `json:"id"`
	ClubID    string    `json:"clubId"`
	Name      string    `json:"name"`
	Format    string    `json:"format"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Membership struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	ClubID    string    `json:"clubId"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}
