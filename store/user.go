package store

import (
	"errors"
	"strconv"
	"database/sql"
	"github.com/calenaur/raidtime/model"
	"github.com/calenaur/raidtime/config"
	_ "github.com/go-sql-driver/mysql"
)

type UserStore struct {
	db *sql.DB
	cfg *config.Config
}

func NewUserStore(db *sql.DB, cfg *config.Config) *UserStore {
	return &UserStore{
		db: db,
		cfg: cfg,
	}
}

func (us *UserStore) GetByID(id int64) (*model.User, error) {
	stmt, err := us.db.Prepare(`
		SELECT 
			u.id, u.username, u.discriminator, u.avatar,
			c.id, c.name, c.color,
			gr.id, gr.name, gr.protected,
			p.id, p.name, p.manage_users, p.manage_events
		FROM user u
		JOIN class c ON u.class = c.id
		JOIN guild_rank gr ON u.guild_rank = gr.id
		JOIN permissions p ON u.permissions = p.id
		WHERE u.id = ?
	`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	row := stmt.QueryRow(id)
	user, err := us.CreateUserFromRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStore) GetBySession(session string) (*model.User, error) {
	stmt, err := us.db.Prepare(`
		SELECT 
			u.id, u.username, u.discriminator, u.avatar,
			c.id, c.name, c.color,
			gr.id, gr.name, gr.protected,
			p.id, p.name, p.manage_users, p.manage_events
		FROM user u
		JOIN class c ON u.class = c.id
		JOIN guild_rank gr ON u.guild_rank = gr.id
		JOIN permissions p ON u.permissions = p.id
		WHERE session=? AND UNIX_TIMESTAMP(session_creation_time) + ? > UNIX_TIMESTAMP()
	`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	row := stmt.QueryRow(session, us.cfg.Session.SessionDuration)
	user, err := us.CreateUserFromRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStore) ValidateSession(session string) bool {
	stmt, err := us.db.Prepare(`
		SELECT id FROM user WHERE session=? AND UNIX_TIMESTAMP(session_creation_time) + ? > UNIX_TIMESTAMP()`)
	if err != nil {
		return false
	}

	defer stmt.Close()
	row := stmt.QueryRow(session, us.cfg.Session.SessionDuration)
	
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return false
	}

	return true
}

func (us *UserStore) StartSession(user *model.User) (string, error) {
	stmt, err := us.db.Prepare("UPDATE user SET session=?, session_creation_time=NOW() WHERE id=?")
	if err != nil {
		return "", err
	}

	defer stmt.Close()
	session := user.GenerateSession(us.cfg.Session.SessionSecret)
	result, err := stmt.Exec(session, user.ID)
	if result == nil || err != nil {
		return "", err
	}

	return session, nil
}

func (us *UserStore) SignupToEvent(user *model.User, event int, signupType int) error {
	stmt, err := us.db.Prepare(`
		INSERT INTO 
			signup (event, user, type) 
		VALUES 
			(?, ?, ?)
		ON DUPLICATE KEY UPDATE
			type=?, date=NOW()
	`)
	if err != nil {
		return err
	}

	defer stmt.Close()
	result, err := stmt.Exec(event, user.ID, signupType, signupType)
	if result == nil || err != nil {
		return err
	}

	return nil
}

func (us *UserStore) CancelSignup(user *model.User, event int) error {
	stmt, err := us.db.Prepare("DELETE FROM signup WHERE event=? AND user=?")
	if err != nil {
		return err
	}

	defer stmt.Close()
	result, err := stmt.Exec(event, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected < 1 {
		return errors.New("Could not cancel signup: No signup found")
	}

	return nil
}

func (us *UserStore) Login(credentials *model.UserCredentials) (*model.User, string, error) {
	stmt, err := us.db.Prepare(`
		INSERT INTO 
			user (id, username, discriminator, avatar)
		VALUES
			(?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			id=?, username=?, discriminator=?, avatar=?
	`)
	if err != nil {
		return nil, "", err
	}

	ID, err := strconv.ParseInt(credentials.ID, 10, 64)
	if err != nil {
		return nil, "", err
	}

	defer stmt.Close()
	result, err := stmt.Exec(
		ID, 
		credentials.Username, 
		credentials.Discriminator, 
		credentials.Avatar, 
		ID, 
		credentials.Username, 
		credentials.Discriminator, 
		credentials.Avatar, 
	)
	if result == nil || err != nil {
		return nil, "", err
	}

	user, err := us.GetByID(ID)
	if err != nil {
		return nil, "", err
	}

	session, err := us.StartSession(user)
	if err != nil {
		return nil, "", err
	}

	return user, session, nil
}

func (us *UserStore) Logout(user *model.User) error {
	stmt, err := us.db.Prepare("UPDATE user SET session='invalid', session_creation_time=0 WHERE id=?")
	if err != nil {
		return err
	}

	defer stmt.Close()
	result, err := stmt.Exec(user.ID)
	if result == nil || err != nil {
		return err
	}

	return nil
}

func (us *UserStore) CreateUserFromRow(row *sql.Row) (*model.User, error) {
	var (
		manageUsers byte
		manageEvents byte
		protected byte
	)
	user := &model.User{
		Class: &model.Class{},
		GuildRank: &model.GuildRank{},
		Permissions: &model.Permissions{},
	}
	err := row.Scan(
		&user.ID, 
		&user.Username, 
		&user.Discriminator, 
		&user.Avatar,
		&user.Class.ID,
		&user.Class.Name,
		&user.Class.Color,
		&user.GuildRank.ID,
		&user.GuildRank.Name,
		&protected,
		&user.Permissions.ID, 
		&user.Permissions.Name, 
		&manageUsers, 
		&manageEvents,
	)
	if err != nil {
		return nil, err
	}

	user.GuildRank.Protected = protected == 1
	user.Permissions.ManageUsers = manageUsers == 1
	user.Permissions.ManageEvents = manageEvents == 1
	return user, nil
}