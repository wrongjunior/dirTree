
# dirTree

## Описание

**dirTree** — это утилита для сканирования файловой системы с выводом дерева директорий. Программа поддерживает как текстовый интерфейс командной строки (CLI), так и интерактивный текстовый пользовательский интерфейс (TUI), который позволяет выбирать директории вручную. Утилита предоставляет возможность генерировать список файлов и папок с указанием относительных или абсолютных путей, а также исключать определённые директории или файлы с конкретными расширениями.

## Навигация
- [Пример вывода программы](#пример-вывода-программы)
- [Основные возможности](#основные-возможности)
- [Установка и запуск](#установка-и-запуск)
- [Параметры командной строки](#параметры-командной-строки)
- [Конфигурация](#конфигурация)
- [Интерактивный режим](#интерактивный-режим)
- [Зависимости](#зависимости)
- [Лицензия](#лицензия)

## Пример вывода программы

Пример 1: Вывод относительных путей с игнорированием директории `.git` и файлов с расширением `.mod`.

```
dirTree/
├── CONTRIBUTING.md (4.34 KB)
├── LICENSE (1.04 KB)
├── README.md (5.27 KB)
├── cmd/
│   └── main.go (808 B)
├── conf.txt (42 B)
└── internal/
    ├── cli/
    │   └── cli.go (723 B)
    ├── config/
    │   └── config.go (3.03 KB)
    ├── fileinfo/
    │   └── fileinfo.go (178 B)
    ├── output/
    │   └── output.go (3.76 KB)
    ├── scanner/
    │   └── scanner.go (5.04 KB)
    └── tui/
        └── tui.go (5.65 KB)
```

Пример 2: Вывод абсолютных путей без игнорирования директорий и файлов.

```
/Users/nameUser/Program/nameUserdirTree
/Users/nameUser/Program/nameUserdirTree/CONTRIBUTING.md
/Users/nameUser/Program/nameUserdirTree/LICENSE
/Users/nameUser/Program/nameUserdirTree/README.md
/Users/nameUser/Program/nameUserdirTree/cmd
.
.
/Users/nameUser/Program/nameUserdirTree/internal/tui
/Users/nameUser/Program/nameUserdirTree/internal/tui/tui.go
```

## Основные возможности

- Вывод дерева директорий.
- Интерактивный выбор директорий через TUI.
- Поддержка флагов для вывода относительных или абсолютных путей.
- Возможность игнорирования директорий и файлов по расширению.
- Копирование результатов в буфер обмена.

## Установка и запуск

Чтобы установить проект, необходимо клонировать репозиторий и скомпилировать программу с помощью Go:

```bash
git clone https://github.com/wrongjunior/dirTree.git
cd dirTree
go build -o dirTree main.go
```

## Параметры командной строки

### Командный интерфейс (CLI)

Пример команды для вывода относительных путей всех файлов и директорий:

```bash
./dirTree --relative
```


https://github.com/user-attachments/assets/42b938a0-c7e6-4ced-bf9c-0c9c5d03d03d



Пример команды для вывода абсолютных путей:

```bash
./dirTree --absolute
```



https://github.com/user-attachments/assets/04cdd0a1-ecda-4acc-9a1b-bc1c5a2a22b0




### Игнорирование директорий и файлов

Для игнорирования определённых директорий используется флаг `--ignore-dirs`:

```bash
./dirTree --ignore-dirs "vendor,.git"
```



https://github.com/user-attachments/assets/62e9034b-81ad-44fd-adfa-ff0fca7d98e2




Игнорирование файлов по расширению осуществляется с помощью флага `--ignore-exts`(вместе с игнорированием директорий):

```bash
./dirTree --ignore-exts "go,tmp"
```

<img width="377" alt="image" src="https://github.com/user-attachments/assets/f0fc701a-4cae-4a94-bcca-5aa1f0ee239c">


## Конфигурация

Для более гибкой настройки можно использовать файл конфигурации, в котором указываются игнорируемые директории и расширения файлов. Пример файла конфигурации(c теми же настройками что выше, по отдельности):

```
dir:vendor
dir:.git
dir:.idea
ext:go
ext:tmp
```

Запуск программы с указанием файла конфигурации:

```bash
./dirTree --ignore-config ./config.txt
```

<img width="377" alt="image" src="https://github.com/user-attachments/assets/8fdf7426-294d-4650-9260-69a18a1cc2d8">


## Интерактивный режим

Для запуска интерактивного режима выбора директорий используется команда:

```bash
./dirTree --tui
```



https://github.com/user-attachments/assets/92d0a92e-6107-4da8-a038-dc632e752aa5



## Зависимости

Для корректной работы проекта необходимы следующие зависимости:

- [Go](https://golang.org/doc/install) версии 1.22 или выше.
- Библиотеки, установленные через `go mod`.

## Лицензия

Проект распространяется под лицензией MIT. Подробности см. в файле [LICENSE](LICENSE).
