package db

import (
	"database/sql"

	"social-network/pkg/models"
)

type NotifRepository struct {
	DB *sql.DB
}

func (repo *NotifRepository) Save(n models.Notification) error {
	stmt, err := repo.DB.Prepare(`
		INSERT INTO notifications (notif_id, user_id, type, content, sender) 
		VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(n.ID, n.TargetID, n.Type, n.Content, n.Sender)
	return err
}

func (repo *NotifRepository) Delete(notificationId string) error {
	_, err := repo.DB.Exec(`DELETE FROM notifications WHERE notif_id = ?`, notificationId)
	return err
}

func (repo *NotifRepository) DeleteByType(n models.Notification) error {
	_, err := repo.DB.Exec(`
		DELETE FROM notifications 
		WHERE user_id = ? AND type = ? AND content = ?`, n.TargetID, n.Type, n.Content)
	return err
}

func (repo *NotifRepository) GetGroupRequests(groupId string) ([]models.Notification, error) {
	var notifs []models.Notification
	rows, err := repo.DB.Query(`
		SELECT content, notif_id, type, sender, user_id 
		FROM notifications 
		WHERE user_id = ? AND type = 'GROUP_REQUEST'`, groupId)
	if err != nil {
		return notifs, err
	}
	defer rows.Close()

	for rows.Next() {
		var n models.Notification
		rows.Scan(&n.Content, &n.ID, &n.Type, &n.Sender, &n.TargetID)
		notifs = append(notifs, n)
	}
	return notifs, nil
}

func (repo *NotifRepository) GetGroupInvites(groupId string) ([]models.Notification, error) {
	var notifs []models.Notification
	rows, err := repo.DB.Query(`
		SELECT notif_id, user_id, type, content, sender 
		FROM notifications 
		WHERE content = ? AND type = 'GROUP_INVITE'`, groupId)
	if err != nil {
		return notifs, err
	}
	defer rows.Close()

	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.TargetID, &n.Type, &n.Content, &n.Sender); err == nil {
			notifs = append(notifs, n)
		}
	}
	return notifs, nil
}

func (repo *NotifRepository) DeleteGroupInvite(userId, groupId string) error {
	_, err := repo.DB.Exec(`
		DELETE FROM notifications 
		WHERE user_id = ? AND content = ? AND type = 'GROUP_INVITE'`, userId, groupId)
	return err
}

func (repo *NotifRepository) GetUserFromRequest(notificationId string) (string, error) {
	var userId string
	err := repo.DB.QueryRow(`
		SELECT content FROM notifications 
		WHERE notif_id = ?`, notificationId).Scan(&userId)
	return userId, err
}

func (repo *NotifRepository) CheckIfExists(n models.Notification) (bool, error) {
	var count int
	err := repo.DB.QueryRow(`
		SELECT COUNT(*) FROM notifications 
		WHERE user_id = ? AND content = ? AND type = ?`, n.TargetID, n.Content, n.Type).Scan(&count)
	return count > 0, err
}

func (repo *NotifRepository) GetGroupId(notificationId string) (string, error) {
	var groupId string
	err := repo.DB.QueryRow(`
		SELECT content FROM notifications 
		WHERE notif_id = ?`, notificationId).Scan(&groupId)
	return groupId, err
}

func (repo *NotifRepository) GetAll(userId string) ([]models.Notification, error) {
	var notifs []models.Notification
	rows, err := repo.DB.Query(`
		SELECT content, notif_id, type, sender, user_id, COALESCE(read, FALSE) 
		FROM notifications 
		WHERE user_id = ? 
		OR (SELECT administrator FROM groups WHERE group_id = notifications.user_id) = ?
		ORDER BY notif_id DESC`, userId, userId)
	if err != nil {
		return notifs, err
	}
	defer rows.Close()

	for rows.Next() {
		var n models.Notification
		rows.Scan(&n.Content, &n.ID, &n.Type, &n.Sender, &n.TargetID, &n.Read)
		notifs = append(notifs, n)
	}
	return notifs, nil
}

func (repo *NotifRepository) GetCahtNotifById(notificationId string) (models.Notification, error) {
	var n models.Notification
	err := repo.DB.QueryRow(`
		SELECT content, user_id, sender 
		FROM notifications 
		WHERE notif_id = ?`, notificationId).Scan(&n.Content, &n.TargetID, &n.Sender)
	return n, err
}

func (repo *NotifRepository) CheckIfChatRequestExists(senderId, receiverId string) (bool, error) {
	var count int
	err := repo.DB.QueryRow(`
		SELECT COUNT(*) FROM notifications 
		WHERE user_id = ? AND sender = ? AND type = 'CHAT_REQUEST'`, receiverId, senderId).Scan(&count)
	return count > 0, err
}

func (repo *NotifRepository) GetContentFromChatRequest(senderId, receiverId string) (string, error) {
	var content string
	err := repo.DB.QueryRow(`
		SELECT content FROM notifications 
		WHERE user_id = ? AND sender = ? AND type = 'CHAT_REQUEST'`, receiverId, senderId).Scan(&content)
	return content, err
}

func (repo *NotifRepository) MarkAsRead(notificationId string) error {
	_, err := repo.DB.Exec(`
		UPDATE notifications 
		SET read = TRUE 
		WHERE notif_id = ?`, notificationId)
	return err
}

func (repo *NotifRepository) MarkAllAsRead(userId string) error {
	_, err := repo.DB.Exec(`
		UPDATE notifications 
		SET read = TRUE 
		WHERE user_id = ?`, userId)
	return err
}
