# frf-via-stats

Скрипт для сбора статистики по источникам постов (via) в Clio-архивах.

Установка:
```
go get github.com/davidmz/frf-via-stats
```

Использование:
```
frf-via-stats clio/*.zip > stats.jsons
```

Результат: в файле stats.jsons каждая строка — JSON-объект вида `{username, sources: [{name, url, count}]}`.