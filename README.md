# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## Структура проекта

Приведённая в этом репозитории структура проекта является рекомендуемой, но не обязательной.

Это лишь пример организации кода, который поможет вам в реализации сервиса.

При необходимости можно вносить изменения в структуру проекта, использовать любые библиотеки и предпочитаемые структурные паттерны организации кода приложения, например:
- **DDD** (Domain-Driven Design)
- **Clean Architecture**
- **Hexagonal Architecture**
- **Layered Architecture**

## Результаты профилирования памяти

После оптимизации инициализации логгера, мы получили следующие результаты, которые показывают уменьшение потребления памяти:

```
Showing nodes accounting for -242.13kB, 18.21% of 1329.76kB total
      flat  flat%   sum%        cum   cum%
 -768.26kB 57.77% 57.77%  -768.26kB 57.77%  go.uber.org/zap/zapcore.newCounters (inline)
  526.13kB 39.57% 18.21%   526.13kB 39.57%  github.com/jackc/pgx/v5/pgtype.(*Map).buildReflectTypeToType (inline)
         0     0% 18.21%  -768.26kB 57.77%  github.com/RussiaFPS/shortlink/internal/logger.NewLogger
         0     0% 18.21%   526.13kB 39.57%  github.com/jackc/pgx/v5.ConnectConfig
         0     0% 18.21%   526.13kB 39.57%  github.com/jackc/pgx/v5.connect
         0     0% 18.21%   526.13kB 39.57%  github.com/jackc/pgx/v5/pgtype.NewMap
         0     0% 18.21%   526.13kB 39.57%  github.com/jackc/pgx/v5/pgtype.initDefaultMap
         0     0% 18.21%   526.13kB 39.57%  github.com/jackc/pgx/v5/pgxpool.NewWithConfig.func1
         0     0% 18.21%   526.13kB 39.57%  github.com/jackc/puddle/v2.(*Pool[go.shape.*uint8]).initResourceValue.func1
         0     0% 18.21%  -768.26kB 57.77%  go.uber.org/zap.(*Logger).WithOptions
         0     0% 18.21%  -768.26kB 57.77%  go.uber.org/zap.Config.Build
         0     0% 18.21%  -768.26kB 57.77%  go.uber.org/zap.Config.buildOptions.WrapCore.func5
         0     0% 18.21%  -768.26kB 57.77%  go.uber.org/zap.Config.buildOptions.func1
         0     0% 18.21%  -768.26kB 57.77%  go.uber.org/zap.New
         0     0% 18.21%  -768.26kB 57.77%  go.uber.org/zap.optionFunc.apply
         0     0% 18.21%  -768.26kB 57.77%  go.uber.org/zap/zapcore.NewSamplerWithOptions
         0     0% 18.21%  -768.26kB 57.77%  main.main
         0     0% 18.21%  -768.26kB 57.77%  runtime.main
         0     0% 18.21%   526.13kB 39.57%  sync.(*Once).Do (inline)
         0     0% 18.21%   526.13kB 39.57%  sync.(*Once).doSlow
```
