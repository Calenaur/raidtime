package store
import (
	"time"
	"database/sql"
	"github.com/calenaur/raidtime/config"
	"github.com/calenaur/raidtime/model"
	_ "github.com/go-sql-driver/mysql"
)

type EventStore struct {
	db *sql.DB
	cfg *config.Config
}

func NewEventStore(db *sql.DB, cfg *config.Config) *EventStore {
	return &EventStore{
		db: db,
		cfg: cfg,
	}
}

func (es *EventStore) GetEventsByDateRange(start time.Time, end time.Time) ([]*model.Event, error) {
	stmt, err := es.db.Prepare(`
		SELECT 
			e.id, e.name, e.date, 
			c.id, c.name, c.color,
			u.id, u.username, u.discriminator, u.avatar, u.guild_rank,
			class.id, class.name, class.color,
			p.id, p.name, p.manage_users, p.manage_events 
		FROM event e 
		JOIN color c ON e.color = c.id 
		JOIN user u ON e.creator = u.id
		JOIN class ON u.class = class.id
		JOIN permissions p ON u.permissions = p.id
		WHERE e.date >= ? AND e.date <= ?
		ORDER BY e.name ASC
	`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(start, end)
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()
	events, err := es.CreateEventsFromRows(rows)
	if err != nil {
		return nil, err
	}

	return events, nil
} 

func (es *EventStore) GetSignupsByEvent(event *model.Event) ([]*model.Signup, error){
	stmt, err := es.db.Prepare(`
		SELECT 
			u.id, u.username, u.discriminator, u.avatar, u.guild_rank,
			c.id, c.name, c.color,
			p.id, p.name, p.manage_users, p.manage_events,
			s.date,
			st.id, st.will_attend, st.description
		FROM signup s
		JOIN user u ON s.user = u.id
		JOIN class c ON u.class = c.id
		JOIN permissions p ON u.permissions = p.id
		JOIN signup_type st ON s.type = st.id
		WHERE s.event = ?
		ORDER BY st.id, c.id, u.guild_rank, u.username ASC
	`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(event.ID)
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()
	signups, err := es.CreateSignupsFromRows(event, rows)
	if err != nil {
		return nil, err
	}

	return signups, nil
}

func (es *EventStore) CreateEventsFromRows(rows *sql.Rows) ([]*model.Event, error) {
	events := []*model.Event{}
	for rows.Next() {
		event := &model.Event{
			Color: &model.Color{},
			Creator: &model.User{
				Class: &model.Class{},
				Permissions: &model.Permissions{},
			},
			Signups: []*model.Signup{},
		}

		var (
			manageUsers byte
			manageEvents byte
		)
		err := rows.Scan(
			&event.ID, 
			&event.Name, 
			&event.Date, 
			&event.Color.ID, 
			&event.Color.Name, 
			&event.Color.Color, 
			&event.Creator.ID, 
			&event.Creator.Username,
			&event.Creator.Discriminator, 
			&event.Creator.Avatar,
			&event.Creator.GuildRank,
			&event.Creator.Class.ID,
			&event.Creator.Class.Name,
			&event.Creator.Class.Color,
			&event.Creator.Permissions.ID,
			&event.Creator.Permissions.Name,
			&manageUsers,
			&manageEvents,
		)
		if err != nil {
			return nil, err
		}

		event.Creator.Permissions.ManageUsers = manageUsers == 1
		event.Creator.Permissions.ManageEvents = manageEvents == 1
		signups, err := es.GetSignupsByEvent(event)
		if err != nil {
			return nil, err
		}

		event.Signups = signups
		events = append(events, event)
	}
	return events, nil
}

func (es *EventStore) CreateSignupsFromRows(event *model.Event, rows *sql.Rows) ([]*model.Signup, error) {
	signups := []*model.Signup{}
	for rows.Next() {
		signup := &model.Signup{
			Event: event,
			User: &model.User{
				Class: &model.Class{},
				Permissions: &model.Permissions{},
			},
			SignupType: &model.SignupType{},
		}

		var (
			manageUsers byte
			manageEvents byte
			willAttend byte
		)
		err := rows.Scan(
			&signup.User.ID, 
			&signup.User.Username, 
			&signup.User.Discriminator, 
			&signup.User.Avatar,
			&signup.User.GuildRank,
			&signup.User.Class.ID,
			&signup.User.Class.Name,
			&signup.User.Class.Color,
			&signup.User.Permissions.ID,
			&signup.User.Permissions.Name,
			&manageUsers,
			&manageEvents, 
			&signup.Date,
			&signup.SignupType.ID,
			&willAttend,
			&signup.SignupType.Description,
		)
		if err != nil {
			return nil, err
		}

		signup.User.Permissions.ManageUsers = manageUsers == 1
		signup.User.Permissions.ManageEvents = manageEvents == 1
		signup.SignupType.WillAttend = willAttend == 1
		signups = append(signups, signup)
	}
	return signups, nil
}