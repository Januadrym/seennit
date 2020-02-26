package community

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/app/status"
	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/http/respond"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type (
	service interface {
		SearchCommunity(ctx context.Context, name string) (*types.Community, error)
		CreateCommunity(ctx context.Context, req *types.Community) (*types.Community, error)
		PrivateCommunity(ctx context.Context, Com string) error
		GetCommunity(ctx context.Context, name string) (*types.Community, error)
		EnrollUser(ctx context.Context, idCom string) error
		GetAll(ctx context.Context) ([]*types.Community, error)
		UpdateInfo(ctx context.Context, idCom string, comm *types.Community) error
		// Policy & User
		PromoteUser(ctx context.Context, idUser string, idCom string) error
		GetUsers(ctx context.Context, idCom string) ([]*types.User, error)
		GetAllMods(ctx context.Context, idCom string) ([]*types.User, error)
		TransferAdmin(ctx context.Context, idUser, nameCom string) error
		// Post
		SubmitPost(ctx context.Context, nameComm string, req *types.Post) (*types.Post, error)
		GetAllPosts(ctx context.Context, nameComm string) ([]*types.Post, error)
	}

	Handler struct {
		Svc service
	}
)

func NewHandler(svc service) *Handler {
	return &Handler{
		Svc: svc,
	}
}

func (h *Handler) CreateCommunity(w http.ResponseWriter, r *http.Request) {
	var req types.Community
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	comm, err := h.Svc.CreateCommunity(r.Context(), &req)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: comm,
	})
}

func (h *Handler) DeleteCom(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if name == "" {
		logrus.WithContext(r.Context()).Info("invalid id")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := h.Svc.PrivateCommunity(r.Context(), name); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: types.IDResponse{
			ID: name,
		},
	})
}

func (h *Handler) GetCommunity(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if name == "" {
		logrus.WithContext(r.Context()).Info("invalid name")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
		return
	}
	com, err := h.Svc.GetCommunity(r.Context(), name)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	if com.Status == types.CommunityStatusPublic {
		respond.JSON(w, http.StatusOK, types.BaseResponse{
			Data: com,
		})
		return
	}
	user := auth.FromContext(r.Context())
	if com.Status == types.CommunityStatusPrivate && user != nil && user.UserID == com.CreatedByID {
		respond.JSON(w, http.StatusOK, types.BaseResponse{
			Data: com,
		})
		return
	}
	respond.JSON(w, http.StatusNotFound, types.BaseResponse{
		Status: status.Gen().NotFound,
	})
}

func (h *Handler) GetAllCommunity(w http.ResponseWriter, r *http.Request) {
	list, err := h.Svc.GetAll(r.Context())
	if err != nil {
		logrus.WithContext(r.Context()).Errorf("failed to get all communities, err: %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: list,
	})
}

func (h *Handler) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	var req *types.Community
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	comName := mux.Vars(r)["name"]
	if comName == "" {
		logrus.WithContext(r.Context()).Info("invalid name")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
		return
	}
	com, err := h.Svc.SearchCommunity(r.Context(), comName)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	if err := h.Svc.UpdateInfo(r.Context(), com.ID, req); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: com,
	})
}

// Users

func (h *Handler) EnrollUser(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if name == "" {
		logrus.WithContext(r.Context()).Info("invalid name")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
		return
	}
	com, err := h.Svc.SearchCommunity(r.Context(), name)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	if err := h.Svc.EnrollUser(r.Context(), com.ID); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: com,
	})
}
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if name == "" {
		logrus.WithContext(r.Context()).Info("invalid name")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
		return
	}
	com, err := h.Svc.SearchCommunity(r.Context(), name)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	users, er := h.Svc.GetUsers(r.Context(), com.ID)
	if er != nil {
		respond.Error(w, er, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: users,
	})
}

func (h *Handler) GetAllMods(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if name == "" {
		logrus.WithContext(r.Context()).Info("invalid name")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
		return
	}
	com, err := h.Svc.SearchCommunity(r.Context(), name)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	mods, er := h.Svc.GetAllMods(r.Context(), com.ID)
	if er != nil {
		respond.Error(w, er, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: mods,
	})
}

// Policy stuff

func (h *Handler) PromoteMod(w http.ResponseWriter, r *http.Request) {
	var req *UserIDPolicyRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	comName := mux.Vars(r)["name"]
	if comName == "" {
		logrus.WithContext(r.Context()).Info("invalid name")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
		return
	}
	com, err := h.Svc.SearchCommunity(r.Context(), comName)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	if err := h.Svc.PromoteUser(r.Context(), req.ID, com.ID); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: types.IDResponse{
			ID: req.ID,
		},
	})
}

func (h *Handler) TransferAdmin(w http.ResponseWriter, r *http.Request) {
	var req UserIDPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	comName := mux.Vars(r)["name"]
	if comName == "" {
		logrus.WithContext(r.Context()).Info("invalid name")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
		return
	}
	if err := h.Svc.TransferAdmin(r.Context(), req.ID, comName); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: types.IDResponse{
			ID: req.ID,
		},
	})
}

// Post stuff

func (h *Handler) SubmitPost(w http.ResponseWriter, r *http.Request) {
	var req *types.Post
	comName := mux.Vars(r)["name"]
	if comName == "" {
		logrus.WithContext(r.Context()).Info("invalid name")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	post, err := h.Svc.SubmitPost(r.Context(), comName, req)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: post,
	})
}

func (h *Handler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	comName := mux.Vars(r)["name"]
	if comName == "" {
		logrus.WithContext(r.Context()).Info("invalid name")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	posts, err := h.Svc.GetAllPosts(r.Context(), comName)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: posts,
	})
}
