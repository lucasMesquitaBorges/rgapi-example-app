## Proxy API
This app is proxy for developer API with testing and learning purposes.

### How to use?
Any request made to `localhost:${APP_PORT}/*` will be proxied to RIOT Developer api (the region is defined at the environment of this proxy `RIOT_API_BASE_URI`).

Only HTTP Methods and request body will be proxied at request and response.

## Hot Reload
The dev docker container uses a tool called [CompileDaemon](https://github.com/githubnemo/CompileDaemon) to hot reload any changes made to the code.

## Setup
To setup this app is needed `docker`, `docker-compose` and `make`.

Steps:
1. Create `.env` by running `cp .env.example .env`
2. Set env var `RIOT_API_TOKEN` using your api credential
3. Run `make up`

## Available commands
- Start server `make up`
- Shutdown server `make down`
- Rebuild docker image `make build`
- Watch app logs `make logs`
- Enter container bash `make bash`
