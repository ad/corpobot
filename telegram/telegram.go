package telegram

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
	// "strconv"

	dlog "github.com/amoghe/distillog"
	"golang.org/x/net/proxy"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	database "github.com/ad/corpobot/db"
	sql "github.com/lazada/sqle"
)

// InitTelegram ...
func InitTelegram(token, proxyHost, proxyPort, proxyUser, proxyPassword string, debug bool) (bot *tgbotapi.BotAPI, err error) {
	var tr http.Transport

	if proxyHost != "" {
		tr = http.Transport{
			DialContext: func(_ context.Context, network, addr string) (net.Conn, error) {
				socksDialer, err2 := proxy.SOCKS5(
					"tcp",
					fmt.Sprintf("%s:%s", proxyHost, proxyPort),
					&proxy.Auth{User: proxyUser, Password: proxyPassword},
					proxy.Direct,
				)
				if err2 != nil {
					return nil, err2
				}

				return socksDialer.Dial(network, addr)
			},
		}
	}

	bot, err = tgbotapi.NewBotAPIWithClient(token, &http.Client{Transport: &tr})
	if err != nil {
		return nil, err
	}

	bot.Debug = debug

	dlog.Debugf("Authorized on account @%s", bot.Self.UserName)

	return bot, nil
}


func ProcessTelegramMessages(db *sql.DB, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		dlog.Infof("%s [%d] %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)

		message := database.TelegramMessage{
			UserID:   update.Message.From.ID,
			UserName: update.Message.From.UserName,
			Message:  update.Message.Text,
			Date:     time.Unix(int64(update.Message.Date), 0),
		}

		err2 := database.StoreTelegramMessage(db, message)
		if err2 != nil {
			dlog.Errorf("%s", err2)
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start", "startgroup", "repos":
				// ghuser := &database.GithubUser{
				// 	TelegramUserID: strconv.Itoa(update.Message.From.ID),
				// }

				if update.Message.Command() != "repos" && update.Message.CommandArguments() != "" {
					// if user, err3 := client.GetGithubUser(update.Message.CommandArguments()); err3 == nil {
					// 	if user.Name != "" {
					// 		msg.Text = "Hi, " + user.Name
					// 	} else {
					// 		msg.Text = "Hi, " + user.UserName
					// 	}

					// 	ghuser.Name = user.Name
					// 	ghuser.UserName = user.UserName
					// 	ghuser.Token = update.Message.CommandArguments()

					// 	dbuser, err4 := database.AddUserIfNotExist(db, ghuser)
					// 	if err4 != nil && err4.Error() != database.AlreadyExists {
					// 		msg.Text += "\nError on save your token, try /start again\n" + err4.Error()
					// 		_, err5 := bot.Send(msg)
					// 		if err5 != nil {
					// 			dlog.Errorln(err5)
					// 		}
					// 		continue
					// 	}
					// 	ghuser.ID = dbuser.ID
					// }
				} else {
					// if user, err20 := database.GetGithubUserFromDB(db, ghuser.TelegramUserID); err20 == nil {
					// 	ghuser.ID = user.ID
					// 	ghuser.Name = user.Name
					// 	ghuser.UserName = user.UserName
					// 	ghuser.Token = user.Token
					// }
				}

				// if ghuser.ID != 0 {
				// 	// if repos, err6 := client.GetGithubUserRepos(ghuser.Token, ghuser.UserName); err6 == nil {
				// 	// 	msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "You are watching:\n")
				// 	// 	for _, repo := range repos {
				// 	// 		msg2.Text += "[" + repo.FullName + "](https://github.com/" + repo.FullName + ") updated at:" + repo.UpdatedAt.Format("2006-01-02 15:04:05") + "\n"

				// 	// 		ghrepo := &database.GithubRepo{
				// 	// 			Name:     repo.Name,
				// 	// 			RepoName: repo.FullName,
				// 	// 		}

				// 	// 		if dbrepo, err7 := database.AddRepoIfNotExist(db, ghrepo); err7 != nil && err7.Error() != database.AlreadyExists {
				// 	// 			dlog.Errorln(err7)
				// 	// 		} else if err8 := database.AddRepoLinkIfNotExist(db, ghuser, dbrepo, repo.UpdatedAt); err8 != nil && err8.Error() != database.AlreadyExists {
				// 	// 			dlog.Errorln(err8)
				// 	// 		}
				// 	// 	}
				// 	// 	msg2.ParseMode = "Markdown"
				// 	// 	msg2.DisableWebPagePreview = true
				// 	// 	_, err9 := bot.Send(msg2)
				// 	// 	if err9 != nil {
				// 	// 		dlog.Errorln(err9)
				// 	// 	}

				// 	// 	continue
				// 	// } else {
				// 	// 	dlog.Errorln(err6)
				// 	// }
				// }

				// text := `[Click here to authorize bot in github](https://github.com/login/oauth/authorize?client_id=` + clientID + `&redirect_uri=` + httpRedirectURI + `), and then press START again`
				msg.ParseMode = "Markdown"
				msg.Text = "text"
				msg.DisableWebPagePreview = true
			case "me":
				// if user, err10 := database.GetGithubUserFromDB(db, strconv.Itoa(update.Message.From.ID)); err10 == nil {
				// 	if user.Name != "" {
				// 		msg.Text = "Hi, " + user.Name
				// 	} else {
				// 		msg.Text = "Hi, " + user.UserName
				// 	}
				// } else {
				// 	msg.Text = "type /start\n"
				// 	msg.Text += err10.Error()
				// }
			case "delete":
				// if checkRepoName(update.Message.CommandArguments()) {
				// 	if ghuser, err10 := database.GetGithubUserFromDB(db, strconv.Itoa(update.Message.From.ID)); err10 == nil {
				// 		if ghrepo, errGetRepo := database.GetGithubRepoByNameFromDB(db, update.Message.CommandArguments()); errGetRepo == nil {
				// 			if errDeleteRepo := database.DeleteRepoUserLinkDB(db, ghuser, ghrepo); err == nil {
				// 				dlog.Infof("%s %s %s", ghuser.Name, "removed", ghrepo.RepoName)
				// 				msg.Text = ghrepo.RepoName + " removed, uncheck Watching in Github interface"
				// 			} else {
				// 				msg.Text += errDeleteRepo.Error()
				// 			}
				// 		} else {
				// 			msg.Text = errGetRepo.Error()
				// 		}
				// 	} else {
				// 		msg.Text = "type /start\n"
				// 		msg.Text += err10.Error()
				// 	}
				// } else {
				// 	msg.Text = "wrong repo format, try username/reponame instead"
				// }
			case "add":
				// if checkRepoName(update.Message.CommandArguments()) {
					// if ghuser, err10 := database.GetGithubUserFromDB(db, strconv.Itoa(update.Message.From.ID)); err10 == nil {
					// 	if ghrepo, errCheckRepo := database.GetGithubRepoByNameFromDB(db, update.Message.CommandArguments()); errCheckRepo == nil {
					// 		if err8 := database.AddRepoLinkIfNotExist(db, ghuser, ghrepo, time.Now()); err8 != nil && err8.Error() != database.AlreadyExists {
					// 			dlog.Errorln(err8)
					// 		} else {
					// 			msg.Text = ghrepo.RepoName + " added"
					// 		}
					// 	} else {
					// 		if repo, errGetRepo := client.GetGithubRepo(ghuser.Token, update.Message.CommandArguments()); errGetRepo == nil {
					// 			ghrepo := &database.GithubRepo{
					// 				Name:     repo.Name,
					// 				RepoName: repo.FullName,
					// 			}
					// 			if dbrepo, err7 := database.AddRepoIfNotExist(db, ghrepo); err7 != nil && err7.Error() != database.AlreadyExists {
					// 				dlog.Errorln(err7)
					// 			} else if err8 := database.AddRepoLinkIfNotExist(db, ghuser, dbrepo, repo.UpdatedAt); err8 != nil && err8.Error() != database.AlreadyExists {
					// 				dlog.Errorln(err8)
					// 			} else {
					// 				msg.Text = ghrepo.RepoName + " added"
					// 			}
					// 		} else {
					// 			msg.Text = ghrepo.RepoName + " not found"
					// 		}
					// 	}
					// } else {
					// 	msg.Text = "type /start\n"
					// 	msg.Text += err10.Error()
					// }
				// } else {
				// 	msg.Text = "wrong repo format, try username/reponame instead"
				// }
			case "help":
				msg.Text = "type /start"
			default:
				msg.Text = "I don't know that command"
			}
			msg.ReplyToMessageID = update.Message.MessageID
			_, err11 := bot.Send(msg)
			if err11 != nil {
				dlog.Errorln(err11)
			}
		}
	}
}
