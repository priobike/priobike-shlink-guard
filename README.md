# priobike-shlink-guard

This is a proxy that can be used in front of a Shlink service to make sure that only valid requests and espacially valid long links are sent to the Shlink. If a link is valid, it gets proxied to the Shlink. Otherwise, an error is returned.

Valid means in our PrioBike-context that the long link represents a shortcut-sharing-link consisting of a base64 encoded json shortcut object.

This service is only useful when used in combination with a Shlink service. It is not a standalone service.

[Learn more about PrioBike](https://github.com/priobike)

## Quickstart

Run locally using Docker:
```bash
docker-compose up
```

## CLI

Run tests:
```bash
go test
```

## API

This service mirrors the [Shlink API](https://api-spec.shlink.io/). It is a subset of the Shlink API. The following endpoints are supported and get validated by this service:

- POST /rest/v3/short-urls - Create a new short URL
- GET /rest/v3/short-urls/{shortCode} - Get a short URL by its short code

Every other endpoint is not supported and will be invalid by default.

### Example requests 

#### POST

Should not work:
```bash
curl -X POST --header "X-Api-Key: secret" -H "Content-Type: application/json" -d @example_long_link_base64_invalid.json  http://localhost/rest/v3/short-urls
```

Should work:
```bash
curl -X POST --header "X-Api-Key: secret" -H "Content-Type: application/json" -d @example_long_link_shortcut_location.json  http://localhost/rest/v3/short-urls
```

#### GET

Should not work:
```bash
curl -X GET --header "X-Api-Key: secret" -H "Content-Type: application/json" http://localhost/rest/v3/short-urls
```

Should work (if short link exists):
```bash
curl -X GET --header "X-Api-Key: secret" -H "Content-Type: application/json" http://localhost/rest/v3/short-urls/segrs4
```

## Contributing

We highly encourage you to open an issue or a pull request. You can also use our repository freely with the `MIT` license. 

Every service runs through testing before it is deployed in our release setup. Read more in our [PrioBike deployment readme](https://github.com/priobike/.github/blob/main/wiki/deployment.md) to understand how specific branches/tags are deployed.

## Anything unclear?

Help us improve this documentation. If you have any problems or unclarities, feel free to open an issue.