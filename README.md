# Mah Moneh

[![test status](https://github.com/d-ashesss/mah-moneh/workflows/test/badge.svg?branch=main)](https://github.com/d-ashesss/mah-moneh/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/d-ashesss/mah-moneh)](https://goreportcard.com/report/github.com/d-ashesss/mah-moneh)
[![Go version](https://img.shields.io/github/go-mod/go-version/d-ashesss/mah-moneh)](https://github.com/d-ashesss/mah-moneh/blob/main/go.mod)
[![latest tag](https://img.shields.io/github/v/tag/d-ashesss/mah-moneh?include_prereleases&sort=semver)](https://github.com/d-ashesss/mah-moneh/tags)
[![MIT license](https://img.shields.io/github/license/d-ashesss/mah-moneh?color=blue)](https://opensource.org/licenses/MIT)
[![feline reference](https://img.shields.io/badge/may%20contain%20cat%20fur-%F0%9F%90%88-blueviolet)](https://github.com/d-ashesss/mah-moneh)

Personal finance management API.

## Running the app

The easiest way to run the app is to run it in Docker (provide environment variables as needed):

```bash
docker run -p 8080:8080 \
  -e 'DB_HOST=postgres' \
  -e 'DB_PASSWORD=postgres-password' \
  -e 'AUTH_OPENID_CONFIGURATION_URL=https://accounts.google.com/.well-known/openid-configuration' \
  ashesss/mah-moneh:latest
```

Or even with docker compose, where you can configure the environment in `docker-compose.override.yml` (see [example](https://github.com/d-ashesss/mah-moneh/wiki/example-docker%E2%80%90compose.override.yml)), then simply run:
```bash
docker-compose up
```

## Configuration

The configuration is done via environment variables.

* `PORT` - the port to run the web server on, default: 8080

### Authentication

The API handles authentication using JWT tokens. The token is passed in the `Authorization` header as a bearer token. It is possible to use OIDC service like [Auth0](https://auth0.com) to obtain the token and then use it with this API. It is required to provide the app with the URL to the OIDC discovery documents in the `AUTH_OPENID_CONFIGURATION_URL` environment variable.

### Database

The API uses PostgreSQL database. The connection must be with following environment variables:

* `DB_HOST` - the hostname of the database server, default: localhost
* `DB_PORT` - the port of the database server, default: 5432
* `DB_NAME` - the name of the database, default: postgres
* `DB_USER` - database user, default: postgres
* `DB_PASSWORD` - database user's password, required
* `DB_DEBUG` - whether to use the database in debug mode

### CORS

To configure CORS to allow access from a specific domain, set the `CORS_ALLOWED_ORIGINS` environment variable to semicolon-separated list of allowed URLs,
for example `CORS_ALLOWED_ORIGINS=http://localhost:5000;http://example.com`.
