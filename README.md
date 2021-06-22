# Water API

## Local Development

1. Type `docker compose up` to bring `db`,`api`,`flyway`,`pgadmin` containers online.

2. `tiger-data` adds TIGER state and county boundary datasets to the local development database (see https://postgis.net/docs/Extras.html).  Type `docker compse up tiger-data` after other containers are already running (step 1 above).
