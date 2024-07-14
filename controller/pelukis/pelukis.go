package pelukis

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"audy/database"
	"audy/model/pelukis"

	"github.com/gorilla/mux"
)

func GetPelukis(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM pelukis")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pelukisList []pelukis.Pelukis
	for rows.Next() {
		var c pelukis.Pelukis
		if err := rows.Scan(&c.PelukisId, &c.Nama, &c.Alamat); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pelukisList = append(pelukisList, c)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pelukisList)
}

func PostPelukis(w http.ResponseWriter, r *http.Request) {
	var pc pelukis.Pelukis
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//prepare the SQL statement for inserting a new course
	query := `
        INSERT INTO pelukis (nama, alamat)
        VALUES (?, ?)`

	//Execute the SQL statement
	ress, err := database.DB.Exec(query, pc.Nama, pc.Alamat)
	if err != nil {
		http.Error(w, "Failed to insert pelukis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//Get the newly created ID
	id, _ := ress.LastInsertId()

	//Return the newly created ID in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Pelukis added successfully",
		"id":      id,
	})
}

func PutPelukis(w http.ResponseWriter, r *http.Request) {
	//AMBIL ID DARI URL
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	//DECODE JSON BODY
	var pc pelukis.Pelukis
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//PREPARE THE SQL STATEMENT FOR UPDATING THE CATEGORY ADMIN
	query := `
		UPDATE pelukis
		SET nama=?, alamat=?
		WHERE PelukisId=?`

	//EXECUTE THE SQL STATEMENT
	result, err := database.DB.Exec(query, pc.Nama, pc.Alamat, idStr)
	if err != nil {
		http.Error(w, "Failed to update pelukis: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	//GET THE NUMBER OF ROWS AFFECTED
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	//CHECK IF ANY ROWS WERE UPDATED
	if rowsAffected == 0 {
		http.Error(w, "No rows were updated", http.StatusNotFound)
		return
	}

	//RETURN SUCCES MESSAGE
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Course updated succesfully",
	})

}

func DeletePelukis(w http.ResponseWriter, r *http.Request) {
	//EXTRACT ID FROM URL
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

	//PREPARE THE SQL STATEMENT FOR DELETING A CATEGORY ADMIN
	query := `
		DELETE FROM pelukis
		WHERE PelukisId = ?`

	//EXECUTE THE SQL STATEMENT
	result, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete pelukis: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	//CHECK IF ANY ROWS WERE AFFECTED
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows were deleted", http.StatusNotFound)
		return
	}

	//RETURN THE RESPONSE
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mesage": "Pelukis deleted successfully",
	})

}

// func GetPelukisByID handles GET request to fetch a single pelukis by its ID
func GetPelukisByID(w http.ResponseWriter, r *http.Request) {
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

	var pelukis pelukis.Pelukis
	query := "SELECT * FROM pelukis WHERE PelukisId = ?"
	err = database.DB.QueryRow(query, id).Scan(&pelukis.PelukisId, &pelukis.Nama, &pelukis.Alamat)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Pelukis not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pelukis)
}
