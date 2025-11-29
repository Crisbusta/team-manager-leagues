# Team Manager Leagues Service

Microservice responsible for managing Leagues, Series, and Team Registrations.

- **Language**: Go 1.25
- **Framework**: Gin
- **Database**: PostgreSQL (shared)

## Endpoints

### Leagues
- `POST /leagues` - Create league
- `GET /leagues` - List leagues
- `GET /leagues/:id` - Get league details
- `PUT /leagues/:id` - Update league
- `DELETE /leagues/:id` - Delete league

### Series
- `GET /leagues/:id/series` - List series in a league
- `POST /leagues/:id/series` - Create series
- `PUT /leagues/:id/series/:seriesId` - Update series
- `DELETE /leagues/:id/series/:seriesId` - Delete series

### Registrations
- `POST /registrations` - Register team to series
- `GET /registrations` - List registrations (by team or series)
- `PUT /registrations/:id` - Update registration status
- `DELETE /registrations/:id` - Cancel registration

## Environment Variables

- `PORT` (default `8080`)
- `DATABASE_URL` (required)
- `JWT_SECRET` (required)

## Docker

Build and run:
```bash
docker build -t team-manager-leagues .
docker run -p 8083:8080 --env-file .env team-manager-leagues
```
