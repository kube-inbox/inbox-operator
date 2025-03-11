package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	inboxv1 "kubeinbox.com/inbox-operator/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// InboxHandler handles HTTP requests for Inbox resources
type InboxHandler struct {
	client client.Client
	logger logr.Logger
}

// NewInboxHandler creates a new InboxHandler
func NewInboxHandler(client client.Client, logger logr.Logger) *InboxHandler {
	return &InboxHandler{
		client: client,
		logger: logger,
	}
}

// ListInboxes handles GET requests to list Inbox resources
func (h *InboxHandler) ListInboxes(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	inboxList := &inboxv1.InboxList{}
	opts := []client.ListOption{
		client.InNamespace(namespace),
	}

	err := h.client.List(context.Background(), inboxList, opts...)
	if err != nil {
		h.logger.Error(err, "Failed to list inboxes")
		http.Error(w, "Failed to list inboxes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inboxList)
}

// GetInbox handles GET requests to retrieve a specific Inbox resource
func (h *InboxHandler) GetInbox(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	name := r.URL.Query().Get("name")

	if namespace == "" || name == "" {
		http.Error(w, "Namespace and name are required", http.StatusBadRequest)
		return
	}

	inbox := &inboxv1.Inbox{}
	err := h.client.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: name}, inbox)
	if err != nil {
		if errors.IsNotFound(err) {
			http.Error(w, "Inbox not found", http.StatusNotFound)
			return
		}
		h.logger.Error(err, "Failed to get inbox")
		http.Error(w, "Failed to get inbox", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inbox)
}
