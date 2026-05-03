package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"strconv" 

	"github.com/joho/godotenv"
)

type AIModelConnector struct {
	Client *http.Client
}

type Inputs struct {
	Table map[string][]string `json:"table"`
	Query string              `json:"query"`
}

type Response struct {
	Answer      string   `json:"answer"`
	Coordinates [][]int  `json:"coordinates"`
	Cells       []string `json:"cells"`
	Aggregator  string   `json:"aggregator"`
}


func readCSVFile(filename string) (map[string][]string, error) {
	// Baca file CSV
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Baca semua data dari file CSV
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Konversi data CSV ke dalam map
	data, err := CsvToSlice(string(content))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getAIResponse(query string, data map[string][]string, token string) (Response, error) {
	// Buat payload berdasarkan query yang diberikan dan data yang ada
	payload := Inputs{
		Table: data,
		Query: query,
	}

	// Inisialisasi AIModelConnector
	connector := AIModelConnector{Client: &http.Client{}}

	// Panggil AIModelConnector.ConnectAIModel dengan payload dan token
	response, err := connector.ConnectAIModel(payload, token)
	if err != nil {
		return Response{}, err
	}

	return response, nil
}


func CsvToSlice(data string) (map[string][]string, error) {
	reader := csv.NewReader(strings.NewReader(data))
	
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	for i, record := range records {
		if i == 0 {
			for _, column := range record {
				result[column] = []string{}
			}
		} else {
			for j, value := range record {
				columnName := records[0][j]
				result[columnName] = append(result[columnName], value)
			}
		}
	}

	return result, nil
}


func (c *AIModelConnector) ConnectAIModel(payload Inputs, token string) (Response, error) {
	url := "https://router.huggingface.co/hf-inference/models/google/tapas-base-finetuned-wtq"

	// UBAH: bungkus payload dengan key "inputs"
	wrapped := map[string]interface{}{
		"inputs": payload,
	}
	jsonData, err := json.Marshal(wrapped)
	if err != nil {
		return Response{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	// fmt.Println("DEBUG:", string(body))

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Response{}, err
	}

	return response, nil
}


func processResponse(resp Response, data map[string][]string) string {
    result := ""

    switch resp.Aggregator {
    case "SUM":
        total := 0.0
        for _, cell := range resp.Cells {
            val, err := strconv.ParseFloat(cell, 64)
            if err == nil {
                total += val
            }
        }
        result = fmt.Sprintf("Total Penggunaan : %.2f kWh", total)

    case "AVERAGE":
        if len(resp.Cells) == 0 {
            return "Tidak ada data"
        }
        total := 0.0
        for _, cell := range resp.Cells {
            val, err := strconv.ParseFloat(cell, 64)
            if err == nil {
                total += val
            }
        }
        avg := total / float64(len(resp.Cells))
        result = fmt.Sprintf("Rata-rata        : %.2f kWh", avg)

    case "NONE":
        result = fmt.Sprintf("Jawaban          : %s", strings.Join(resp.Cells, ", "))

    default:
        result = fmt.Sprintf("Hasil            : %s", resp.Answer)
    }

    result += fmt.Sprintf("\nCells            : %s", strings.Join(resp.Cells, ", "))
    result += fmt.Sprintf("\nCoordinates      : %v", resp.Coordinates)
    result += fmt.Sprintf("\nAggregator       : %s", resp.Aggregator)

    return result
}


func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Ambil token dari environment variable
	token := os.Getenv("HUGGINGFACE_TOKEN")
	if token == "" {
		fmt.Println("Please set the HUGGINGFACE_TOKEN environment variable in .env file")
		return
	}

	// Baca data dari file CSV
	data, err := readCSVFile("data-series.csv")
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	fmt.Println("Welcome to AI Chatbot CLI!")
	fmt.Println("=================================")
	fmt.Println("1. Berapa total penggunaan listrik Refrigerator?")
	fmt.Println("2. Berapa rata-rata penggunaan listrik Washing Machine?")
	fmt.Println("3. Perangkat apa yang paling boros listrik?")
	fmt.Println("4. Berapa total semua penggunaan listrik?")
	fmt.Println("5. Tanya sendiri (english only)")
	fmt.Println("6. Keluar")
	fmt.Println("=================================")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nPilih menu (1-5): ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "1":
			response, err := getAIResponse("What is the total energy consumption of the Refrigerator?", data, token)
			if err != nil { fmt.Println("Error:", err); continue }
			fmt.Println("Jawaban:", processResponse(response, data))

		case "2":
			response, err := getAIResponse("What is the average energy consumption of the Washing Machine?", data, token)
			if err != nil { fmt.Println("Error:", err); continue }
			fmt.Println("Jawaban:", processResponse(response, data))

		case "3":
			response, err := getAIResponse("Which appliance has the highest energy consumption?", data, token)
			if err != nil { fmt.Println("Error:", err); continue }
			fmt.Println("Jawaban:", processResponse(response, data))

		case "4":
			response, err := getAIResponse("What is the total energy consumption of all appliances?", data, token)
			if err != nil { fmt.Println("Error:", err); continue }
			fmt.Println("Jawaban:", processResponse(response, data))

		case "5":
			fmt.Print("Pertanyaanmu: ")
			scanner.Scan()
			query := strings.TrimSpace(scanner.Text())
			response, err := getAIResponse(query, data, token)
			if err != nil { fmt.Println("Error:", err); continue }
			fmt.Println("Jawaban:", processResponse(response, data))

		case "6":
			fmt.Println("Sampai jumpa!")
			return

			fmt.Println("6. Tanya sendiri (ketik pertanyaan bebas)")


		default:
			fmt.Println("Pilihan tidak ada, coba lagi (1-5).")
		}
	}
}
