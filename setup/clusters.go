package setup

import (
	"bufio"
	"log"
	"strconv"
)

func setupClusters(res map[string]interface{}, scanner *bufio.Scanner) int64 {

	var clusterId int64
	for clusterId == 0 {

		log.Println("3. Please select a cluster:")

		clusters := res["clusters"].([]interface{})
		for i, cluster := range clusters {

			cluster := cluster.(map[string]interface{})
			log.Println(i, cluster["name"], " (", cluster["id"], ") - ", cluster["country"])
		}

		scanner.Scan()
		input, err := strconv.Atoi(scanner.Text())

		if err != nil {
			log.Println("Please enter a valid number.")
			continue
		}

		if input < 0 || input >= len(clusters) {
			log.Println("Please enter a valid number.")
			continue
		}

		clusterId = int64(clusters[input].(map[string]interface{})["id"].(float64))
	}

	return clusterId
}
