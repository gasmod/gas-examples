package files

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gasmod/gas"
	"github.com/gasmod/gas-examples/api-server/auth"
	"github.com/gasmod/gas-examples/api-server/db"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// FileUploaded is emitted after a file is successfully stored. Other services
// can subscribe to react to new uploads.
var FileUploaded = gas.Event[FileUploadedPayload]{Name: "file:uploaded"}

type FileUploadedPayload struct {
	FileID uuid.UUID
	UserID uuid.UUID
	Name   string
}

type Service struct {
	router  *gas.Router
	bus     *gas.EventBus
	db      gas.DatabaseProvider
	mgr     gas.MigrationManager
	storage gas.StorageProvider
	cache   gas.CacheProvider
	auth    *auth.Service
	queries *db.Queries
}

func New(
	router *gas.Router,
	bus *gas.EventBus,
	dbProvider gas.DatabaseProvider,
	mgr gas.MigrationManager,
	storage gas.StorageProvider,
	cache gas.CacheProvider,
	authSvc *auth.Service,
) *Service {
	return &Service{
		router:  router,
		bus:     bus,
		db:      dbProvider,
		mgr:     mgr,
		storage: storage,
		cache:   cache,
		auth:    authSvc,
		queries: db.New(dbProvider.DB()),
	}
}

func (s *Service) Name() string { return "files" }

func (s *Service) Init() error {
	if err := s.mgr.RegisterFS(s.Name(), migrationsFS); err != nil {
		return err
	}

	// All file routes require authentication.
	s.router.Group(func(sub *gas.Router) {
		sub.UseMiddlewareFunc(s.auth.Middleware())

		sub.Handle(s.Name(), http.MethodPost, "/api/files", s.handleUpload)
		sub.Handle(s.Name(), http.MethodGet, "/api/files", s.handleList)
		sub.Handle(s.Name(), http.MethodGet, "/api/files/{id}", s.handleGet)
		sub.Handle(s.Name(), http.MethodGet, "/api/files/{id}/download", s.handleDownload)
		sub.Handle(s.Name(), http.MethodDelete, "/api/files/{id}", s.handleDelete)
	})

	return nil
}

func (s *Service) Close() error { return nil }

// --- Response types ---

type fileResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type downloadResponse struct {
	URL string `json:"url"`
}

// --- Handlers ---

func (s *Service) handleUpload(ctx gas.Context) error {
	// Parse multipart form — 32MB max memory.
	r := ctx.Request()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return &apiError{Status: http.StatusBadRequest, Message: "invalid multipart form"}
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return &apiError{Status: http.StatusBadRequest, Message: "missing file field"}
	}
	defer file.Close()

	principal := gas.PrincipalFromContext(ctx)
	userID, _ := uuid.Parse(principal.Subject())

	// Generate a unique storage key.
	storageKey := fmt.Sprintf("%s/%s/%s", userID, uuid.New(), header.Filename)

	// Detect content type from the multipart header.
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to S3 with explicit content type so downloads serve the
	// correct MIME type without needing to guess from the key.
	if err := s.storage.Upload(ctx, storageKey, file, gas.WithContentType(contentType)); err != nil {
		return fmt.Errorf("upload to storage: %w", err)
	}

	f, err := s.queries.CreateFile(ctx, db.CreateFileParams{
		UserID:      userID,
		Name:        header.Filename,
		Size:        header.Size,
		ContentType: contentType,
		StorageKey:  storageKey,
	})
	if err != nil {
		return fmt.Errorf("create file record: %w", err)
	}

	// Emit event for other services to react to.
	gas.Emit(s.bus, FileUploaded, FileUploadedPayload{
		FileID: f.ID,
		UserID: userID,
		Name:   f.Name,
	})

	return ctx.JSON(http.StatusCreated, fileResponse{
		ID:          f.ID,
		Name:        f.Name,
		Size:        f.Size,
		ContentType: f.ContentType,
		CreatedAt:   f.CreatedAt,
	})
}

func (s *Service) handleList(ctx gas.Context) error {
	principal := gas.PrincipalFromContext(ctx)
	userID, _ := uuid.Parse(principal.Subject())

	files, err := s.queries.ListFilesByUser(ctx, userID)
	if err != nil {
		return err
	}

	resp := make([]fileResponse, len(files))
	for i, f := range files {
		resp[i] = fileResponse{
			ID:          f.ID,
			Name:        f.Name,
			Size:        f.Size,
			ContentType: f.ContentType,
			CreatedAt:   f.CreatedAt,
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Service) handleGet(ctx gas.Context) error {
	fileID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return &apiError{Status: http.StatusBadRequest, Message: "invalid file id"}
	}

	// Check cache first.
	cacheKey := "file:" + fileID.String()
	if data, err := s.cache.Get(ctx, cacheKey); err == nil {
		var resp fileResponse
		if json.Unmarshal(data, &resp) == nil {
			return ctx.JSON(http.StatusOK, resp)
		}
	}

	f, err := s.queries.GetFileByID(ctx, fileID)
	if err != nil {
		return &apiError{Status: http.StatusNotFound, Message: "file not found"}
	}

	// Verify ownership.
	principal := gas.PrincipalFromContext(ctx)
	userID, _ := uuid.Parse(principal.Subject())
	if f.UserID != userID {
		return &apiError{Status: http.StatusNotFound, Message: "file not found"}
	}

	resp := fileResponse{
		ID:          f.ID,
		Name:        f.Name,
		Size:        f.Size,
		ContentType: f.ContentType,
		CreatedAt:   f.CreatedAt,
	}

	// Cache for 5 minutes.
	if data, err := json.Marshal(resp); err == nil {
		s.cache.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Service) handleDownload(ctx gas.Context) error {
	fileID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return &apiError{Status: http.StatusBadRequest, Message: "invalid file id"}
	}

	// Check cache for presigned URL.
	cacheKey := "download:" + fileID.String()
	if data, err := s.cache.Get(ctx, cacheKey); err == nil {
		return ctx.JSON(http.StatusOK, downloadResponse{URL: string(data)})
	}

	f, err := s.queries.GetFileByID(ctx, fileID)
	if err != nil {
		return &apiError{Status: http.StatusNotFound, Message: "file not found"}
	}

	// Verify ownership.
	principal := gas.PrincipalFromContext(ctx)
	userID, _ := uuid.Parse(principal.Subject())
	if f.UserID != userID {
		return &apiError{Status: http.StatusNotFound, Message: "file not found"}
	}

	// Generate presigned URL valid for 15 minutes.
	url, err := s.storage.PresignDownloadURL(ctx, f.StorageKey, 15*time.Minute)
	if err != nil {
		return fmt.Errorf("presign url: %w", err)
	}

	// Cache the presigned URL for 14 minutes (slightly less than expiry).
	s.cache.Set(ctx, cacheKey, []byte(url), 14*time.Minute)

	return ctx.JSON(http.StatusOK, downloadResponse{URL: url})
}

func (s *Service) handleDelete(ctx gas.Context) error {
	fileID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return &apiError{Status: http.StatusBadRequest, Message: "invalid file id"}
	}

	f, err := s.queries.GetFileByID(ctx, fileID)
	if err != nil {
		return &apiError{Status: http.StatusNotFound, Message: "file not found"}
	}

	principal := gas.PrincipalFromContext(ctx)
	userID, _ := uuid.Parse(principal.Subject())
	if f.UserID != userID {
		return &apiError{Status: http.StatusNotFound, Message: "file not found"}
	}

	// Delete from S3 first, then DB.
	if err := s.storage.Delete(ctx, f.StorageKey); err != nil {
		return fmt.Errorf("delete from storage: %w", err)
	}

	if err := s.queries.DeleteFile(ctx, db.DeleteFileParams{
		ID:     fileID,
		UserID: userID,
	}); err != nil {
		return fmt.Errorf("delete file record: %w", err)
	}

	// Invalidate cache.
	s.cache.Delete(ctx, "file:"+fileID.String())
	s.cache.Delete(ctx, "download:"+fileID.String())

	return ctx.NoContent()
}

type apiError struct {
	Status  int    `json:"-"`
	Message string `json:"error"`
}

func (e *apiError) Error() string      { return e.Message }
func (e *apiError) StatusCode() int    { return e.Status }
