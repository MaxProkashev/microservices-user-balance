# microservices-user-balance
 console work with user balance
 
Язык написания: GoLang

База данных: Postgres (с помощью облачной PaaS-платформы Heroku) + драйвер: <a href="https://github.com/lib/pq">github.com/lib/pq</a>

### Таблица all_users (таблица всех пользователей)
|        id        | balance |
|:----------------:|:-------:|
| TEXT PRIMARY KEY | NUMERIC |

 - balance в рублях

### Таблица transactions (таблица всех транзакций)
|        id        | idFrom |   sum   | date | info |
|:----------------:|:------:|:-------:|:----:|:----:|
| TEXT PRIMARY KEY |  TEXT  | NUMERIC |  INT | TEXT |
 
 - id того, с чьим счетом была произведена операция
 - idFrom кто произвел операцию
 - sum в рублях
 - date в time.Time.Unix()
 - info краткая информация об операции
 
## Реализация API методов (маршрутизатор запросов: <a href="https://github.com/gorilla/mux">gorilla/mux</a>)
Везде реализован метод GET

#### - ```http://localhost:9000/user/{id:[0-9]+}/{act:(?:add|del)}?{sum:[0-9]+}```
Принимает переменные id,act,sum (проверяет что sum не ноль). Если del то проверка что на счете достаточно средств для списания такой sum.

#### - ```http://localhost:9000/user/transfer?{sum:[0-9]+};{idFrom:[0-9]+};{idTo:[0-9]+}```
Принимает переменные sum - скольок перевести в рублях,idFrom - от кого перевести,idTo - кому перевести. Проверка что тот с кого списывает имеет достаточно денег

#### - ```http://localhost:9000/user/{id:[0-9]+}/balance```
Возвращает баланс пользователя в рублях

#### - ```http://localhost:9000/user/{id:[0-9]+}/balance?{currency}```
Возвращает баланс пользователя в валюте currency 

#### - ```http://localhost:9000/user/{id:[0-9]+}/transactions?{sort:last|new|high|low}```
Возвращает пронумерованный список транзакция пользователя id
Сортировка:
 - last: По убыванию давности транзакции
 - new: По возрастанию давности транзакции
 - heigh: По убыванию суммы
 - low: По возрастанию суммы
 
## Способ запуска

клонировать репу --> в консоле запустить команду **go run *.go**

##### Пример запросов и ответов клиента: 

```
>>>curl -X GET http://localhost:9000/user/1/add?sum=1000
{"id":"1","balance":1000,"currency":"RUB"}

>>>curl -X GET http://localhost:9000/user/1/add?sum=200
{"id":"1","balance":1200,"currency":"RUB"}

>>>curl -X GET http://localhost:9000/user/1/del?sum=2000
608 Insufficient funds for card transaction

>>>curl -X GET http://localhost:9000/user/1/del?sum=400
{"id":"1","balance":800,"currency":"RUB"}

C:\Users\pmpav>curl -X GET http://localhost:9000/user/transfer?sum=500,idFrom=1,idTo=2
404 page not found

>>>curl -X GET http://localhost:9000/user/transfer?sum=500;idFrom=1;idTo=2
{"id":"1","balance":300,"currency":"RUB"}
{"id":"2","balance":500,"currency":"RUB"}

>>>curl -X GET http://localhost:9000/user/1/balance
{"id":"1","balance":300,"currency":"RUB"}

>>>curl -X GET http://localhost:9000/user/1/balance?currency=EUR
{"id":"1","balance":3.4258072800000003,"currency":"EUR"}

>>>curl -X GET http://localhost:9000/user/1/transactions?sort=high
[{"num":0,"operation":{"idFrom":"1","sum":1000,"date":1597928231,"info":"user has deposited money"}},{"num":1,"operation":{"idFrom":"2","sum":500,"date":1597928291,"info":"user transferred money"}},{"num":2,"operation":{"idFrom":"1","sum":400,"date":1597928254,"info":"user has withdrawn money from the account"}},{"num":3,"operation":{"idFrom":"1","sum":200,"date":1597928237,"info":"user has deposited money"}}]

>>>curl -X GET http://localhost:9000/user/1/transactions?sort=last
[{"num":0,"operation":{"idFrom":"1","sum":1000,"date":1597928231,"info":"user has deposited money"}},{"num":1,"operation":{"idFrom":"1","sum":200,"date":1597928237,"info":"user has deposited money"}},{"num":2,"operation":{"idFrom":"1","sum":400,"date":1597928254,"info":"user has withdrawn money from the account"}},{"num":3,"operation":{"idFrom":"2","sum":500,"date":1597928291,"info":"user transferred money"}}]
```

##### пример лога:

```
$ go run *.go
2020/08/20 15:57:05 [OK] Drop all_users table
2020/08/20 15:57:05 [OK] Drop transactions table
2020/08/20 15:57:05 [OK] Create all_users table
2020/08/20 15:57:05 [OK] Create transactions table
2020/08/20 15:57:11 [OK] Add new user id=1
2020/08/20 15:57:11 [OK] Update balance for user id=1 to balance=1000.00
2020/08/20 15:57:11 [OK] Add new operation
2020/08/20 15:57:17 [OK] Update balance for user id=1 to balance=1200.00
2020/08/20 15:57:18 [OK] Add new operation
2020/08/20 15:57:34 [OK] Update balance for user id=1 to balance=800.00
2020/08/20 15:57:34 [OK] Add new operation
2020/08/20 15:58:10 [OK] Add new user id=2
2020/08/20 15:58:10 [OK] Update balance for user id=1 to balance=300.00
2020/08/20 15:58:11 [OK] Update balance for user id=2 to balance=500.00
2020/08/20 15:58:11 [OK] Add new operation
2020/08/20 15:58:11 [OK] Add new operation
```
