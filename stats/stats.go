package stats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	statisticsFile = "statistics.json"
)

type Statistics struct {
	Players map[string]int `json:"players"`
}

var (
	StatisticsData = Statistics{Players: make(map[string]int)}
)

func LoadStatistics() {
	fmt.Println("Loading statistics...")
	file, err := os.Open(statisticsFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Statistics file does not exist. Creating new one.")
			return
		}
		fmt.Println("Error opening statistics file:", err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading statistics file:", err)
		return
	}

	err = json.Unmarshal(data, &StatisticsData)
	if err != nil {
		fmt.Println("Error unmarshaling statistics data:", err)
	} else {
		fmt.Println("Statistics loaded successfully.")
	}
}

func GetStatistics() Statistics {
	return StatisticsData
}

/*func SaveStatistics(stats Statistics) error {
	file, err := os.Create(statisticsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(stats)
	if err != nil {
		return err
	}

	return nil
}
*/

func SaveStatistics() {
	data, err := json.Marshal(StatisticsData)
	if err != nil {
		fmt.Println("Error marshaling statistics data:", err)
		return
	}

	err = os.WriteFile(statisticsFile, data, 0644)
	if err != nil {
		fmt.Println("Error writing statistics file:", err)
	}

}
