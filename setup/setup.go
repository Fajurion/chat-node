package setup

import (
	"bufio"
	"chat-node/util"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var StartPort = "3001"

func Setup() bool {

	scanner := bufio.NewScanner(os.Stdin)

	log.Println("Do you want to run this node in testing mode? (y/n)")

	scanner.Scan()
	input := scanner.Text()
	util.Testing = input == "y"

	if util.Testing {
		return true
	}

	if len(strings.Split(input, ":")) > 1 {
		StartPort = strings.Split(input, ":")[1]
		log.Println("Starting on port " + StartPort + "!")
	}

	log.Println("Please provide the file name of the data file (e.g. data (.node will be appended automatically))")

	scanner.Scan()
	input = scanner.Text()

	log.Println("Continuing in normal mode..")

	if readData(util.FilePath + "/" + input + ".node") {
		log.Println("Ready to start.")
		return true
	}

	var creationToken, nodeDomain string
	log.Println("No data file found. Please enter the following information:")

	log.Println("1. Base Path (e.g. http://localhost:3000)")
	scanner.Scan()
	util.BasePath = scanner.Text()

	log.Println("2. Creation Token (Received from a creation request in the admin panel)")
	scanner.Scan()
	creationToken = scanner.Text()

	log.Println("Getting clusters..")
	res, err := util.PostRequest("/node/manage/clusters", fiber.Map{
		"token": creationToken,
	})

	if err != nil {
		log.Println("Your creation token is invalid.")
		return false
	}

	clusterId := setupClusters(res, scanner)

	log.Println("4. App id (e.g. 1)")
	scanner.Scan()
	appId, err := strconv.Atoi(scanner.Text())

	if err != nil {
		log.Println("Please enter a valid number.")
		return false
	}

	log.Println("5. The domain of this node (e.g. node-1.example.com)")
	scanner.Scan()
	nodeDomain = scanner.Text()

	log.Println("6. The performance level (relative to all other nodes) of this node (e.g. 0.75)")
	scanner.Scan()
	performanceLevel, err := strconv.ParseFloat(scanner.Text(), 64)

	if err != nil {
		log.Println("Please enter a valid number.")
		return false
	}

	log.Println("Creating node..")

	res, err = util.PostRequest("/node/manage/new", fiber.Map{
		"token":             creationToken,
		"domain":            nodeDomain,
		"performance_level": performanceLevel,
		"app":               appId,
		"cluster":           clusterId,
	})

	if err != nil {
		log.Println("Error while creating node.")
		return false
	}

	if !res["success"].(bool) {
		log.Println("Error while creating node. Please check your input.")
		return false
	}

	log.Println("Node created successfully.")

	util.NODE_TOKEN = res["token"].(string)
	util.NODE_ID = int(res["id"].(float64))

	// Save data to file
	f, err := os.Create(util.FilePath + "/" + input + ".node")
	if err != nil {
		log.Println("Error while saving data file.")
		return false
	}
	defer f.Close()

	// Write data to file
	f.WriteString(util.BasePath + "\n")
	f.WriteString(util.NODE_TOKEN + "\n")
	f.WriteString(fmt.Sprintf("%d", util.NODE_ID) + "\n")

	log.Println("Data saved to file.")

	return true
}

func readData(path string) bool {

	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	util.BasePath = scanner.Text()

	scanner.Scan()
	util.NODE_TOKEN = scanner.Text()

	scanner.Scan()
	util.NODE_ID, err = strconv.Atoi(scanner.Text())

	if err != nil {
		log.Println("Error while reading data file.")
		return false
	}

	return true
}
