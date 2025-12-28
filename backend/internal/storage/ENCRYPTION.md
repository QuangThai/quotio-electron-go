# Encryption at Rest

Sensitive fields in the Account model are encrypted at rest using AES-256-GCM.

## Encrypted Fields

The following Account fields are automatically encrypted before storage and decrypted after retrieval:
- `APIKey` - Provider API key
- `OAuthToken` - OAuth access token
- `RefreshToken` - OAuth refresh token

## How It Works

**Encryption**: When an Account is saved (via `BeforeSave` hook), sensitive fields are encrypted using AES-256-GCM and stored as base64-encoded strings.

**Decryption**: When an Account is loaded from the database (via `AfterFind` hook), encrypted fields are automatically decrypted back to plaintext in memory.

**Backwards Compatibility**: If decryption fails (e.g., for legacy plaintext data), the value is returned as-is. New data will be encrypted.

## Encryption Key Management

### Key Persistence Strategy

The application uses a three-tier approach to manage encryption keys:

1. **Environment Variable** (highest priority): If `QUOTIO_ENCRYPTION_KEY` is set, it will be used
2. **File-based Storage** (recommended for desktop): Key is persisted to `~/.quotio/.encryption.key`
3. **Auto-generation** (fallback): If neither above exists, a new key is generated and persisted to file

### Desktop Application (Default)

By default, a random 256-bit key is generated on first startup and persisted to:
```
~/.quotio/.encryption.key
```

This file is protected with restrictive permissions (0600). The key will be reused on subsequent restarts, ensuring encrypted data remains accessible.

### Using Environment Variables (Optional)

To override the file-based key or use a custom key:
```bash
# Generate a key
openssl rand -base64 32

# Set in environment
export QUOTIO_ENCRYPTION_KEY="your-base64-encoded-key"
```

**Important**: Store the encryption key securely. Without it, encrypted data cannot be recovered.

## Implementation Details

- **Algorithm**: AES-256-GCM (authenticated encryption with associated data)
- **Key Size**: 256-bit (32 bytes)
- **Nonce**: 96-bit (12 bytes), randomly generated per encryption
- **Encoding**: Base64 for safe storage in text fields

## Troubleshooting

- **"encryption key not initialized"**: The `InitEncryption()` function was not called. This happens automatically during `storage.Initialize()`.
- **Decryption errors on legacy data**: Plaintext values from before encryption was enabled are automatically returned as-is.
