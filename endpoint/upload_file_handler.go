package endpoint

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Nhanderu/brdoc"
	"github.com/albertojnk/neoway-db-manipulation/datasource"
	"github.com/cuducos/go-cnpf"
	"github.com/labstack/echo"
)

func uploadFileHandler(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("error while reading file, err: ", err)
		return c.JSON(http.StatusBadRequest, "please check if file was sent")
	}

	src, err := file.Open()
	if err != nil {
		log.Println("error while openning file, err: ", err)
		return c.JSON(http.StatusInternalServerError, "please check if file is not corrupted")
	}

	defer src.Close()

	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()
	defer os.Remove(file.Filename)

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	lines := readFileLines(file.Filename)

	indexes := getStartIndexes(lines[1])

	rawData := lines[1:]

	data := parseData(rawData, indexes)

	err = datasource.BulkCreateClientInfo(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to save data into db")
	}

	return c.JSON(http.StatusCreated, data)
}

func parseData(rawData []string, indexes []int) []datasource.ClientInfo {

	client := make([]datasource.ClientInfo, 0)
	for _, rd := range rawData {
		line := make([]string, 0)

		// get data from rawData based on the startIndex of each value and append the values to line
		for i, v := range indexes {
			if i < len(indexes)-1 {
				newLine := rd[v:indexes[i+1]]
				newLine = strings.ReplaceAll(newLine, " ", "")
				line = append(line, newLine)
				continue
			}
			newLine := rd[v:]
			newLine = strings.ReplaceAll(newLine, " ", "")
			line = append(line, newLine)
		}

		cpf := cnpf.Unmask(line[0])
		private, _ := strconv.ParseBool(line[1])
		incomplete, _ := strconv.ParseBool(line[2])
		lastPurchaseDate, _ := time.Parse("2006-01-02", line[3])
		avgBudget, _ := strconv.ParseFloat(strings.ReplaceAll(line[4], ",", "."), 64)
		lastBudget, _ := strconv.ParseFloat(strings.ReplaceAll(line[5], ",", "."), 64)
		freqStore := cnpf.Unmask(line[6])
		lastStore := cnpf.Unmask(line[7])

		client = append(client, datasource.ClientInfo{
			CPF:                  strings.ToUpper(cpf),
			IsValidCPF:           brdoc.IsCPF(line[0]),
			Private:              private,
			Incomplete:           incomplete,
			LastPurchaseDate:     lastPurchaseDate,
			AverageBudget:        avgBudget,
			LastPurchaseBudget:   lastBudget,
			MostFrequentStore:    strings.ToUpper(freqStore),
			IsValidFrequentStore: brdoc.IsCNPJ(line[6]),
			LastPurchaseStore:    strings.ToUpper(lastStore),
			IsValidLastStore:     brdoc.IsCNPJ(line[7]),
		})
	}

	return client
}

func readFileLines(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) < 2 {
		log.Println("Database arquive must have atleast 2 lines")
		return nil
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}

func getStartIndexes(s string) []int {
	indexes := make([]int, 0)

	for i := range s {
		if i == 0 {
			indexes = append(indexes, i)
			continue
		}
		if s[i] != 32 && s[i-1] == 32 {
			indexes = append(indexes, i)
		}
	}

	return indexes
}
