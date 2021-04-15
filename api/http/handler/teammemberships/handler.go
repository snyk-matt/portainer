package teammemberships

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/http/security"
	"github.com/portainer/portainer/api/internal/authorization"
)

const (
	handlerActivityContext = "Portainer"
)

// Handler is the HTTP handler used to handle team membership operations.
type Handler struct {
	*mux.Router
	AuthorizationService *authorization.Service
	DataStore            portainer.DataStore
	UserActivityStore    portainer.UserActivityStore
}

// NewHandler creates a handler to manage team membership operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/team_memberships",
		bouncer.AdminAccess(httperror.LoggerHandler(h.teamMembershipCreate))).Methods(http.MethodPost)
	h.Handle("/team_memberships",
		bouncer.AdminAccess(httperror.LoggerHandler(h.teamMembershipList))).Methods(http.MethodGet)
	h.Handle("/team_memberships/{id}",
		bouncer.AdminAccess(httperror.LoggerHandler(h.teamMembershipUpdate))).Methods(http.MethodPut)
	h.Handle("/team_memberships/{id}",
		bouncer.AdminAccess(httperror.LoggerHandler(h.teamMembershipDelete))).Methods(http.MethodDelete)

	return h
}
