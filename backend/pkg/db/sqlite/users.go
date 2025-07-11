package db

import (
	"database/sql"
	"strings"

	"social-network/pkg/models"
)

type UserRepository struct {
	DB *sql.DB
}

// Insert new user in db
func (repo *UserRepository) Add(user models.User) error {
	// example code
	stmt, err := repo.DB.Prepare("INSERT INTO users(user_id, email, first_name, last_name, nickname, about, password, birthday, image) values(?,?,?,?,(NULLIF(?,'')),?,?,?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(user.ID, user.Email, user.FirstName, user.LastName, user.Nickname, user.About, user.Password, user.DateOfBirth, user.ImagePath); err != nil {
		return err
	}
	return nil
}

// check if email already registered -> must be  unique
func (repo *UserRepository) EmailNotTaken(email string) (bool, error) {
	row := repo.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? LIMIT 1", email)
	var result int
	if err := row.Scan(&result); err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
	}
	if result != 0 {
		return false, nil
	}
	return true, nil
}

// find user by email and return user_id and password / mainly for login funcionality
func (repo *UserRepository) FindUserByEmail(email string) (models.User, error) {
	row := repo.DB.QueryRow("SELECT user_id,password FROM users WHERE email = ? LIMIT 1", email)
	var user models.User
	if err := row.Scan(&user.ID, &user.Password); err != nil {
		return user, err
	}
	return user, nil
}

// Return list of all users except current and following/follower info
func (repo *UserRepository) GetAllAndFollowing(userID string) ([]models.User, error) {
	var users []models.User

	rows, err := repo.DB.Query("SELECT user_id, IFNULL(nickname, first_name || ' ' || last_name), (SELECT COUNT(*) FROM followers WHERE followers.user_id = '"+userID+"' AND follower_id = users.user_id) as follower,(SELECT COUNT(*) FROM followers WHERE followers.user_id = users.user_id AND follower_id = '"+userID+"') as following, image FROM users WHERE user_id != ? ;", userID)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user models.User
		var follower int
		var following int
		rows.Scan(&user.ID, &user.Nickname, &follower, &following, &user.ImagePath)
		// configure followers
		if follower == 1 {
			user.Follower = true
		}
		if following == 1 {
			user.Following = true
		}
		users = append(users, user)
	}
	return users, nil
}

// returns id, nickname and image
// min package for returning user data
func (repo *UserRepository) GetDataMin(userID string) (models.User, error) {
	row := repo.DB.QueryRow("SELECT user_id, first_name, last_name, nickname, image FROM users WHERE user_id = ? LIMIT 1", userID)
	var user models.User
	var nickname sql.NullString
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &nickname, &user.ImagePath); err != nil {
		return user, err
	}
	if nickname.Valid {
		user.Nickname = nickname.String
	}
	return user, nil
}

// returns true if current user is following
func (repo *UserRepository) IsFollowing(userID, currentUserID string) (bool, error) {
	row := repo.DB.QueryRow("SELECT COUNT() FROM followers WHERE user_id = ? AND follower_id = ?;", userID, currentUserID)
	var result int
	if err := row.Scan(&result); err != nil {
		return false, err
	}
	if result == 0 {
		return false, nil
	}
	return true, nil
}

// returns true lif profile public
func (repo *UserRepository) ProfileStatus(userID string) (string, error) {
	row := repo.DB.QueryRow("SELECT status FROM users WHERE user_id = ?;", userID)
	var status string
	if err := row.Scan(&status); err != nil {
		return "", err
	}
	return status, nil
}

// returns user information evaluating propperties
// if current user returns full data set
// if public profile -> returns full data set
// if private profile and following- full data set
func (repo *UserRepository) GetProfileMax(userID string) (models.User, error) {
	row := repo.DB.QueryRow("SELECT IFNULL(nickname, first_name || ' ' || last_name),first_name, last_name, image, email, strftime('%d.%m.%Y', birthday), about FROM users WHERE user_id = ? LIMIT 1", userID)
	var user models.User
	if err := row.Scan(&user.Nickname, &user.FirstName, &user.LastName, &user.ImagePath, &user.Email, &user.DateOfBirth, &user.About); err != nil {
		return user, err
	}
	user.ID = userID
	return user, nil
}

// returns user information
// if private profile and  current user not following -> small data set
func (repo *UserRepository) GetProfileMin(userID string) (models.User, error) {
	row := repo.DB.QueryRow("SELECT  IFNULL(nickname, first_name || ' ' || last_name), image FROM users WHERE user_id = ? LIMIT 1", userID)
	var user models.User
	if err := row.Scan(&user.Nickname, &user.ImagePath); err != nil {
		return user, err
	}
	user.ID = userID

	// Compte les followers
	row = repo.DB.QueryRow("SELECT COUNT(*) FROM followers WHERE user_id = ?", userID)
	if err := row.Scan(&user.FollowersCount); err != nil {
		user.FollowersCount = 0
	}
	// Compte les abonnements
	row = repo.DB.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ?", userID)
	if err := row.Scan(&user.FollowingCount); err != nil {
		user.FollowingCount = 0
	}

	return user, nil
}

func (repo *UserRepository) GetFollowers(userID string) ([]models.User, error) {
	var users []models.User

	rows, err := repo.DB.Query("SELECT user_id, first_name, last_name, IFNULL(nickname, first_name || ' ' || last_name), image FROM users WHERE (SELECT COUNT(*) FROM followers WHERE followers.user_id = ? AND follower_id = users.user_id) = 1", userID)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user models.User
		rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Nickname, &user.ImagePath)
		users = append(users, user)
	}
	return users, nil
}

func (repo *UserRepository) GetFollowing(userID string) ([]models.User, error) {
	var users []models.User

	rows, err := repo.DB.Query("SELECT user_id, first_name, last_name, IFNULL(nickname, first_name || ' ' || last_name), image FROM users WHERE (SELECT COUNT(*) FROM followers WHERE followers.follower_id = ? AND user_id = users.user_id) = 1", userID)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user models.User
		rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Nickname, &user.ImagePath)
		users = append(users, user)
	}
	return users, nil
}

// get current user status
func (repo *UserRepository) GetStatus(userID string) (string, error) {
	row := repo.DB.QueryRow("SELECT status FROM users WHERE user_id = ? LIMIT 1", userID)
	var status string
	if err := row.Scan(&status); err != nil {
		return "", err
	}
	return status, nil
}

// set new status
func (repo *UserRepository) SetStatus(user models.User) error {
	stmt, err := repo.DB.Prepare("UPDATE USERS SET (status) = ? WHERE user_id = (?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Status, user.ID)
	if err != nil {
		return err
	}
	return nil
}

// save new follower
func (repo *UserRepository) SaveFollower(userId, followerId string) error {
	// Start a transaction to ensure all operations succeed or fail together
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Add to followers table
	_, err = tx.Exec("INSERT INTO followers(user_id, follower_id) VALUES (?,?)", userId, followerId)
	if err != nil {
		return err
	}

	// Grant access to all existing almost_private posts created by the user being followed
	_, err = tx.Exec(`
		INSERT INTO almost_private (user_id, post_id)
		SELECT ?, post_id 
		FROM posts 
		WHERE created_by = ? 
		AND visibility = 'ALMOST_PRIVATE'
		AND post_id NOT IN (SELECT post_id FROM almost_private WHERE user_id = ?)
	`, followerId, userId, followerId)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// delete follower
func (repo *UserRepository) DeleteFollower(userId, followerId string) error {
	// Start a transaction to ensure all operations succeed or fail together
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove from followers table
	_, err = tx.Exec("DELETE FROM followers WHERE (user_id = ? AND follower_id = ?)", userId, followerId)
	if err != nil {
		return err
	}

	// Remove access from almost_private posts created by the user being unfollowed
	_, err = tx.Exec(`
		DELETE FROM almost_private 
		WHERE user_id = ? 
		AND post_id IN (SELECT post_id FROM posts WHERE created_by = ?)
	`, followerId, userId)
	if err != nil {
		return err
	}

	// Remove access from private_post_access for posts created by the user being unfollowed
	_, err = tx.Exec(`
		DELETE FROM private_post_access 
		WHERE user_id = ? 
		AND post_id IN (SELECT post_id FROM posts WHERE created_by = ?)
	`, followerId, userId)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// SearchUsers searches for users by first name, last name, or nickname
// Excludes the current user and returns limited info
func (repo *UserRepository) SearchUsers(query, currentUserID string) ([]models.User, error) {
	searchTerm := "%" + strings.ToLower(query) + "%"

	sqlQuery := `
		SELECT user_id, first_name, last_name, nickname, image
		FROM users
		WHERE user_id != ? 
		AND (LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(nickname) LIKE ?)
		ORDER BY first_name, last_name
		LIMIT 20
	`

	rows, err := repo.DB.Query(sqlQuery, currentUserID, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Nickname,
			&user.ImagePath,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
