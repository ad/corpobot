# Corpobot

Бот для организации доступа к чатам компании.

### Новый пользователь

- Пользователя просят написать боту
- Пользователь пишет боту /start
- Бот записывает пользователя в базу со статусом new
- Бот отправляет уведомление о новом пользователе администратору(или нескольким)
- Администратор меняет роль пользователя на "участник"
- Администратор помещает пользователя в группу (или несколько)
- Бот отправляет пользователю список доступных команд
- Пользователь может отправив команду боту получить список доступных чатов и ссылки для входа в них

### Удаление пользователя

- Администратор выбирает какого пользователя удалить из каких чатов (пользователь банится в этих чатах)
- Администратор удаляет пользователя, чем запрещает ему выполнение любых команд бота

## Команды

Those are my commands: 
- /broadcast - Send message to all users
- /groupaddgroupchat - Add groupchat to group
- /groupadduser - Add user to group
- /groupchatdelete - Delete groupchat
- /groupchatinvitegenerate - Generate groupchat invite link
- /groupchatlist - Groupchat list
- /groupchatmembers - List groupchat members
- /groupchatuserban - Ban user in groupchat
- /groupchatuserunban - Unban user in groupchat
- /groupcreate - Create group
- /groupdelete - Delete group
- /groupdeletegroupchat - Delete groupchat from group
- /groupdeleteuser - Delete user from group
- /grouplist - Group list
- /grouprename - Rename group
- /groupundelete - Undelete group
- /help - Display this help
- /me - Your ID/username
- /meetingroomactivate - Activate meetingroom
- /meetingroomblock - Block meetingroom
- /meetingroombook - Book meetingroom
- /meetingroomcreate - Add meetingroom
- /meetingroomdelete - Delete meetingroom
- /meetingroomlist - Return list of meetingrooms
- /meetingroomrebook - Rebook meetingroom
- /meetingroomrename - Rename meetingrooms
- /meetingroomschedule - Return schedule of meetingroom
- /meetingroomscheduleinfo - Return schedule info
- /meetingroomunbook - Unbook meetingroom
- /message - Send message to user
- /plugindisable - Disable plugin
- /pluginenable - Enable plugin
- /pluginlist - List of plugins
- /start - Bot /start command
- /user - User actions
- /userbirthday - Set user birthday
- /userblock - Block user
- /userdelete - Delete user
- /userlist - User list
- /userpromote - Change user role
- /userunblock - Unblock user
- /userundelete - Undelete user

## TODO
- [ ] Поздравлять пользователя с днем рождения и уведомлять админов (заранее)

## Идеи для плагинов
- [ ] Отпуска
- [ ] Онбординг
- [ ] Контакты
