package repository

import (
	"log"
	"database/sql"

	"ginapp/internal/model/db"
	"ginapp/internal/model/entity"
	"ginapp/internal/dto"
)


func init() {
	db := db.GetDB()

	cmd := `
		CREATE TABLE IF NOT EXISTS USERS (
			UID INTEGER PRIMARY KEY AUTOINCREMENT,
			USERNAME VARCHAR(100) NOT NULL UNIQUE,
			PASSWORD VARCHAR(100) NOT NULL,
			CREATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
			UPDATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
		);

		CREATE TRIGGER IF NOT EXISTS trigger_users_updated_at AFTER UPDATE ON USERS
		BEGIN
    		UPDATE USERS 
    			SET UPDATE_AT = DATETIME('now', 'localtime') 
    			WHERE rowid == NEW.rowid;
		END;`

	_, err := db.Exec(cmd)

	if err != nil {
		log.Panic(err)
	}
}


type UserRepository interface {
	Select() ([]entity.User, error)
    SelectByUId(uid int) (entity.User, error)
    UpdateByUId(uid int, user *entity.User) error
    DeleteByUId(uid int) error

    Signup(sd *dto.SignupDto) error
    SelectByUsername(username string) (entity.User, error)
    UpdatePasswordByUId(uid int, password string) error
    UpdateUsernameByUId(uid int, username string) error
}


type userRepository struct {
	db *sql.DB
}


func NewUserRepository() UserRepository {
	db := db.GetDB()
	return &userRepository{db}
}


func (ur *userRepository) Select() ([]entity.User, error){
	var users []entity.User

	rows, err := ur.db.Query(
		`SELECT 
			UID,
			USERNAME, 
			CREATE_AT, 
			UPDATE_AT 
		 FROM USERS`,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := entity.User{}
		err = rows.Scan(&user.UId, &user.Username, &user.CreateAt, &user.UpdateAt)
		if err != nil {
			break
		}
		users = append(users, user)
	}

	return users, err
}


func (ur *userRepository) SelectByUId(uid int) (entity.User, error){
	var user entity.User
	err := ur.db.QueryRow(
		`SELECT 
			UID, 
			USERNAME, 
			CREATE_AT, 
			UPDATE_AT 
		 FROM USERS 
		 WHERE UID = ?`, uid,
	).Scan(
		&user.UId, &user.Username, &user.CreateAt, &user.UpdateAt,
	)

	return user, err
}


func (ur *userRepository) UpdateByUId(uid int, user *entity.User) error {
	_, err := ur.db.Exec(
		`UPDATE USERS 
		 SET 
			USERNAME = ? 
		 WHERE UID = ?`,
		user.Username, uid,
	)
	return err
}


func (ur *userRepository) DeleteByUId(uid int) error {
	_, err := ur.db.Exec(`DELETE FROM USERS WHERE UID = ?`, uid)

	return err
}


func (ur *userRepository) Signup(sd *dto.SignupDto) error {
	_, err := ur.db.Exec(
		`INSERT INTO USERS (
			USERNAME, PASSWORD
		 ) VALUES(?,?)`,
		sd.Username, sd.Password,
	)

	return err
}


func (ur *userRepository) SelectByUsername(username string) (entity.User, error){
	var user entity.User
	err := ur.db.QueryRow(
		`SELECT 
			UID, 
			USERNAME, 
			PASSWORD, 
			CREATE_AT, 
			UPDATE_AT 
		 FROM USERS 
		 WHERE USERNAME = ?`, username,
	).Scan(
		&user.UId, &user.Username, &user.Password, &user.CreateAt, &user.UpdateAt,
	)

	return user, err
}


func (ur *userRepository) UpdatePasswordByUId(uid int, password string) error {
	_, err := ur.db.Exec(
		`UPDATE USERS 
		 SET 
			PASSWORD = ? 
		 WHERE UID = ?`, password, uid,
	)
	return err
}


func (ur *userRepository) UpdateUsernameByUId(uid int, username string) error {
	_, err := ur.db.Exec(
		`UPDATE USERS 
		 SET 
		 	USERNAME = ? 
		 WHERE UID = ?`, username, uid,
	)
	return err
}
