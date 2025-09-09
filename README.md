# Запуск

```shell
cp .env.example .env
```

```shell
cp config_example.json config.json
```

Заполнить `TELEGRAM_BOT_API_TOKEN` в `.env`

Обязательно заполнить `chat_id`, опционально заполнение `users` в `config.json`, пицца будет работать, статистика будет собираться, но слежки за аватарками не будет.

Ключи объектов в `users` == `user_id`. 
Чтобы узнать `user_id`, можно воспользоваться сайтом https://tg-user.id/ или любым другим способом, который позволит узнать `id` пользователя телеграм.

```shell
docker compose up --build
```

1 экземпляр приложения работает с одним телеграм бот токеном и одним чатом. Возможно когда-нибудь сделаю чтобы бот работал в нескольких чатах.  

Логи сохраняются в файл в `docker volume`, если надо посмотреть их, то путь можно узнать через 
```shell
docker volume inspect telegram_pfp_spy_error_logs
```
В поле `Mountpoint` будет путь к папке, там errors.log.

```shell
sudo cat /var/lib/docker/volumes/telegram_pfp_spy_error_logs/_data/errors.log
```

# Миграции

Нужен `goose` https://github.com/pressly/goose

Создать файл с миграцией

```shell
goose -dir ./migrations create migration_name sql
```
