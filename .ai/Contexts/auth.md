# Authentication System

## Overview

Listenarr uses API key authentication for the MVP. This provides a simple, secure way to protect the API while keeping implementation straightforward.

## Implementation

### API Key Authentication

- **Method**: API key in HTTP header or query parameter
- **Header**: `X-API-Key: <api-key>`
- **Query Parameter**: `?apikey=<api-key>`
- **Generation**: Automatically generated on first run if not set
- **Storage**: Stored in config file (`config/config.yml`)

### Security Features

- API key is cryptographically secure (32 bytes, base64 encoded)
- Minimum 16 characters (format validation)
- Health check endpoint (`/api/health`) is publicly accessible
- All other endpoints require valid API key

### Configuration

```yaml
auth:
  enabled: true
  api_key: ""  # Auto-generated if empty
```

### Usage

#### From Frontend (Axios)
```typescript
axios.get('/api/v1/library', {
  headers: {
    'X-API-Key': apiKey
  }
})
```

#### From Command Line (curl)
```bash
curl -H "X-API-Key: your-api-key" http://localhost:8686/api/v1/library
```

Or with query parameter:
```bash
curl "http://localhost:8686/api/v1/library?apikey=your-api-key"
```

### Future Enhancements

- OAuth 2.0 support
- Multiple API keys per user
- API key rotation
- Rate limiting per API key
- API key expiration

## Endpoints

### Public Endpoints
- `GET /api/health` - Health check (no auth required)

### Protected Endpoints
All `/api/v1/*` endpoints require valid API key.

## Error Responses

Invalid or missing API key returns:
```json
{
  "success": false,
  "error": "Invalid or missing API key"
}
```
Status: `401 Unauthorized`

