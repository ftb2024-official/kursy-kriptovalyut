/*
project_root/
├── cmd/
│   └── server/                        # Главная точка входа в приложение (сервер)
│       └── main.go                    # Запуск приложения и настройка зависимостей
├── internal/
│   ├── domain/
│   │   └── coin.go                    # Определение сущности Coin
│   ├── usecases/
│   │   └── service.go                 # Реализация Service - бизнес-логика приложения
│   ├── ports/
│   │   ├── crypto_provider.go         # Интерфейс CryptoProvider
│   │   ├── storage.go                 # Интерфейс Storage
│   │   └── service.go                 # Интерфейс Service
│   ├── adapters/
│   │   ├── api/
│   │   │   └── server.go              # Сервер для REST API, реализация Server
│   │   ├── provider/
│   │   │   └── cryptocompare.go       # Адаптер CryptoProvider для интеграции с API CryptoCompare
│   │   └── storage/
│   │       └── postgres.go            # Адаптер Storage для PostgreSQL
├── config/
│   └── config.go                      # Конфигурация приложения (настройки БД, URL внешних API и т.д.)
└── pkg/
    └── utils/                         # Утилиты и вспомогательные функции
        └── http_client.go             # HTTP клиент для запросов к внешнему API

*/

package main

func main() {

}
