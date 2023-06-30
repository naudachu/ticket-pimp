# Сборка и запуск:
Первые шаги делаю на локальной машине:
1. Поменять в коде файл окружения на '.env'
2. Собрать контейнер: `docker build -t naudachu/ticket-pimp:latest --pull .`
3. Затолкать контейнер в docker hub: `docker push naudachu/ticket-pimp:latest`

Далее с сервера:
1. Вытягиваем новый образ: `docker pull naudachu/ticket-pimp`
2. Запускаем в фоне: `docker run -d naudachu/ticket-pimp`

# To-do P1:

- [x] Делать запросы в Git, ownCloud параллельно;
- [x] Сохранять правильную ссылку на Git;
- [x] Сохранять правильную ссылку на GitBuild;
- [x] Сделать бота в Telegram;


# To-do P2*:
- [ ] В уведомлении об успешном создании сообщать всю инфу: 
    - git;
    - git-build url + ssh url;
    - ссылку на графику;
    - добавлять название игры;
- [ ] Сохранять короткую ссылку на графику;
- [ ] Сохранять внешнюю ссылку на графику;
- [ ] Сделать бота в Discord;
- [ ] Подумать над нормальной обработкой ошибок, сейчас достаточно всрато;
- [ ] Складывать в описание репозитория ссылку на тикет;
- [ ] Сделать базулю с достойными пользователями;

- [x] Run bot on docker scratch: https://github.com/jeremyhuiskamp/golang-docker-scratch/blob/main/README.mdа