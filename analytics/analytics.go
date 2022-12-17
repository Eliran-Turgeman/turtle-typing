package analytics

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

var (
	GameDataFile = os.Getenv("HOME") + "/.turtle-typing/game_data.csv"
)

func SaveGameData(wpm int, errorRate int) {
	data := []string{time.Now().String(), fmt.Sprint(wpm), fmt.Sprint(errorRate)}
	fmt.Println("data:", data)
	csvFile, err := os.Open(GameDataFile)

	if err != nil {
		fmt.Printf("failed opening file file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)
	err = csvwriter.Write(data)
	if err != nil {
		fmt.Printf("Unable to write data: %s", err)
	}
	csvwriter.Flush()
	csvFile.Close()
}
