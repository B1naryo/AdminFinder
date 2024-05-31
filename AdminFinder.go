package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Estrutura para armazenar o resultado de uma verificação de diretório
type DirectoryCheckResult struct {
	Directory string
	Err       error
}

// Função que verifica um único diretório
func checkDirectory(baseURL, directory string, wg *sync.WaitGroup, results chan<- DirectoryCheckResult) {
	defer wg.Done()

	testURL := fmt.Sprintf("%s/%s", baseURL, directory)
	fmt.Println("Verificando diretório:", testURL) // Log de depuração

	for attempt := 1; attempt <= 3; attempt++ { // Tentar 3 vezes
		resp, err := http.Get(testURL)
		if err != nil {
			fmt.Printf("Erro na tentativa %d ao verificar %s: %v\n", attempt, testURL, err)
			time.Sleep(1 * time.Second) // Esperar um pouco antes de tentar novamente
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			results <- DirectoryCheckResult{Directory: testURL, Err: nil}
			return
		}

		// Se não for um status OK, retornar um erro
		results <- DirectoryCheckResult{Directory: testURL, Err: fmt.Errorf("Status de resposta não OK: %d", resp.StatusCode)}
		return
	}

	// Se todas as tentativas falharem, retornar um erro
	results <- DirectoryCheckResult{Directory: testURL, Err: fmt.Errorf("Todas as tentativas falharam")}
}

// Função para verificar todos os diretórios para uma única URL base
func crawl(baseURL string, directories []string, results chan<- DirectoryCheckResult) {
	var wg sync.WaitGroup
	for _, directory := range directories {
		wg.Add(1)
		go checkDirectory(baseURL, directory, &wg, results)
	}
	wg.Wait()
}

// Função para salvar o diretório encontrado em um arquivo
func saveDirectory(directory string) {
	file, err := os.OpenFile("diretorios_encontrados.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Erro ao abrir ou criar o arquivo de saída:", err)
		return
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s\n", directory)
	if err != nil {
		fmt.Println("Erro ao escrever no arquivo de saída:", err)
	}
}

func main() {
	// Argumentos de linha de comando
	urlFlag := flag.String("u", "", "URL base")
	directoriesFileFlag := flag.String("d", "admin.txt", "Arquivo contendo diretórios")
	flag.Parse()

	// Verificar se pelo menos a flag -u foi fornecida
	if *urlFlag == "" {
		fmt.Println("Por favor, forneça uma URL base com -u.")
		return
	}

	fmt.Println("URL base:", *urlFlag) // Log de depuração

	// Ler os diretórios do arquivo
	file, err := os.Open(*directoriesFileFlag)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo de diretórios:", err)
		return
	}
	defer file.Close()

	fmt.Println("Lendo diretórios do arquivo:", *directoriesFileFlag) // Log de depuração

	scanner := bufio.NewScanner(file)
	directories := []string{}
	for scanner.Scan() {
		directory := strings.TrimSpace(scanner.Text())
		if directory != "" {
			directories = append(directories, directory)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro ao ler o arquivo de diretórios:", err)
		return
	}

	fmt.Println("Diretórios a serem verificados:", directories) // Log de depuração

	// Canal para coletar os resultados
	results := make(chan DirectoryCheckResult, len(directories))

	// Verificar diretórios para uma única URL base
	go crawl(*urlFlag, directories, results)

	// Exibir os resultados e salvar diretórios encontrados
	for range directories {
		result := <-results
		if result.Err != nil {
			fmt.Printf("Erro ao verificar %s: %v\n", result.Directory, result.Err)
		} else {
			fmt.Println("Diretório encontrado:", result.Directory)
			saveDirectory(result.Directory)
		}
	}
}

