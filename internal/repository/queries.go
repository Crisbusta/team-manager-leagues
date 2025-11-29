package repository

// DDL statements for schema creation
var SchemaStatements = []string{
	`CREATE TABLE IF NOT EXISTS leagues (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL UNIQUE,
        slug TEXT NOT NULL UNIQUE,
        region TEXT NOT NULL,
        created_by TEXT NOT NULL, -- References users(id) logically
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );`,

	// Series: 1:N league -> series
	`CREATE TABLE IF NOT EXISTS series (
        id TEXT PRIMARY KEY,
        league_id TEXT NOT NULL REFERENCES leagues(id) ON DELETE CASCADE,
        name TEXT NOT NULL,
        format TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        UNIQUE(league_id, name)
    );`,
	`CREATE UNIQUE INDEX IF NOT EXISTS series_league_name_lower_uidx ON series (league_id, lower(name));`,

	// Team Registrations: M:N team <-> series
	// Note: REFERENCES teams(id) assumes teams table exists in the same DB.
	// Since we assume shared DB, this is fine.
	`CREATE TABLE IF NOT EXISTS team_registrations (
        id TEXT PRIMARY KEY,
        team_id TEXT NOT NULL, -- References teams(id) logically
        series_id TEXT NOT NULL REFERENCES series(id) ON DELETE CASCADE,
        status TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        UNIQUE(team_id, series_id)
    );`,
}

// DML queries
const (
	// Leagues CRUD
	QInsertLeague     = `INSERT INTO leagues (id, name, slug, region, created_by, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,now(),now())`
	QSelectLeagueByID = `SELECT id, name, slug, region, created_by, created_at, updated_at FROM leagues WHERE id=$1`
	QSelectLeagues    = `SELECT id, name, slug, region, created_by, created_at, updated_at FROM leagues ORDER BY created_at`
	QUpdateLeague     = `UPDATE leagues SET name=$2, slug=$3, region=$4, updated_at=now() WHERE id=$1`
	QDeleteLeague     = `DELETE FROM leagues WHERE id=$1`

	// Series CRUD
	QInsertSeries         = `INSERT INTO series (id, league_id, name, format, created_at, updated_at) VALUES ($1,$2,$3,$4,now(),now())`
	QSelectSeriesByLeague = `SELECT id, league_id, name, format, created_at, updated_at FROM series WHERE league_id=$1 ORDER BY created_at`
	QSelectSeriesByID     = `SELECT id, league_id, name, format, created_at, updated_at FROM series WHERE id=$1`
	QUpdateSeries         = `UPDATE series SET name=$2, format=$3, updated_at=now() WHERE id=$1`
	QDeleteSeries         = `DELETE FROM series WHERE id=$1`

	// Team Registrations
	QInsertTeamRegistration      = `INSERT INTO team_registrations (id, team_id, series_id, status, created_at) VALUES ($1,$2,$3,$4,now())`
	QSelectRegistrationsByTeam   = `SELECT id, team_id, series_id, status, created_at FROM team_registrations WHERE team_id=$1`
	QSelectRegistrationsBySeries = `SELECT id, team_id, series_id, status, created_at FROM team_registrations WHERE series_id=$1`

	// Read-only queries for validation (assuming shared DB)
	QSelectTeamByID        = `SELECT id, club_id, name, format, created_at, updated_at FROM teams WHERE id=$1`
	QOwnerMembershipExists = `SELECT 1 FROM memberships WHERE user_id=$1 AND club_id=$2 AND role='owner' AND status='active' LIMIT 1`
)
