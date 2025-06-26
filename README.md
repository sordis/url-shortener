# URL Shortener Service

[![CI Status](https://github.com/sordis/url-shortener/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/sordis/url-shortener/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/sordis/url-shortener/badge.svg?branch=master)](https://coveralls.io/github/sordis/url-shortener?branch=master)
![Go Version](https://img.shields.io/github/go-mod/go-version/sordis/url-shortener)

## Описание
Сервис для сокращения URL-адресов


### Локально
```bash
# Сборка
go build -o url-shortener ./cmd/url-shortener

# Запуск
CONFIG_PATH=./config/prod.yml AUTH_PASS=your_password ./url-shortener
```

### Через Docker
```bash
docker build -t url-shortener .
docker run -e AUTH_PASS=your_password -p 8080:8080 url-shortener
```

## Конфигурация
Основные параметры (config/prod.yml):
```yaml
http_server:
  address: ":8080"
  timeout: "4s"
  user: "admin"
  password: "${AUTH_PASS}"
storage:
  path: "./storage/storage.db"
```





## Лицензия
MIT License

