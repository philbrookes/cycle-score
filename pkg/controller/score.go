package controller

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/strava/go.strava"
	"encoding/json"
	"strconv"
	"time"
	"fmt"
)

type score struct {
	Score int64 `json:"score"`
	Rides int `json:"rides"`
}

func ConfigureScore(router *mux.Router) {
	router.HandleFunc("/generate/from/{from}/to/{to}", generateScoreFromTo)
	router.HandleFunc("/generate/from/{from}", generateScoreFrom)
	router.HandleFunc("/generate", generateScore)
}

func generateScoreFromTo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	token, err := r.Cookie("strava_token")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	athleteId, err := r.Cookie("strava_athlete_id")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	client := strava.NewClient(token.Value, http.DefaultClient)
	activities, err := getActivities(client, athleteId.Value, params["from"], params["to"])

	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	results := calculateTotalScore(activities)

	w.Header().Add("Content-Type", "Application/JSON")
	json.NewEncoder(w).Encode(results)
}

func generateScoreFrom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	token, err := r.Cookie("strava_token")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	athleteId, err := r.Cookie("strava_athlete_id")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	client := strava.NewClient(token.Value, http.DefaultClient)
	activities, err := getActivities(client, athleteId.Value, params["from"], strconv.Itoa(int(time.Now().Unix())))

	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	results := calculateTotalScore(activities)

	w.Header().Add("Content-Type", "Application/JSON")
	json.NewEncoder(w).Encode(results)
}


func generateScore(w http.ResponseWriter,r *http.Request) {
	token, err := r.Cookie("strava_token")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	athleteId, err := r.Cookie("strava_athlete_id")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	client := strava.NewClient(token.Value, http.DefaultClient)
	activities, err := getActivities(client, athleteId.Value, "0", strconv.Itoa(int(time.Now().Unix())))

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	results := calculateTotalScore(activities)

	w.Header().Add("Content-Type", "Application/JSON")
	json.NewEncoder(w).Encode(results)
}


func getActivities(client *strava.Client, athleteStr, fromStr, toStr string) ([]*strava.ActivitySummary, error) {
	allActivities := make([]*strava.ActivitySummary, 0)
	athlete, err := strconv.ParseInt(athleteStr, 10, 64)
	if err != nil {
		return nil, err
	}
	from, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		return nil, err
	}
	to, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		return nil, err
	}


	pageOn := 1
	for {
		activities, err := strava.NewAthletesService(client).ListActivities(athlete).Page(pageOn).After(from).Before(to).Do()
		if err != nil {
			return allActivities, err
		}

		for _, activity := range activities {
			allActivities = append(allActivities, activity)
		}
		pageOn++

		if len(activities) == 0 {
			return allActivities, nil
		}
	}
}

func calculateTotalScore(activities []*strava.ActivitySummary) score {
	results := score{
		Score: 0,
		Rides: 0,
	}
	for _, activity := range activities {
		if activity.Type != strava.ActivityTypes.Ride {
			continue
		}
		results.Rides++
		results.Score += calculateActivityScore(activity)
	}

	return results
}

func calculateActivityScore(activity *strava.ActivitySummary) (int64) {
	return int64(activity.Distance + float64(activity.MovingTime) + activity.TotalElevationGain + activity.AverageSpeed)
}