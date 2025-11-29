package http

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"strings"

	"base_app/internal/entity"
	v1 "base_app/internal/handler/http/v1"
	"base_app/internal/usecase"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
)

// Handler is the implementation of the ogen-generated interface.
// It connects the generated server with the application's use cases.
type Handler struct {
	authUsecase    usecase.AuthUsecase
	dataUsecase    usecase.DataUsecase
	catalogUsecase usecase.CatalogUsecase
	sessionManager *scs.SessionManager
	contentFS      fs.FS
}

// NewHandler creates a new handler implementation.
func NewHandler(
	authUC usecase.AuthUsecase,
	dataUC usecase.DataUsecase,
	catalogUC usecase.CatalogUsecase,
	sm *scs.SessionManager,
	contentFS fs.FS,
) *Handler {
	return &Handler{
		authUsecase:    authUC,
		dataUsecase:    dataUC,
		catalogUsecase: catalogUC,
		sessionManager: sm,
		contentFS:      contentFS,
	}
}

// --- Interface Implementation ---

// Login implements login operation.
func (h *Handler) Login(ctx context.Context, req *v1.LoginRequest) (v1.LoginRes, error) {
	user, err := h.authUsecase.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		return &v1.LoginUnauthorized{}, nil // Return specific error type for 401
	}

	if err := h.sessionManager.RenewToken(ctx); err != nil {
		return nil, err // Let the framework handle this as a 500
	}

	h.sessionManager.Put(ctx, "userID", user.ID.String())
	h.sessionManager.Put(ctx, "userEmail", user.Email)

	// Convert entity.User to v1.User
	response := &v1.User{
		ID:        v1.NewOptUUID(user.ID),
		Email:     v1.NewOptString(user.Email),
		CreatedAt: v1.NewOptDateTime(user.CreatedAt),
	}
	return response, nil
}

// Logout implements logout operation.
func (h *Handler) Logout(ctx context.Context) (v1.LogoutRes, error) {
	if err := h.sessionManager.Destroy(ctx); err != nil {
		return nil, err
	}
	return &v1.LogoutOK{}, nil
}

// GetMe implements getMe operation.
func (h *Handler) GetMe(ctx context.Context) (v1.GetMeRes, error) {
	userID := h.sessionManager.GetString(ctx, "userID")
	if userID == "" {
		return &v1.GetMeUnauthorized{}, nil
	}

	userEmail := h.sessionManager.GetString(ctx, "userEmail")
	parsedID, _ := uuid.Parse(userID)

	response := &v1.User{
		ID:    v1.NewOptUUID(parsedID),
		Email: v1.NewOptString(userEmail),
		// CreatedAt is not available in session, so it's omitted
	}
	return response, nil
}

// PostData implements postData operation.
func (h *Handler) PostData(ctx context.Context, req *v1.DataRequest) (v1.PostDataRes, error) {
	data := &entity.Data{
		Key:   req.Key,
		Value: req.Value,
	}
	if err := h.dataUsecase.SaveData(ctx, data); err != nil {
		return nil, err
	}
	return &v1.PostDataCreated{}, nil
}

// GetCatalog implements getCatalog operation.
func (h *Handler) GetCatalog(ctx context.Context) (v1.GetCatalogRes, error) {
	items, err := h.catalogUsecase.GetCatalogItems(ctx)
	if err != nil {
		return nil, err
	}

	// Convert []entity.CatalogItem to v1.GetCatalogOKApplicationJSON
	response := make(v1.GetCatalogOKApplicationJSON, len(items))
	for i, item := range items {
		response[i] = v1.CatalogItem{
			ID:          v1.NewOptUUID(item.ID),
			Title:       v1.NewOptString(item.Title),
			Description: v1.NewOptString(item.Description),
			Disabled:    v1.NewOptBool(item.Disabled),
		}
	}
	return &response, nil
}

// --- Security Handler ---

// HandleCookieAuth implements cookieAuth security scheme.
func (h *Handler) HandleCookieAuth(ctx context.Context, operationName string, t v1.CookieAuth) (context.Context, error) {
	userID := h.sessionManager.GetString(ctx, "userID")
	if userID == "" {
		// Return a standard error that ogen will interpret as a 401 Unauthorized.
		return ctx, errors.New("unauthorized")
	}
	return ctx, nil
}

// --- Static File Server ---

// ServeHTTP serves the frontend files.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if strings.HasSuffix(p, "/") {
		p = strings.TrimSuffix(p, "/") + "/index.html"
	}
	p = strings.TrimPrefix(p, "/")

	f, err := h.contentFS.Open(p)
	if err != nil {
		http.ServeFileFS(w, r, h.contentFS, "index.html")
		return
	}
	_ = f.Close()

	http.FileServer(http.FS(h.contentFS)).ServeHTTP(w, r)
}
