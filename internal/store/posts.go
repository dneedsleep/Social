package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64      `json:"id"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	UserID    int64      `json:"user_id"`
	Tags      []string   `json:"tags"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Version   int64      `json:"version"`
	Comments  []Comment  `json:comments`
}

type PostWithMetaData struct {
	ID        int64      `json:"id"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	UserID    int64      `json:"user_id"`
	Tags      []string   `json:"tags"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Version   int64      `json:"version"`
	Comments  []Comment  `json:"comments"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`

	CommentCount int64 `json:"comment_count"`
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		INSERT INTO posts(content , title ,user_id , tags)
		VALUES ($1 , $2 , $3 ,$4) RETURNING id,created_at,updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostsStore) GetById(ctx context.Context, userId int64) (*Post, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		SELECT id, user_id ,title , content , created_at , tags , updated_at , version
		FROM posts
		WHERE id = $1
	`

	var post Post
	err := s.db.QueryRowContext(ctx, query, userId).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		pq.Array(&post.Tags),
		&post.UpdatedAt,
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err

		}
	}

	return &post, nil
}

func (s *PostsStore) DeleteById(ctx context.Context, postId int64) error {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		DELETE FROM posts 
		WHERE id = $1
	`

	result, err := s.db.ExecContext(ctx, query, postId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil

}

func (s *PostsStore) UpdateById(ctx context.Context, p *Post) error {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		UPDATE posts
		SET title = $1 , content = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		&p.Title,
		&p.Content,
		&p.ID,
		&p.Version,
	).Scan(&p.Version)

	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil

}

func (s *PostsStore) GetUserFeed(ctx context.Context, p *User, pg PaginatedFeedQuery) ([]PostWithMetaData, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		 SELECT p.id , p.title , p.user_id , p.content , p.created_at , p.tags , p.version , u.first_name , u.last_name ,
   COUNT(*) AS comment_count 
   FROM posts p
  LEFT JOIN comments c
   ON c.post_id = p.id
  JOIN users u
   ON u.id = p.user_id
  WHERE  p.user_id = $1
  (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
			(p.tags @> $5 OR $5 = '{}')
  GROUP BY p.id , u.first_name , u.last_name
  ORDER BY p.created_at ` + pg.Sort + `
	 LIMIT $2 OFFSET $3
	`

	rows, err := s.db.QueryContext(ctx, query, &p.ID, pg.Limit, pg.Offset, pg.Search, pg.Tags)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	posts := []PostWithMetaData{}

	for rows.Next() {
		var r PostWithMetaData
		err := rows.Scan(&r.ID, &r.Title, &r.UserID, &r.Content, &r.CreatedAt, pq.Array(&r.Tags), &r.Version, &r.FirstName, &r.LastName, &r.CommentCount)

		if err != nil {
			return nil, err
		}

		posts = append(posts, r)
	}

	return posts, nil

}
