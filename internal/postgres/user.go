package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"

	gohtmx "github.com/falagansoftware/auto-repair/internal"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db: db}
}

func (u *UserService) CreateUser(ctx context.Context, user *gohtmx.User) (*gohtmx.User, error) {
	tx, err := u.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Execute query
	query := `INSERT INTO users (uid, name, surname, email, active) VALUES ($1, $2, $3, $4, $5) RETURNING uid, name, surname, email, active, created_at, updated_at`
	log.Print(query)
	res, err := tx.QueryContext(ctx, query, user.Uid, user.Name, user.Surname, user.Email, user.Active)
	if err != nil {
		return nil, err
	}
	defer res.Close()



	

func (u *UserService) FindUserByUid(ctx context.Context, uid string) (*gohtmx.User, error) {
	tx, err := u.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := findUserByUid(ctx, tx, uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) FindUsers(ctx context.Context, filters *gohtmx.UserFilters) ([]*gohtmx.User, error) {
	tx, err := u.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	users, _, err := findUsers(ctx, tx, filters)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserService) FindUsersGlobally(ctx context.Context, search *string) ([]*gohtmx.User, error) {
	tx, err := u.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	users, _, err := findUsersGlobally(ctx, tx, search)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Helpers

func createUser(ctx context.Context, tx *Tx, user *gohtmx.User) (*gohtmx.User, error) {
	result, err := tx.ExecContext(ctx, `
		INSERT INTO users (
			name,
			surname,
			email,
			password
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		user.Name,
		user.Surname,
		user.Email,
		user.Password,
	)
	if err != nil {
		return FormatError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)

	return nil
}

func findUserByUid(ctx context.Context, tx *Tx, uid string) (*gohtmx.User, error) {
	users, _, err := findUsers(ctx, tx, &gohtmx.UserFilters{Uid: &uid})

	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, &gohtmx.Error{Code: gohtmx.ENOTFOUND, Message: "User not found"}
	}
	return users[0], nil

}

func findUsers(ctx context.Context, tx *Tx, filter *gohtmx.UserFilters) (u []*gohtmx.User, n int, e error) {
	// Where clause based on filters props
	where := []string{"1 = 1"}
	orderBy := "name"
	direction := "ASC"

	if v := filter.Uid; v != nil {
		condition := fmt.Sprintf("uid = '%v'", *v)
		where = append(where, condition)
	}

	if v := filter.Name; v != nil {
		condition := fmt.Sprintf("name = '%v'", *v)
		where = append(where, condition)
	}

	if v := filter.Surname; v != nil {
		condition := fmt.Sprintf("surname = '%v'", *v)
		where = append(where, condition)
	}

	if v := filter.Email; v != nil {
		condition := fmt.Sprintf("email = '%v'", *v)
		where = append(where, condition)
	}

	if v := filter.Active; v {
		condition := fmt.Sprintf("active = '%v'", v)
		where = append(where, condition)
	}

	if v := filter.Sort; v != "" {
		orderBy = v
	}

	if v := filter.Order; v != "" {
		direction = v
	}

	// Execute query
	query := `SELECT uid, name, surname, email, active, created_at, updated_at FROM users WHERE ` + strings.Join(where, " AND ") + " ORDER BY " + orderBy + " " + direction + " " + FormatLimitOffset(filter.Limit, filter.Offset)
	log.Print(query)
	rows, err := tx.QueryContext(ctx, query)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	// Map rows to struct
	users := make([]*gohtmx.User, 0)

	for rows.Next() {
		var user gohtmx.User
		err := rows.Scan(&user.Uid, &user.Name, &user.Surname, &user.Email, &user.Active, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}

	// Check rows error
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, len(users), nil

}

func findUsersGlobally(ctx context.Context, tx *Tx, search *string) (u []*gohtmx.User, n int, e error) {
	// Execute query
	query := `SELECT uid, name, surname, email, active, created_at, updated_at FROM users WHERE CONCAT(name,'||',surname,'||', email) LIKE '%` + *search + `%' ORDER BY name ASC`
	log.Print(query)
	rows, err := tx.QueryContext(ctx, query)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	// Map rows to struct
	users := make([]*gohtmx.User, 0)

	for rows.Next() {
		var user gohtmx.User
		err := rows.Scan(&user.Uid, &user.Name, &user.Surname, &user.Email, &user.Active, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}

	// Check rows error
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, len(users), nil

}
