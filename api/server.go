package api

import (
	"akv/controller"
	"encoding/json"
	"fmt"
	"net/http"
)

func StartServer(controller *controller.PodController) {
	http.HandleFunc("/pods", func(w http.ResponseWriter, r *http.Request) {
		podList := controller.GetPods()
		writeJson(w, podList)

	})

	http.HandleFunc("/summary", func(w http.ResponseWriter, r *http.Request) {
		podSummary := controller.GetSummary()
		writeJson(w, podSummary)
	})
	fmt.Println("🌐 REST API running at http://localhost:8080")
	go http.ListenAndServe(":8080", nil)
}
func writeJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}
