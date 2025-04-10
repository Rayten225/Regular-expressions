package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// calculateExpression вычисляет математическое выражение и возвращает результат в виде строки.
func calculateExpression(s string, re *regexp.Regexp) (string, error) {
	matches := re.FindStringSubmatch(s)
	if len(matches) != 4 {
		return "", fmt.Errorf("неверный формат выражения")
	}

	num1, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", err
	}

	num2, err := strconv.Atoi(matches[3])
	if err != nil {
		return "", err
	}

	var result int
	switch matches[2] {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	default:
		return "", fmt.Errorf("неподдерживаемый оператор: %s", matches[2])
	}

	expression := s[:len(s)-2]
	return expression + "=" + strconv.Itoa(result), nil
}

func main() {
	inputFile := os.Args[1]
	var outputFile string
	if len(os.Args) == 2 {
		outputFile = "output.txt"
	} else {
		outputFile = os.Args[2]
	}

	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Ошибка при чтении входного файла: %v", err)
	}

	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Ошибка при создании директории: %v", err)
	}

	out, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Ошибка при открытии выходного файла: %v", err)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	defer writer.Flush()

	re := regexp.MustCompile(`(\d+)([\+\-])(\d+)=\?`)

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if re.MatchString(line) {
			result, err := calculateExpression(line, re)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка в выражении '%s': %v\n", line, err)
				continue
			}
			_, err = writer.WriteString(result + "\n")
			if err != nil {
				log.Fatalf("Ошибка при записи в файл: %v", err)
			}
		}
	}
}
