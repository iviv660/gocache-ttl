# GoCache TTL 🗑️

In-memory кэш с ограничением по времени жизни записей (TTL) и максимальным размером (MaxSize).  
Реализован на Go с поддержкой автоматической очистки устаревших элементов.

## Возможности
- Установка ключа с TTL `Set`
- Получение значения по ключу (`Get`)
- Удаление ключа `Delete`
- Проверка существования (`Exists`)
- Получение всех ключей (`Keys`)
- Автоматический `cleanup` устаревших значений
- Ограничение по максимальному количеству записей (`MaxSize`, вытеснение старых)

## Структура проекта
```
├── cache.go # Реализация кэша 
├── cache_test.go # Набор unit-тестов 
└── README.md # Документация
```

## Установка и запуск тестов
```bash
git clone https://github.com/your-username/gocache-ttl.git
cd gocache-ttl
go test ./...
```

## Пример использованния
```bash
cache := gocachettl.NewCache(time.Second, 10)
defer cache.Close()

cache.Set("foo", "bar", 5*time.Second)

if val, ok := cache.Get("foo"); ok {
    fmt.Println(val) // bar
}
```




