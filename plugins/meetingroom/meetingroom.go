package meetingroom

import (
	"strings"

	database "github.com/ad/corpobot/db"
	"github.com/ad/corpobot/plugins"
	"github.com/ad/corpobot/telegram"
	dlog "github.com/amoghe/distillog"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Plugin struct{}

func init() {
	plugins.RegisterPlugin(&Plugin{})
}

func (m *Plugin) OnStart() {
	if !plugins.CheckIfPluginDisabled("meetingroom.Plugin", "enabled") {
		return
	}

	err := database.ExecSQL(plugins.DB, `CREATE TABLE IF NOT EXISTS meetingrooms (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"name" TEXT NOT NULL,
		"state" VARCHAR(32) NOT NULL,
		"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT "meetingroom_name" UNIQUE ("name") ON CONFLICT IGNORE
	);
	CREATE TRIGGER IF NOT EXISTS meetingrooms_updated_at_Trigger
	AFTER UPDATE On meetingrooms
	BEGIN
	   UPDATE meetingrooms SET updated_at = STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW') WHERE id = NEW.id;
	END;
	
	CREATE TABLE IF NOT EXISTS meetingroom_schedule (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"meetingroom_id" INTEGER NOT NULL,
		"start" timestamp NOT NULL,
		"end" timestamp NOT NULL,
		"creator" integer NOT NULL,
		"created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		"updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT "meeting_room_id" FOREIGN KEY ("meetingroom_id") REFERENCES "meetingrooms" ("id") ON DELETE CASCADE,
		CONSTRAINT "meeting_room_creator_id" FOREIGN KEY ("creator") REFERENCES "users" ("id") ON DELETE CASCADE,
		CONSTRAINT "uniq_schedule_start" UNIQUE ("meetingroom_id", "start") ON CONFLICT IGNORE,
		CONSTRAINT "uniq_schedule_end" UNIQUE ("meetingroom_id", "end") ON CONFLICT IGNORE
	  );
	  CREATE TRIGGER IF NOT EXISTS meetingroom_schedule_updated_at_trigger
	  AFTER UPDATE ON meetingroom_schedule
	  BEGIN
		 UPDATE meetingroom_schedule SET updated_at = STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW') WHERE id = NEW.id;
	  END;
	  
	  CREATE TRIGGER IF NOT EXISTS validate_meetingroom_schedule_trigger
	  BEFORE INSERT ON meetingroom_schedule
	  BEGIN
		SELECT
		  CASE
			  WHEN
				NEW.start < CURRENT_TIMESTAMP THEN
					RAISE ( ABORT, 'Invalid start, less than now' ) 
		END;
		SELECT
		  CASE
			  WHEN
				NEW.end <= CURRENT_TIMESTAMP THEN
					RAISE ( ABORT, 'Invalid end, less than now' ) 
		END;
		SELECT
		  CASE
			  WHEN
				NEW.end <= NEW.start THEN
					RAISE ( ABORT, 'Invalid end, less than start' ) 
		END;
	  END;`)

	if err != nil {
		dlog.Errorf("%s", err)
	}

	plugins.RegisterCommand("meetingroomcreate", "Add meetingroom", []string{"admin", "owner"}, meetingroomCreate)
	plugins.RegisterCommand("meetingroomlist", "Return list of meetingrooms", []string{"member", "admin", "owner"}, meetingroomList)
	plugins.RegisterCommand("meetingroomrename", "Rename meetingrooms", []string{"member", "admin", "owner"}, meetingroomRename)
	plugins.RegisterCommand("meetingroomdelete", "Delete meetingroom", []string{"admin", "owner"}, meetingroomDelete)
	plugins.RegisterCommand("meetingroomblock", "Block meetingroom", []string{"admin", "owner"}, meetingroomBlock)
	plugins.RegisterCommand("meetingroomactivate", "Activate meetingroom", []string{"admin", "owner"}, meetingroomActivate)

	plugins.RegisterCommand("meetingroomschedule", "Return schedule of meetingroom", []string{"member", "admin", "owner"}, meetingroomschedule)
	plugins.RegisterCommand("meetingroombook", "Book meetingroom", []string{"member", "admin", "owner"}, meetingroombook)
	plugins.RegisterCommand("meetingroomrebook", "Rebook meetingroom", []string{"member", "admin", "owner"}, meetingroomrebook)
	plugins.RegisterCommand("meetingroomunbook", "Unbook meetingroom", []string{"member", "admin", "owner"}, meetingroomunbook)
}

func (m *Plugin) OnStop() {
	dlog.Debugln("[meetingroom.Plugin] Stopped")

	plugins.UnregisterCommand("meetingroomcreate")
	plugins.UnregisterCommand("meetingroomlist")
	plugins.UnregisterCommand("meetingroomrename")
	plugins.UnregisterCommand("meetingroomdelete")
	plugins.UnregisterCommand("meetingroomblock")
	plugins.UnregisterCommand("meetingroomactivate")

	plugins.UnregisterCommand("meetingroomschedule")
	plugins.UnregisterCommand("meetingroombook")
	plugins.UnregisterCommand("meetingroomrebook")
	plugins.UnregisterCommand("meetingroomunbook")
}

var meetingroomList plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	meetingrooms, err := database.GetMeetingrooms(plugins.DB, strings.Fields(args))
	if err != nil {
		return err
	}

	if len(meetingrooms) == 0 {
		return telegram.Send(user.TelegramID, "meetingroom list is empty")
	}

	var meetingroomsList []string

	for _, m := range meetingrooms {
		meetingroomsList = append(meetingroomsList, "* "+m.String())
	}

	return telegram.Send(user.TelegramID, strings.Join(meetingroomsList, "\n"))
}

var meetingroomCreate plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	if args == "" {
		return telegram.Send(user.TelegramID, "failed: empty meetingroom name")
	}

	meetingroom := &database.Meetingroom{
		Name:  args,
		State: "active",
	}

	_, err := database.AddMeetingroomIfNotExist(plugins.DB, meetingroom)
	if err != nil {
		return telegram.Send(user.TelegramID, "failed: "+err.Error())
	}

	return telegram.Send(user.TelegramID, "meetingroom created")
}

var meetingroomRename plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	names := strings.Split(args, "\n")

	if len(names) != 2 {
		return telegram.Send(user.TelegramID, "failed: you must provide the names of the two meetingrooms with a new line between them")
	}

	oldName, newName := strings.TrimSpace(names[0]), strings.TrimSpace(names[1])

	if oldName == "" || newName == "" {
		return telegram.Send(user.TelegramID, "failed: you must provide the names of the two meetingrooms with a new line between them")
	}

	rows, err := database.UpdateMeetingroomName(plugins.DB, oldName, newName)
	if err != nil {
		return err
	}

	if rows != 1 {
		return telegram.Send(user.TelegramID, "failed")
	}

	return telegram.Send(user.TelegramID, "success")
}

var meetingroomDelete plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	meetingroom := &database.Meetingroom{
		Name:  args,
		State: "deleted",
	}

	rows, err := database.UpdateMeetingroomState(plugins.DB, meetingroom)
	if err != nil {
		return err
	}

	if rows != 1 {
		return telegram.Send(user.TelegramID, "failed")
	}

	return telegram.Send(user.TelegramID, "success")
}

var meetingroomBlock plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	meetingroom := &database.Meetingroom{
		Name:  args,
		State: "blocked",
	}

	rows, err := database.UpdateMeetingroomState(plugins.DB, meetingroom)
	if err != nil {
		return err
	}

	if rows != 1 {
		return telegram.Send(user.TelegramID, "failed")
	}

	return telegram.Send(user.TelegramID, "success")
}
var meetingroomActivate plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	meetingroom := &database.Meetingroom{
		Name:  args,
		State: "active",
	}

	rows, err := database.UpdateMeetingroomState(plugins.DB, meetingroom)
	if err != nil {
		return err
	}

	if rows != 1 {
		return telegram.Send(user.TelegramID, "failed")
	}

	return telegram.Send(user.TelegramID, "success")
}

var meetingroomschedule plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	return telegram.Send(user.TelegramID, "not implemented")
}

var meetingroombook plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	return telegram.Send(user.TelegramID, "not implemented")
}

var meetingroomrebook plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	return telegram.Send(user.TelegramID, "not implemented")
}

var meetingroomunbook plugins.CommandCallback = func(update *tgbotapi.Update, command, args string, user *database.User) error {
	return telegram.Send(user.TelegramID, "not implemented")
}
