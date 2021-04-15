package users

import (
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/http/useractivity"
	"github.com/portainer/portainer/api/internal/authorization"
	consts "github.com/portainer/portainer/api/useractivity"
)

type adminInitPayload struct {
	Username string
	Password string
}

func (payload *adminInitPayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Username) || govalidator.Contains(payload.Username, " ") {
		return errors.New("Invalid username. Must not contain any whitespace")
	}
	if govalidator.IsNull(payload.Password) {
		return errors.New("Invalid password")
	}
	return nil
}

// POST request on /api/users/admin/init
func (handler *Handler) adminInit(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload adminInitPayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	users, err := handler.DataStore.User().UsersByRole(portainer.AdministratorRole)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve users from the database", err}
	}

	if len(users) != 0 {
		return &httperror.HandlerError{http.StatusConflict, "Unable to create administrator user", errAdminAlreadyInitialized}
	}

	user := &portainer.User{
		Username:                payload.Username,
		Role:                    portainer.AdministratorRole,
		PortainerAuthorizations: authorization.DefaultPortainerAuthorizations(),
	}

	user.Password, err = handler.CryptoService.Hash(payload.Password)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to hash user password", errCryptoHashFailure}
	}

	err = handler.DataStore.User().CreateUser(user)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist user inside the database", err}
	}

	payload.Password = consts.RedactedValue
	useractivity.LogHttpActivity(handler.UserActivityStore, "", r, payload)

	return response.JSON(w, user)
}
