## Quickstart

```bash
docker-compose up
```

```bash
curl -X POST --header "X-Api-Key: secret" -H "Content-Type: application/json" -d @example_long_link.json  http://localhost/rest/v3/short-urls
```

Then

```bash
curl -X GET --header "X-Api-Key: secret" -H "Content-Type: application/json" -d @example_long_link.json  http://localhost/rest/v3/short-urls/<shortcode>
```
