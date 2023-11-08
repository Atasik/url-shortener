package v1

import (
	"encoding/json"
	"io"
	"link-shortener/internal/domain"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) initLinkRoutes(api *mux.Router) {
	r := api.PathPrefix("/link").Subrouter()
	r.Methods("POST").HandlerFunc(h.createToken)
	r.HandleFunc("/{token}", h.getOriginalURL).Methods("GET")
}

// @Summary Create Token
// @Tags link
// @ID create-token
// @Accept json
// @Product json
// @Param input body domain.CreateTokenRequest true "link"
// @Success	201		    {object}	tokenResponse     "token"
// @Failure	400,404		{object}	errResponse
// @Failure	500			{object}	errResponse
// @Failure	default		{object}	errResponse
// @Router		/api/v1/link [post]
func (h *Handler) createToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	if r.Header.Get("Content-Type") != appJSON {
		newErrorResponse(w, "unknown payload", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()

	var inp domain.CreateTokenRequest
	if err = json.Unmarshal(body, &inp); err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = inp.ValidateURL(); err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	link := domain.Link{
		OriginalURL: inp.OriginalURL,
	}

	token, err := h.services.Link.CreateToken(link)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(tokenResponse{Token: token})
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(resp) //nolint:errcheck
}

// @Summary Get long URL
// @Tags link
// @ID get-long-url
// @Product json
// @Param		token	path		string	true	"link"
// @Success	200		    {object}	linkResponse     "link"
// @Failure	400,404		{object}	errResponse
// @Failure	500			{object}	errResponse
// @Failure	default		{object}	errResponse
// @Router		/api/v1/link/{token} [get]
func (h *Handler) getOriginalURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	vars := mux.Vars(r)
	inp := domain.GetOriginalURLRequest{Token: vars["token"]}

	if err := inp.ValidateToken(); err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	originalURL, err := h.services.Link.GetOriginalURL(inp.Token)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(linkResponse{originalURL})
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp) //nolint:errcheck
}
