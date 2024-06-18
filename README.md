# PrioBike Shlink Guard

This is a proxy that can be used in front of the Shlink to make sure that only valid requests and espacially valid long links are sent to the Shlink.

## Quickstart

Run locally:
```bash
docker-compose up
```

Run tests:
```bash
go test
```

### POST 

Should not work:
```bash
curl -X POST --header "X-Api-Key: secret" -H "Content-Type: application/json" -d @example_long_link_base64_invalid.json  http://localhost/rest/v3/short-urls
```

Should work:
```bash
curl -X POST --header "X-Api-Key: secret" -H "Content-Type: application/json" -d @example_long_link_shortcut_location.json  http://localhost/rest/v3/short-urls
```

### GET

Should not work:
```bash
curl -X GET --header "X-Api-Key: secret" -H "Content-Type: application/json" http://localhost/rest/v3/short-urls
```

Should work (if short link exists):
```bash
curl -X GET --header "X-Api-Key: secret" -H "Content-Type: application/json" http://localhost/rest/v3/short-urls/segrs4
```