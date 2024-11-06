#!/bin/sh

# Author:Mateusz Kędziora https://mkedziora.pl

# Nazwa projektu
PROJECT_NAME="fast-links"

# Domyślny port
DEFAULT_PORT=8080

# Funkcja wyświetlająca sposób użycia
usage() {
    echo "Użycie: $0 [-p port]"
    echo "  -p port    Port na którym ma działać aplikacja (domyślnie: $DEFAULT_PORT)"
    exit 1
}

# Przetwarzanie argumentów
while getopts ":p:" opt; do
    case ${opt} in
        p )
            PORT=$OPTARG
            ;;
        \? )
            usage
            ;;
    esac
done

# Jeśli port nie został podany, użyj domyślnego
PORT=${PORT:-$DEFAULT_PORT}

# Zatrzymaj i usuń stary kontener, jeśli istnieje
echo "Zatrzymywanie i usuwanie starego kontenera..."
docker stop $PROJECT_NAME 2>/dev/null
docker rm $PROJECT_NAME 2>/dev/null

# Zbuduj nowy obraz
echo "Budowanie nowego obrazu..."
docker build -t $PROJECT_NAME .

# Uruchom nowy kontener
echo "Uruchamianie nowego kontenera na porcie $PORT..."
docker run -d --name $PROJECT_NAME -p $PORT:8080 $PROJECT_NAME

# Wyświetl informacje o uruchomionym kontenerze
echo "Kontener uruchomiony. Szczegóły:"
docker ps --filter name=$PROJECT_NAME

echo "Aplikacja dostępna pod adresem: http://localhost:$PORT"