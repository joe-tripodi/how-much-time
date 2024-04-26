package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
  life = 72
  minutesPerYear = 525600 // 365*24*60
)

func renderTemplate(w http.ResponseWriter, tmpl string, data any) {
  t, err := template.ParseFiles(tmpl + ".html")
  if err != nil {
    log.Fatal("There was an error", err)
  }
  t.Execute(w, data)
}

func handleIndex(w http.ResponseWriter, r *http.Request){
  log.Print("User has found our home page")
  log.Print(r.URL)
  renderTemplate(w, "home", nil)
  return
}

func handleNaughty(w http.ResponseWriter, r *http.Request){
  log.Print("User was trying to be naughty")
  http.Redirect(w, r, "home", http.StatusFound)
  return
}

type UserRequest struct {
  Activity string
  Duration int
  Age int
}

func parseUserInput(r *http.Request) *UserRequest {

  duration, _ := strconv.Atoi(r.FormValue("duration"))
  age, _ := strconv.Atoi(r.FormValue("age"))

  ur := UserRequest{
    Activity: r.FormValue("activity"),
    Duration: duration,
    Age: age,
  }

  return &ur
}

func (ur UserRequest) String() string {
  return fmt.Sprintf(`{Activity: "%s", Duration: "%d", Age: "%d"}`, ur.Activity, ur.Duration, ur.Age)
}

func calculatePercentageOfLife(activityTime, yearsLeft, minutesPerYear int) float32 {
  var ret float32 = 0.0
  return ret
}

func handleCalculate(w http.ResponseWriter, r *http.Request){
  log.Print("User wants us to calculate")

  ur := parseUserInput(r)
  log.Print(ur)
  yearsLeft := life - ur.Age
  minutesLeft := yearsLeft*minutesPerYear
  activityTime := yearsLeft*365*ur.Duration
  sleepTime := float32(yearsLeft*(365*8*60))

  p := message.NewPrinter(language.BritishEnglish)

  p.Fprintf(w, "You have %d years left to live\n", yearsLeft)
  p.Fprintf(w, "Which means you have %d minutes left to live.\n", minutesLeft)
  p.Fprintf(w, "You will spend %d minutes %s", activityTime, ur.Activity)
  p.Fprintf(w, "Which is %.2f%% of your remaining life.\n", float32(float32(activityTime)/float32((minutesLeft)))*100.00)
  p.Fprintf(w, "Or %.2f%% of your remaining WAKING life.", float32(float32(activityTime)/(float32(minutesLeft)-sleepTime))*100.00)

  return
}

func handleFavicon(w http.ResponseWriter, r *http.Request){
  return
}

func main(){
  log.SetFlags(0)
  log.SetPrefix("hmt-server: ")

  http.HandleFunc("/*/*", handleNaughty)
  http.HandleFunc("/", handleIndex)
  http.HandleFunc("/calculate", handleCalculate)
  http.HandleFunc("/favicon.ico", handleFavicon)
  log.Fatal(http.ListenAndServe(":8888", nil))
}
