package db

import (
	"database/sql"

	"social-network/pkg/models"
)

type PostRepository struct {
	DB *sql.DB
}

// Returns all posts for user ->
// group posts if is a member
// all public posts
// Private posts if is a follower
// almost_private if has access
// all posts if user is an author
func (repo *PostRepository) GetAll(userID string) ([]models.Post, error) {
	var posts []models.Post
	rows, err := repo.DB.Query(`
		SELECT post_id, created_by, content, image, visibility, created_at FROM posts
		WHERE group_id IS NULL
		  AND (
			visibility = 'PUBLIC'
			OR (visibility = 'ALMOST_PRIVATE' AND (SELECT COUNT(*) FROM almost_private WHERE almost_private.post_id = posts.post_id AND almost_private.user_id = ?) = 1)
			OR (visibility = 'PRIVATE' 
				AND (SELECT COUNT(*) FROM private_post_access WHERE private_post_access.post_id = posts.post_id AND private_post_access.user_id = ?) = 1
				AND (
					created_by = ?
					OR (SELECT COUNT(*) FROM followers WHERE followers.user_id = posts.created_by AND follower_id = ?) = 1
					OR (SELECT status FROM users WHERE user_id = posts.created_by) = 'public'
				)
			)
			OR created_by = ?
		  )
		ORDER BY created_at DESC;
	`, userID, userID, userID, userID, userID)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post models.Post
		rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.ImagePath, &post.Visibility, &post.CreatedAt)
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *PostRepository) GetUserPosts(userID, currentUserID string) ([]models.Post, error) {
	var posts []models.Post
	rows, err := repo.DB.Query(`
		SELECT post_id, created_by, content, image, visibility, created_at FROM posts 
		WHERE group_id IS NULL 
		  AND created_by = ?
		  AND (
			visibility = 'PUBLIC'
			OR (visibility = 'ALMOST_PRIVATE' AND (SELECT COUNT(*) FROM almost_private WHERE almost_private.post_id = posts.post_id AND almost_private.user_id = ?) = 1)
			OR (visibility = 'PRIVATE' 
				AND (SELECT COUNT(*) FROM private_post_access WHERE private_post_access.post_id = posts.post_id AND private_post_access.user_id = ?) = 1
				AND (
					? = ?
					OR (SELECT COUNT(*) FROM followers WHERE followers.user_id = ? AND follower_id = ?) = 1
					OR (SELECT status FROM users WHERE user_id = ?) = 'public'
				)
			)
			OR ? = ?
		  )
		ORDER BY created_at DESC;
	`, userID, currentUserID, currentUserID, currentUserID, userID, userID, currentUserID, userID, currentUserID, userID)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post models.Post
		rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.ImagePath, &post.Visibility, &post.CreatedAt)
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *PostRepository) GetGroupPosts(groupID string) ([]models.Post, error) {
	var posts []models.Post
	rows, err := repo.DB.Query("SELECT post_id , created_by, content, image, created_at  FROM posts WHERE group_id = ? ORDER BY created_at DESC;", groupID)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post models.Post
		rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.ImagePath, &post.CreatedAt)
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *PostRepository) New(post models.Post) error {
	stmt, err := repo.DB.Prepare("INSERT INTO posts (post_id, group_id, created_by, content,image,visibility) values (?,(NULLIF(?,'')),?,?,?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(post.ID, post.GroupID, post.AuthorID, post.Content, post.ImagePath, post.Visibility); err != nil {
		return err
	}
	return nil
}

func (repo *PostRepository) SaveAccess(postId, userId string) error {
	stmt, err := repo.DB.Prepare("INSERT INTO almost_private (post_id, user_id) values (?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(postId, userId); err != nil {
		return err
	}
	return nil
}

func (repo *PostRepository) SavePrivateAccess(postId, userId string) error {
	stmt, err := repo.DB.Prepare("INSERT INTO private_post_access (post_id, user_id) values (?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(postId, userId); err != nil {
		return err
	}
	return nil
}
