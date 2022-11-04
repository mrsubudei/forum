package main

import "forum/internal/app"

func main() {
	app.Run()
	// react := entity.Reaction{Like: true, Date: "dewe"}
	// jsonBlob, _ := json.Marshal(react)
	// re := reaction{}
	// json.Unmarshal(jsonBlob, &re)
	// fmt.Println(re)
}

// type reaction struct {
// 	Like bool `json:"like"`
// 	// Date string `json:"reaction_date"`
// }
