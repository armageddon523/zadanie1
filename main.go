package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

const studentName = "PAWEŁ LADOWSKI"

type Place struct {
	Label  string
	Region string
	Lat    float64
	Lon    float64
}

type ForecastResponse struct {
	Current struct {
		Temperature float64 `json:"temperature"`
		Windspeed   float64 `json:"windspeed"`
		Winddir     float64 `json:"winddirection"`
		Weathercode int     `json:"weathercode"`
		Time        string  `json:"time"`
	} `json:"current_weather"`
}

type ViewData struct {
	Places     map[string]Place
	Selected   Place
	SelectedID string
	Weather    *ForecastResponse
}

var cityList = map[string]Place{
	"los-angeles": {"Los Angeles", "USA", 34.0522, -118.2437},
	"new-york":    {"Nowy Jork", "USA", 40.7128, -74.0060},
	"warsaw":      {"Warszawa", "Polska", 52.2297, 21.0122},
	"vigolo":      {"Vigolo", "Włochy", 45.7167, 10.0333},
	"new-delhi":   {"New Delhi", "Indie", 28.6139, 77.2090},
	"hawaii":      {"Hawaje / Honolulu", "USA", 21.3069, -157.8583},
}

var layout = template.Must(template.New("weather").Parse(`
<!doctype html>
<html lang="pl">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Panel pogody</title>
<style>
*{box-sizing:border-box}
body{
	margin:0;
	min-height:100vh;
	font-family:Inter,Arial,sans-serif;
	background:linear-gradient(135deg,#151515,#26324a 55%,#0e7490);
	color:#f8fafc;
	display:flex;
	align-items:center;
	justify-content:center;
	padding:28px;
}
.card{
	width:100%;
	max-width:820px;
	background:rgba(255,255,255,.10);
	border:1px solid rgba(255,255,255,.18);
	border-radius:28px;
	padding:34px;
	box-shadow:0 24px 80px rgba(0,0,0,.35);
	backdrop-filter:blur(14px);
}
.header{
	display:flex;
	justify-content:space-between;
	gap:20px;
	align-items:flex-start;
	margin-bottom:26px;
}
h1{
	margin:0;
	font-size:38px;
	letter-spacing:-1px;
}
.subtitle{
	margin-top:8px;
	color:#cbd5e1;
	line-height:1.5;
}
.badge{
	background:#22d3ee;
	color:#083344;
	padding:10px 14px;
	border-radius:999px;
	font-weight:700;
	white-space:nowrap;
}
form{
	display:grid;
	grid-template-columns:1fr auto;
	gap:12px;
	margin-top:22px;
}
select,button{
	border:0;
	border-radius:16px;
	padding:15px 16px;
	font-size:16px;
}
select{
	background:#f8fafc;
	color:#0f172a;
}
button{
	background:#facc15;
	color:#1c1917;
	font-weight:800;
	cursor:pointer;
}
button:hover{filter:brightness(.96)}
.weather-box{
	margin-top:28px;
	display:grid;
	grid-template-columns:repeat(2,1fr);
	gap:14px;
}
.tile{
	background:rgba(15,23,42,.55);
	border:1px solid rgba(255,255,255,.12);
	border-radius:20px;
	padding:18px;
}
.tile small{
	display:block;
	color:#94a3b8;
	margin-bottom:8px;
}
.tile strong{
	font-size:26px;
}
.place-title{
	margin-top:28px;
	font-size:24px;
	font-weight:800;
}
@media(max-width:650px){
	form,.weather-box{grid-template-columns:1fr}
	.header{display:block}
	.badge{display:inline-block;margin-top:16px}
	h1{font-size:30px}
}
</style>
</head>
<body>
<section class="card">
	<div class="header">
		<div>
			<h1>Sprawdź aktualną pogodę</h1>
			<p class="subtitle">Wybierz jedną z dostępnych lokalizacji i zobacz bieżące dane pogodowe.</p>
		</div>
		<div class="badge">PAwChO</div>
	</div>

	<form method="POST" action="/check">
		<select name="place">
			{{range $id, $p := .Places}}
<option value="{{$id}}" {{if eq $id $.SelectedID}}selected{{end}}>
	{{$p.Region}} — {{$p.Label}}
</option>
{{end}}
		</select>
		<button type="submit">Pokaż pogodę</button>
	</form>

	{{if .Weather}}
	<div class="place-title">{{.Selected.Region}} — {{.Selected.Label}}</div>
	<div class="weather-box">
		<div class="tile">
			<small>Temperatura</small>
			<strong>{{printf "%.1f" .Weather.Current.Temperature}}°C</strong>
		</div>
		<div class="tile">
			<small>Prędkość wiatru</small>
			<strong>{{printf "%.1f" .Weather.Current.Windspeed}} km/h</strong>
		</div>
		<div class="tile">
			<small>Kierunek wiatru</small>
			<strong>{{printf "%.0f" .Weather.Current.Winddir}}°</strong>
		</div>
		<div class="tile">
			<small>Kod pogody</small>
			<strong>{{.Weather.Current.Weathercode}}</strong>
		</div>
		<div class="tile" style="grid-column:1/-1">
			<small>Czas pomiaru</small>
			<strong>{{.Weather.Current.Time}}</strong>
		</div>
	</div>
	{{end}}
</section>
</body>
</html>
`))

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		checkHealth()
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", showHome)
	mux.HandleFunc("/check", showWeather)
	mux.HandleFunc("/health", healthEndpoint)

	log.Printf("Start aplikacji: %s | Autor: %s | Port TCP: %s",
		time.Now().Format(time.RFC3339), studentName, appPort())

	if err := http.ListenAndServe(":"+appPort(), mux); err != nil {
		log.Fatal(err)
	}
}

func appPort() string {
	value := os.Getenv("PORT")
	if value == "" {
		return "8080"
	}
	return value
}

func showHome(w http.ResponseWriter, r *http.Request) {
	render(w, ViewData{Places: cityList})
}

func showWeather(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("place")

	place, exists := cityList[id]
	if !exists {
		http.Error(w, "Nie wybrano poprawnej lokalizacji", http.StatusBadRequest)
		return
	}

	weather, err := downloadWeather(place)
	if err != nil {
		http.Error(w, "Nie udało się pobrać danych pogodowych", http.StatusBadGateway)
		return
	}

	render(w, ViewData{
	Places:     cityList,
	Selected:   place,
	SelectedID: id,
	Weather:    weather,
})
}

func downloadWeather(place Place) (*ForecastResponse, error) {
	address := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true&timezone=auto",
		place.Lat,
		place.Lon,
	)

	response, err := http.Get(address)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var result ForecastResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func render(w http.ResponseWriter, data ViewData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := layout.Execute(w, data); err != nil {
		http.Error(w, "Błąd renderowania strony", http.StatusInternalServerError)
	}
}

func healthEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func checkHealth() {
	resp, err := http.Get("http://127.0.0.1:" + appPort() + "/health")
	if err != nil || resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}
}
