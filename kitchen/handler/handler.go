package handler

import (
    "encoding/json"
    "errors"
    "net/http"
    "strconv"

    "kitchen/store"
)

type Handler struct {
    Store store.MenuStore
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
    writeJSON(w, status, map[string]any{
        "error": map[string]string{"code": code, "message": message},
    })
}

func (h *Handler) ListMenu(w http.ResponseWriter, r *http.Request) {
    menus := h.Store.List(r.URL.Query().Get("type"))
    writeJSON(w, http.StatusOK, menus)
}

func (h *Handler) GetMenu(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Path[len("/menu/"):]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        writeError(w, http.StatusBadRequest, "BAD_ID", "เลขจานต้องเป็นตัวเลข")
        return
    }
    menu, err := h.Store.Get(id)
    if errors.Is(err, store.ErrNotFound) {
        writeError(w, http.StatusNotFound, "MENU_NOT_FOUND", "ไม่พบเมนูหมายเลขนี้")
        return
    }
    writeJSON(w, http.StatusOK, menu)
}

func (h *Handler) CreateMenu(w http.ResponseWriter, r *http.Request) {
    var m store.Menu
    if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
        writeError(w, http.StatusBadRequest, "BAD_JSON", "อ่านกล่อง JSON ไม่ออก")
        return
    }
    if m.Name == "" || m.Price <= 0 {
        writeError(w, http.StatusBadRequest, "MISSING_FIELD",
            "ต้องมีชื่อเมนู และราคาต้องมากกว่าศูนย์")
        return
    }
    created := h.Store.Add(m)
    writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) DeleteMenu(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Path[len("/menu/"):]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        writeError(w, http.StatusBadRequest, "BAD_ID", "เลขจานต้องเป็นตัวเลข")
        return
    }
    if err := h.Store.Delete(id); errors.Is(err, store.ErrNotFound) {
        writeError(w, http.StatusNotFound, "MENU_NOT_FOUND", "ไม่พบเมนูหมายเลขนี้")
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Routes() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/menu", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            h.ListMenu(w, r)
        case http.MethodPost:
            h.CreateMenu(w, r)
        case http.MethodOptions:
            w.WriteHeader(http.StatusNoContent)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })
    mux.HandleFunc("/menu/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            h.GetMenu(w, r)
        case http.MethodDelete:
            h.DeleteMenu(w, r)
        case http.MethodOptions:
            w.WriteHeader(http.StatusNoContent)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })
    return mux
}

func WithCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
        w.Header().Set("Access-Control-Allow-Methods",
            "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}
