package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Match struct {
	ID            int       `json:"id"`
	BannedMaps    []string  `json:"banned_maps"`
	CapA          int       `json:"cap_a"`
	CapB          int       `json:"cap_b"`
	CurrentTurn   int       `json:"current_turn"`
	TurnCount     int       `json:"turn_count"`
	TurnStartAt   time.Time `json:"turn_start_at"`
}

func main() {
	connStr := "user=u0_a199 dbname=faceit_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil { log.Fatal(err) }
	defer db.Close()

	// Функция выполнения бана
	executeBan := func(mapName string, m Match) {
		newTurnCount := m.TurnCount - 1
		newTurnID := m.CurrentTurn
		resetTimer := false

		if newTurnCount <= 0 {
			resetTimer = true
			if m.CurrentTurn == m.CapA {
				newTurnID = m.CapB
				newTurnCount = 2
			} else {
				newTurnID = m.CapA
				newTurnCount = 2
			}
		}

		query := "UPDATE matches SET banned_maps = array_append(banned_maps, $1), current_turn_id = $2, turn_count = $3"
		if resetTimer {
			query += ", turn_start_at = CURRENT_TIMESTAMP"
		}
		query += " WHERE id = 1"
		db.Exec(query, mapName, newTurnID, newTurnCount)
	}

	http.HandleFunc("/get-match", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		var m Match
		var banned pq.StringArray
		err := db.QueryRow("SELECT id, banned_maps, cap_a_id, cap_b_id, current_turn_id, turn_count, turn_start_at FROM matches WHERE id = 1").
			Scan(&m.ID, &banned, &m.CapA, &m.CapB, &m.CurrentTurn, &m.TurnCount, &m.TurnStartAt)
		
		if err == nil {
			m.BannedMaps = banned
			// Автобан через 30 секунд
			if time.Since(m.TurnStartAt).Seconds() > 30 && len(m.BannedMaps) < 8 {
				allMaps := []string{"Mirage", "Inferno", "Dust II", "Ancient", "Anubis", "Vertigo", "Nuke", "Overpass", "Train"}
				for _, am := range allMaps {
					isAlreadyBanned := false
					for _, bm := range m.BannedMaps { if am == bm { isAlreadyBanned = true } }
					if !isAlreadyBanned {
						executeBan(am, m)
						fmt.Printf("⏰ AUTO-BAN: %s\n", am)
						break
					}
				}
			}
		}
		json.NewEncoder(w).Encode(m)
	})

	http.HandleFunc("/ban", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		mapName := r.URL.Query().Get("map")
		userID := r.URL.Query().Get("user_id")

		var m Match
		var banned pq.StringArray
		db.QueryRow("SELECT cap_a_id, cap_b_id, current_turn_id, turn_count, banned_maps FROM matches WHERE id = 1").
			Scan(&m.CapA, &m.CapB, &m.CurrentTurn, &m.TurnCount, &banned)

		if fmt.Sprintf("%d", m.CurrentTurn) != userID {
			http.Error(w, "Wrong turn", 403); return
		}

		executeBan(mapName, m)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		db.Exec("UPDATE matches SET banned_maps = '{}', current_turn_id = 1, turn_count = 3, turn_start_at = CURRENT_TIMESTAMP WHERE id = 1")
		w.Write([]byte("OK"))
	})

	fmt.Println("🚀 Бэкенд запущен! Автобан (30с) активен.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
