# Media Module Migration Manifest

This document records the inventory of all symbols, contracts, adapters, delivery endpoints, consumers, and migration risks for the `Media` business module.

---

## 1. Owned Symbols
- **Entities**: `MediaAsset` (struct)
- **Validation Rules**: `MaxMediaUploadBytes` (const = 5MB), `AllowedMediaMimeTypes` (map), `IsAllowedMediaMimeType(mimeType string)` (func)

---

## 2. Owned Ports
- **Persistence Port**: `MediaRepo`
  - `Store(ctx context.Context, media *MediaAsset) error`
  - `List(ctx context.Context, page pagination.Pagination, query string) ([]MediaAsset, int, error)`
  - `Get(ctx context.Context, id string) (MediaAsset, error)`
  - `Delete(ctx context.Context, id string) error`
- **Capability Port**: `MediaStorageProvider` (renamed from `usecase.MediaStorage`)
  - `GeneratePresignedURL(ctx context.Context, fileName string, contentType string) (uploadURL string, publicURL string, fileID string, err error)`
  - `Delete(ctx context.Context, key string) error`

---

## 3. Application Services
- **Interface**: `MediaUseCase`
  - `Store(ctx context.Context, media MediaAsset) (MediaAsset, error)`
  - `List(ctx context.Context, page pagination.Pagination, query string) ([]MediaAsset, int, error)`
  - `Delete(ctx context.Context, id string) error`
  - `GeneratePresignedURL(ctx context.Context, fileName string, contentType string) (uploadURL string, publicURL string, fileID string, err error)`
- **Implementation Struct**: `mediaUseCase`
- **Config Struct**: `Config` (`MaxBytes int64`)
- **Constructor**: `NewMediaUseCase(r MediaRepo, s MediaStorageProvider, cfg Config) MediaUseCase`

---

## 4. Infrastructure Adapters
- `internal/repo/persistent/media_postgres.go` (`MediaRepo` struct implementing `MediaRepo`)
- `internal/repo/persistent/media_postgres_test.go`
- `internal/repo/webapi/r2_storage.go` (`R2Storage` implementing `MediaStorageProvider`)
- `internal/repo/webapi/local_storage.go` (`LocalStorage` implementing `MediaStorageProvider`)

---

## 5. Delivery Endpoints & Mappers
- `internal/controller/restapi/v1/media/handler.go`
- `internal/controller/restapi/v1/media/mapper.go`
- `internal/controller/restapi/v1/lead/error_mapper.go`
- `internal/controller/restapi/v1/server.go`
- `internal/controller/restapi/router.go`

---

## 6. App Wiring
- `internal/app/app.go` (`useCases.media *media.UseCase` → `useCases.media media.MediaUseCase`)

---

## 7. Unit & Integration Tests
- `internal/usecase/media/media_test.go` → move to `internal/module/media/media_usecase_test.go`

---

## 8. Cross-Module Dependencies & Risks
- Shared pagination primitive (`internal/shared/pagination`) used for `List`.
- Risk: Low.
