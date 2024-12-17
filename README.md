# Проект: CRUD-приложение с PostgreSQL и Redis

## Описание
Этот проект представляет собой приложение, использующее **PostgreSQL** для хранения данных и **Redis** для кэширования. Приложение развёртывается с использованием **Docker** и **Docker Compose**.

## Требования

Перед запуском проекта на вашей машине необходимо установить:

1. **Docker**  
   Скачать и установить Docker для вашей ОС.  
   Следуйте инструкциям для установки с официального сайта:  
   [Docker Desktop](https://www.docker.com/products/docker-desktop)

2. **Docker Compose**  
   Docker Compose включён в Docker Desktop по умолчанию.

3. **Postman**  
   Для удобства тестирования API и добавления данных скачайте и установите **Postman**:  
   [Postman Download](https://www.postman.com/downloads/)  

---

## Запуск проекта

1. **Клонируйте репозиторий**:
   ```bash
   git clone <ссылка-на-репозиторий>
   cd <папка-с-проектом>

2. **Запустите проект с помощью Docker Compose**:
   В корневой директории проекта выполните команду:
   docker-compose up --build

3. **Доступ к приложению**:
   После запуска проект будет доступен по адресу:
   http://localhost:8080

4. **Остановка проекта**:
   Для остановки всех контейнеров выполните:
   docker-compose down
 
5. **Сервисы проекта**:
   app: Основное приложение.
   db: PostgreSQL база данных.
   redis: Кэш Redis.
   migrate: Выполнение миграций для базы данных.
   golangci-lint: Статический анализ кода.

## Дополнительные команды

1. **Просмотр логов**:
   docker-compose logs

2. **Перезапуск контейнеров**:
   docker-compose restart

3. **Перезапуск контейнеров**:
   docker-compose run migrate

---
  
## Тестирование CRUD API с Postman

В этом разделе описаны шаги для тестирования CRUD API с использованием инструмента **Postman**.

### 1. Создание записи (Create)

**URL:** `http://localhost:8080/book/create`  
**Метод:** POST  
**Тело запроса (JSON):**

{
    "title": "Users",
    "author": "Dia Doe",
    "published_at": "2024-08-01T00:00:00Z"
}


### 2. Получение всех записей (Read)

**URL:** `http://localhost:8080/books`  
**Метод:** GET  


### 3. Получение записи по ID (ReadById)

**URL:** `http://localhost:8080/book/{id}`  
**Метод:** GET  

Замените {id} на нужный идентификатор книги, чтобы получить данные о конкретной книге. 
Например: GET http://localhost:8080/book/1


### 4. Обновление записи (Update)

**URL:** `http://localhost:8080/book/update/{id}`  
**Метод:** PUT 
**Тело запроса (JSON):**

{
    "title": "Users",
    "author": "Dia Doe",
    "published_at": "2024-08-01T00:00:00Z"
}

Замените {id} на нужный идентификатор книги, чтобы изменить данные. 
Например: PUT http://localhost:8080/book/update/1


### 5. Удаление записи (Delete)

**URL:** `http://localhost:8080/book/delete/{id}`  
**Метод:** DELETE

Замените {id} на нужный идентификатор книги, чтобы удалить её из базы данных. 
Например: DELETE http://localhost:8080/book/delete/1

