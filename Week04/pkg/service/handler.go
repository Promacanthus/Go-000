package service

import (
	"net/http"
	"strconv"

	xerrors "github.com/pkg/errors"

	"github.com/Promacanthus/Go-000/Week04/pkg/biz"

	"github.com/gorilla/mux"
)

type Handler interface {
	Save(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

type HandlerImp struct {
	biz.StringService
}

var _ Handler = (*HandlerImp)(nil)

func NewHandler(service biz.StringService) Handler {
	return &HandlerImp{service}
}

// Save godoc
// @Summary Save the target string
// @Description store string to db
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param name path string true "String Name"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Failure default {object} string
// @Router /v1/save/{name} [POST]
func (h *HandlerImp) Save(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	save, err := h.StringService.Save(vars["name"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if xerrors.Is(err, biz.ErrEmpty) {
			w.Write([]byte("no data"))
		}
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(save)))
}

// Get godoc
// @Summary Get the target string
// @Description read string from db
// @Accept  json
// @Produce  json
// @Param name path string true "String Name"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Failure default {object} string
// @Router /v1/get/{name} [GET]
func (h *HandlerImp) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	amount, err := h.StringService.GetAmount(vars["name"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if xerrors.Is(err, biz.ErrEmpty) {
			w.Write([]byte("no data"))
		}
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(amount)))
}
