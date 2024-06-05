package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"

	autorepair "github.com/falagansoftware/auto-repair/internal"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db: db}
}

func (u *UserService) CreateUser(ctx context.Context, user *autorepair.UserCreate) error {
	tx, err := u.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = createUser(ctx, tx, user)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (u *UserService) FindUserByUid(ctx context.Context, uid string) (*autorepair.User, error) {
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

func (u *UserService) FindUsers(ctx context.Context, filters *autorepair.UserFilters) ([]*autorepair.User, error) {
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

func (u *UserService) FindUsersGlobally(ctx context.Context, search *string) ([]*autorepair.User, error) {
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

func createUser(ctx context.Context, tx *Tx, user *autorepair.UserCreate) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO users (
			name,
			surname,
			email,
			password
		)
		VALUES ($1, $2, $3, $4)
	`,
		user.Name,
		user.Surname,
		user.Email,
		user.Password,
	)

	if err != nil {
		return err
	}

	return nil
}

func findUserByUid(ctx context.Context, tx *Tx, uid string) (*autorepair.User, error) {
	users, _, err := findUsers(ctx, tx, &autorepair.UserFilters{Uid: &uid})

	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, &autorepair.Error{Code: autorepair.ENOTFOUND, Message: "User not found"}
	}
	return users[0], nil

}

func findUsers(ctx context.Context, tx *Tx, filter *autorepair.UserFilters) (u []*autorepair.User, n int, e error) {
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
	users := make([]*autorepair.User, 0)

	for rows.Next() {
		var user autorepair.User
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

func findUsersGlobally(ctx context.Context, tx *Tx, search *string) (u []*autorepair.User, n int, e error) {
	// Execute query
	query := `SELECT uid, name, surname, email, active, created_at, updated_at FROM users WHERE CONCAT(name,'||',surname,'||', email) LIKE '%` + *search + `%' ORDER BY name ASC`
	log.Print(query)
	rows, err := tx.QueryContext(ctx, query)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	// Map rows to struct
	users := make([]*autorepair.User, 0)

	for rows.Next() {
		var user autorepair.User
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
