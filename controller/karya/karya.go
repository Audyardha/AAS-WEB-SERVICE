package karya

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"audy/database"
	"audy/model/karya"

	"github.com/gorilla/mux"
)

func GetKarya(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM karya")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var karyaList []karya.Karya
	for rows.Next() {
		var c karya.Karya
		if err := rows.Scan(&c.KaryaId, &c.Judul, &c.PelukisId, &c.TahunDibuat, &c.Media); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		karyaList = append(karyaList, c)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(karyaList)
}

func PostKarya(w http.ResponseWriter, r *http.Request) {
	var pc karya.Karya
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for inserting a new karya
	query := `
		INSERT INTO karya (judul, PelukisId, TahunDibuat, media)
		VALUES (?, ?, ?, ?)`

	// Execute the SQL statement
	res, err := database.DB.Exec(query, pc.Judul, pc.PelukisId, pc.TahunDibuat, pc.Media)
	if err != nil {
		http.Error(w, "Failed to insert karya: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the last inserted ID
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly created ID in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Karya added successfully",
		"id":      id,
	})
}

func PutKarya(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Decode JSON body
	var pc karya.Karya
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for updating the karya
	query := `
		UPDATE karya
		SET judul=?, PelukisId=?, TahunDibuat=?, media=?
		WHERE KaryaId=?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, pc.Judul, pc.PelukisId, pc.TahunDibuat, pc.Media, id)
	if err != nil {
		http.Error(w, "Failed to update karya: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were updated
	if rowsAffected == 0 {
		http.Error(w, "No rows were updated", http.StatusNotFound)
		return
	}

	// Return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Karya updated successfully",
	})
}

func DeleteKarya(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for deleting a karya
	query := `DELETE FROM karya WHERE KaryaId = ?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete karya: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows were deleted", http.StatusNotFound)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Karya deleted successfully",
	})
}

// GetKaryaByID handles GET request to fetch a sigle karya by its ID
func GetKaryaByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var karya karya.Karya
	query := "SELECT * FROM karya WHERE KaryaId = ?"
	err = database.DB.QueryRow(query, id).Scan(&karya.KaryaId, &karya.Judul, &karya.PelukisId, &karya.TahunDibuat, &karya.Media)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Karya not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(karya)
}
