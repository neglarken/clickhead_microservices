package handler

import (
	"encoding/json"
	"net/http"

	"github.com/neglarken/clickhead/some-ms/internal/model"
)

func (h *Handler) CreateItem() http.HandlerFunc {
	type request struct {
		Info  string `json:"info"`
		Price int    `json:"price"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}
		item := &model.Item{
			Info:  req.Info,
			Price: req.Price,
		}
		id, err := h.repo.Item.Create(item)
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}
		type respond struct {
			Id int `json:"id"`
		}
		h.respond(w, r, http.StatusOK, respond{Id: id})
	}
}

func (h *Handler) EditItem() http.HandlerFunc {
	type request struct {
		Id    int    `json:"id"`
		Info  string `json:"info"`
		Price int    `json:"price"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}
		item := &model.Item{
			Id:    req.Id,
			Info:  req.Info,
			Price: req.Price,
		}
		if err := h.repo.Item.Edit(item); err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}
		h.respond(w, r, http.StatusOK, item)
	}
}

func (h *Handler) GetItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := h.repo.Item.Get()
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}
		h.respond(w, r, http.StatusOK, items)
	}
}

func (h *Handler) DeleteItem() http.HandlerFunc {
	type request struct {
		Id int `json:"id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := h.repo.Item.Delete(req.Id); err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}
		type respond struct {
			Id int `json:"id"`
		}
		h.respond(w, r, http.StatusOK, respond{Id: req.Id})
	}
}
