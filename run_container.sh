#!/bin/sh

PROJECT_NAME="fast-links"
DEFAULT_PORT=8080

# Funkcja wyświetlająca sposób użycia
usage() {
    echo "Użycie: $0 [opcje]"
    echo "Opcje:"
    echo "  -p, --port PORT    Port na którym ma działać aplikacja (domyślnie: $DEFAULT_PORT)"
    echo "  -b, --build        Wymuś przebudowanie obrazu"
    echo "  -d, --down         Zatrzymaj i usuń kontenery"
    exit 1
}

# Domyślne wartości
BUILD_FLAG=""
PORT=$DEFAULT_PORT
ACTION="up"

# Przetwarzanie argumentów
ARGS=$(getopt -o p:bd --long port:,build,down -n "$0" -- "$@")
eval set -- "$ARGS"

while true; do
    case "$1" in
        -p|--port)
            PORT="$2"
            shift 2
            ;;
        -b|--build)
            BUILD_FLAG="--build"
            shift
            ;;
        -d|--down)
            ACTION="down"
            shift
            ;;
        --)
            shift
            break
            ;;
        *)
            usage
            ;;
    esac
done

# Eksportuj zmienne środowiskowe
export PROJECT_NAME
export PORT

# Wykonaj akcję
if [ "$ACTION" == "down" ]; then
    echo "Zatrzymywanie i usuwanie kontenerów..."
    docker compose down
else
    echo "Uruchamianie kontenera na porcie $PORT..."
    docker compose up -d $BUILD_FLAG
    
    # Wyświetl informacje o uruchomionym kontenerze
    echo "Kontener uruchomiony. Szczegóły:"
    docker compose ps
    
    echo "Aplikacja dostępna pod adresem: http://0.0.0.0:$PORT"
fi