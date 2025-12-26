# API

Base URL (local): `http://localhost:8080`

## Health

`GET /health`

Response:

```json
{ "status": "ok" }
```

## Ingest a sequence

`POST /v1/sequences/ingest`

Body:

```json
{ "values": [1, -2, 3, 0] }
```

Response `201`:

```json
{
  "sequence": {
    "id": "uuid",
    "created_at": "2025-01-01T00:00:00Z",
    "values": [1, -2, 3, 0],
    "count": 4,
    "sum_fourth_powers_non_positive": 16,
    "min": -2,
    "max": 3,
    "processed": false
  },
  "queued": true
}
```

## Get a sequence

`GET /v1/sequences/{id}`

Response `200`:

```json
{ "sequence": { ... } }
```
