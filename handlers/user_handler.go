package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	internalCtx "project/internal/context"
	"project/internal/models"
	"project/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.service.Register(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loginResponse := loginResp{
		Token: token,
	}

	if err = json.NewEncoder(w).Encode(loginResponse); err != nil {
		log.Println(err)
	}

}

type loginResp struct {
	Token string `json:"token"`
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(internalCtx.UserIDKey)
	userIDInt, ok := userID.(int64)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if userIDInt == 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetMe(r.Context(), userIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(user); err != nil {
		log.Println(err)
	}
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(internalCtx.UserIDKey).(int64)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if userID == 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var request models.UpdateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.UpdateMe(r.Context(), int(userID), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(internalCtx.UserIDKey).(int64)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if userID == 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	err := h.service.DeleteMe(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int) 

	var req models.CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.CreateOrder(r.Context(), userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	orders, err := h.service.GetMyOrders(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

func (h *UserHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	idStr := r.PathValue("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrderByID(r.Context(), orderID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *UserHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	idStr := r.PathValue("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req models.UpdateOrderRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.UpdateOrder(r.Context(), orderID, userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	idStr := r.PathValue("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteOrder(r.Context(), orderID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) AdminGetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.AdminGetAllUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) AdminGetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.AdminGetAllOrders(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

func (h *UserHandler) AdminUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	targetUserID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req models.UpdateRoleRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.AdminUpdateRole(r.Context(), targetUserID, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(internalCtx.UserIDKey).(int64)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if userID == 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var request models.UpdateUserPassword
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.UpdateUserPassword(r.Context(), userID, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	idStr := r.PathValue("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.CancelOrder(r.Context(), orderID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetUserAndOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	res, err := h.service.GetUserAndOrders(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		ID     int            `json:"id"`
		Name   string         `json:"name"`
		Email  string         `json:"email"`
		Orders []models.Order `json:"orders"`
	}{
		ID:     res.User.ID,
		Name:   res.User.Name,
		Email:  res.User.Email,
		Orders: res.Orders,
	}

	json.NewEncoder(w).Encode(response)
}