# Webhook proxy

A service for calling webhooks

# Usage

## Start the proxy server

```
go run main.go
```

## Endpoints and Methods

- `/`
  - `GET` - Fetch Server Info
- `/proxy`
  - `GET` - Fetch the Route info
  - `POST` - Post Payload info to call desired webhooks

# Proxy Call Payload

## Example body for POST request at `/proxy`

```json
{
  "url": "http://localhost:3000",
  "payload": {
    "data_1": "Chat 1",
    "data_2": {
      "data_1": "Chat 1",
      "data_2": "Chat 2"
    }
  },
  "headers": {
    "header_1": "first",
    "header_2": "second"
  }
}
```

## Data structure for post body with go types

- Url

  - Type - `string`
  - Reason - It will be a http end point i.e. a single string of URL

- Payload

  - Type - `json.RawMessage`
  - Reason - It is a JSON data and its structure is not necessary to be parsed as it will be directly copied in bytes to Webhook call

- Headers
  - Type - `map[string]string`
  - Reason - We need to add headers to post request so have parse it in key and value pair

## Persistance Storage

- Currently all info from log will be piped to `logfile.txt` as storage.
- It will contain all the status logs created during webook request and in calling it.

## Retry of Retrival errors

If we get webhook response with following status `502`, `503`, `504`, function will call try to call endpoint again uptill we get any other status, either success or faliure.
