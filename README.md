# Mah Moneh

[![test status](https://github.com/d-ashesss/mah-moneh/workflows/test/badge.svg?branch=main)](https://github.com/d-ashesss/mah-moneh/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/d-ashesss/mah-moneh)](https://goreportcard.com/report/github.com/d-ashesss/mah-moneh)
[![Go version](https://img.shields.io/github/go-mod/go-version/d-ashesss/mah-moneh)](https://github.com/d-ashesss/mah-moneh/blob/main/go.mod)
[![latest tag](https://img.shields.io/github/v/tag/d-ashesss/mah-moneh?include_prereleases&sort=semver)](https://github.com/d-ashesss/mah-moneh/tags)
[![MIT license](https://img.shields.io/github/license/d-ashesss/mah-moneh?color=blue)](https://opensource.org/licenses/MIT)
[![feline reference](https://img.shields.io/badge/may%20contain%20cat%20fur-%F0%9F%90%88-blueviolet)](https://github.com/d-ashesss/mah-moneh)

This piece of software is supposed to help to somehow observe and manage your own finances 🤷

## Configuration

The configuration is done via environment variables. The following variables are supported:

* `PORT` - the port to run the web server on, default: 8080
* `DB_HOST` - the hostname of the database server, default: localhost
* `DB_PORT` - the port of the database server, default: 5432
* `DB_NAME` - the name of the database, default: postgres
* `DB_USER` - database user, default: postgres
* `DB_PASSWORD` - database user's password, required
* `DB_DEBUG` - whether to use the database in debug mode
* `AUTH_OPENID_CONFIGURATION_URL` - OpenID configuration URL for OIDC provider.
