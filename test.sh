#!/bin/bash
echo "Запуск сервера..."
nohup go run cmd/xkcd/main.go &> server.log &
SERVER_PID=$!
echo "Update..."
# Даем серверу немного времени для запуска
sleep 5

# Ожидание завершения обновления базы данных
echo "Ожидание завершения обновления базы данных..."
while ! grep -q "The update is finished." server.log; do
  sleep 1
done

# URL для запроса
URL="http://localhost:4000/pics?search=tst%20apple%20doctor"

# Выполнение запроса и сохранение ответа
RESPONSE=$(curl -s -w "%{http_code}" -o response.json "$URL")

# Извлечение HTTP-кода из ответа
HTTP_CODE=${RESPONSE: -3}

# Проверка кода ответа и содержимого
if [ "$HTTP_CODE" -ne 200 ]; then
  echo "Тест провален: получен HTTP-код $HTTP_CODE"
  sleep 5
  exit 1
fi

# Проверка содержимого ответа (пример проверки на наличие ключевого слова в JSON)
if grep -q '"https://imgs.xkcd.com/comics/an_apple_a_day.png"' response.json; then
  echo "Тест пройден: найдено ключевое слово 'apple doctor' в ответе"
  sleep 5
  exit 0
else
  echo "Тест провален: ключевое слово 'apple doctor' не найдено в ответе"
  sleep 5
  exit 1
fi
PID=$(lsof -t -i:4000)
if [ -n "$PID" ]; then
  kill $PID
  echo "Процесс с PID $PID был завершен"
else
  echo "Процесс на порту 4000 не найден"
fi

kill $SERVER_PID
