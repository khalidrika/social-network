package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/pkg/utils"
)

func (handler *Handler) Notifications(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)

	userId := r.Context().Value(utils.UserKey).(string)
	notifs, err := handler.repos.NotifRepo.GetAll(userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}

	for i := 0; i < len(notifs); i++ {
		switch notifs[i].Type {
		case "GROUP_INVITE":
			notifs[i].Group, _ = handler.repos.GroupRepo.GetData(notifs[i].Content)
			notifs[i].User, _ = handler.repos.UserRepo.GetDataMin(notifs[i].Sender)
		case "FOLLOW":
			notifs[i].User, _ = handler.repos.UserRepo.GetDataMin(notifs[i].Content)
		case "EVENT":
			notifs[i].Event, _ = handler.repos.EventRepo.GetData(notifs[i].Content)
			notifs[i].User, _ = handler.repos.UserRepo.GetDataMin(notifs[i].Sender)
			notifs[i].Group, _ = handler.repos.GroupRepo.GetData(notifs[i].Event.GroupID)
		case "GROUP_REQUEST":
			notifs[i].User, _ = handler.repos.UserRepo.GetDataMin(notifs[i].Content)
			notifs[i].Group, _ = handler.repos.GroupRepo.GetData(notifs[i].TargetID)
		}
		utils.DefineNotificationMsg(&notifs[i])
	}

	utils.RespondWithNotifications(w, notifs, 200)
}

func (handler *Handler) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)

	if r.Method != http.MethodPost {
		utils.RespondWithError(w, "Method not allowed", 405)
		return
	}

	var request struct {
		NotificationID string `json:"notificationId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, "Invalid request body", 400)
		return
	}

	if request.NotificationID == "" {
		utils.RespondWithError(w, "Notification ID is required", 400)
		return
	}

	err := handler.repos.NotifRepo.MarkAsRead(request.NotificationID)
	if err != nil {
		utils.RespondWithError(w, "Error marking notification as read", 500)
		return
	}

	utils.RespondWithSuccess(w, "Notification marked as read", 200)
}

func (handler *Handler) MarkAllNotificationsAsRead(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)

	if r.Method != http.MethodPost {
		utils.RespondWithError(w, "Method not allowed", 405)
		return
	}

	userId := r.Context().Value(utils.UserKey).(string)

	err := handler.repos.NotifRepo.MarkAllAsRead(userId)
	if err != nil {
		utils.RespondWithError(w, "Error marking notifications as read", 500)
		return
	}

	utils.RespondWithSuccess(w, "All notifications marked as read", 200)
}

func (handler *Handler) DismissNotification(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)

	if r.Method != http.MethodDelete {
		utils.RespondWithError(w, "Method not allowed", 405)
		return
	}

	var request struct {
		NotificationID string `json:"notificationId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, "Invalid request body", 400)
		return
	}

	if request.NotificationID == "" {
		utils.RespondWithError(w, "Notification ID is required", 400)
		return
	}

	err := handler.repos.NotifRepo.Delete(request.NotificationID)
	if err != nil {
		utils.RespondWithError(w, "Error dismissing notification", 500)
		return
	}

	utils.RespondWithSuccess(w, "Notification dismissed", 200)
}
