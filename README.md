**Running locally:**
- Put your token in `deploy/dev/txstat.env` (`TXSTAT_GETBLOCKIO_TOKEN` key);
- Run `make devup` to start docker-compose;
- Run `curl GET localhost:8080/stats/greatestBalanceDiff` and wait for response;
- Run `make devdown` to shut down;