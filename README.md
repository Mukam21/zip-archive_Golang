#  Zip Archive Golang Service

## Описание

Это REST API-сервис на Go, который позволяет:

- Создавать задачи на скачивание файлов с
публичных URL

- Добавлять ссылки к задаче (до 3 файлов)

- Валидировать тип файлов (только .pdf, .jpeg, .jpg)

- Скачивать и упаковывать файлы в `.zip` архив

- Возвращать архив (в base64) в теле ответа

- Сообщать об ошибках для недоступных или запрещённых файлов

- Обрабатывать не более 3 задач одновременно

---

## Технологии

- Язык: Go 1.20+

- Роутинг: `chi`

- Архивация: стандартный `archive/zip`

- UUID, MIME, Base64 — из стандартной библиотеки

---

## Структура проекта

---

## Запуск

```bash

      go run cmd/main.go


API Эндпоинты:

Метод	        URL	                    Описание

POST	        /task	              Создать новую задачу

POST	        /task/{id}/url	        Добавить файл по ссылке

GET	        /task/{id}/status	  Получить статус + архив



Пошаговая проверка через Postman или curl

1.Создать задачу

      POST http://localhost:8080/task

Через Postman:

Method: POST

Body: пусто

Через терминал:

      curl -X POST http://localhost:8080/task

Ответ:

      {
        "id": "80f6754d-01ec-4758-b1f3-ee2c91b99982"
      }

2️⃣ Добавить файлы (до 3-х)

Каждый файл добавляется отдельным запросом:

      POST http://localhost:8080/task/{id}/url

Примеры тел запроса (JSON):

Рабочее изображение:

      {
        "url": "https://upload.wikimedia.org/wikipedia/commons/3/3f/JPEG_example_flower.jpg"
      }

Второй файл (недоступен):

      {
        "url": "https://www.africau.edu/images/default/sample.pdf"
      }

Третий файл (недоступен):

      {
        "url": "https://www.orimi.com/pdf-test.pdf"
      }

Используй 3 отдельных запроса подряд. После третьего автоматически начнётся скачивание и архивация.

Получить статус задачи

      GET http://localhost:8080/task/{id}/status

Пример ответа:

До завершения:

      {
        "status": "pending"
      }

После завершения:

      {
        "status": "completed",
        "errors": [
          {
            "url": "https://www.africau.edu/images/default/sample.pdf",
            "message": "download failed"
          },
          {
            "url": "https://www.orimi.com/pdf-test.pdf",
            "message": "download failed"
          }
        ],
        "archive": "UEsDBBQACAAIAAAA..." // base64 zip
      }


Распаковка архива из base64

Скопируй поле "archive" в файл archive.zip:

Windows PowerShell:

      [System.Convert]::FromBase64String("ВАШ_КОД") | Set-Content -Encoding Byte archive.zip

Linux / macOS:

      echo "ВАШ_КОД" | base64 -d > archive.zip

 Автор
      Мукам Усманов

GitHub:

      https://github.com/Mukam21/zip-archive_Golang
