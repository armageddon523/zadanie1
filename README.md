# Weather App – Zadanie 1

Prosta aplikacja webowa napisana w Go, uruchamiana w kontenerze Docker.  
Aplikacja pozwala wybrać miasto z listy i wyświetla aktualną pogodę dla wybranej lokalizacji.

Projekt został przygotowany zgodnie z wymaganiami zadania: aplikacja zapisuje w logach datę uruchomienia, autora oraz port TCP, a Dockerfile wykorzystuje wieloetapowe budowanie obrazu, obraz końcowy `scratch`, etykiety OCI i healthcheck.

## Technologie

- Go
- Docker
- Open-Meteo API

## Pliki w projekcie

- `main.go` – kod aplikacji
- `go.mod` – konfiguracja modułu Go
- `Dockerfile` – plik do budowania obrazu Docker

## Opis działania aplikacji

Po uruchomieniu aplikacja działa jako prosty serwer HTTP na porcie `8080`.

W przeglądarce użytkownik wybiera lokalizację z listy, a aplikacja pobiera aktualne dane pogodowe, takie jak:

- temperatura,
- prędkość wiatru,
- kierunek wiatru,
- kod pogody,
- czas pomiaru.

Dane pogodowe pobierane są z publicznego API Open-Meteo.

## Budowanie obrazu

```bash
docker build -t weatherapp:go .
