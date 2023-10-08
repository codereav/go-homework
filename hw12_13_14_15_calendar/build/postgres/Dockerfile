
# Используем образ PostgreSQL
FROM postgres

# Установка необходимых зависимостей
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        gcc

# Копирование файлов и скриптов
COPY ./migrations/ /docker-entrypoint-initdb.d/

ENV POSTGRES_HOST_AUTH_METHOD trust

# Запуск PostgreSQL
CMD ["postgres"]