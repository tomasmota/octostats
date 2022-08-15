This repository contains a cli for fetching octoups deploy stats

Currently supported actions:
- octostats deployments

TO-DO:
- ability to select more than one projects and projectgroups
- move --project and --projectgroup flags to a higher level cmd when more stats are available
- since deployments are paged, we need to check the next page until we get an item older than the lookback period