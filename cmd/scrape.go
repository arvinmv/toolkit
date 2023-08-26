/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
    "github.com/gocolly/colly/v2"
	"github.com/spf13/cobra"
)

var filename string
// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape redfin addresses for the Estimated price and Appreciation",
	Long: `Scrape addresses from a file passed. Usage: toolkit scrape --file <filename.txt>`,
	Run: func(cmd *cobra.Command, args []string) {
		if filename == "" {
			log.Fatalf("Please provide a file name containing addresses. Usage: toolkit scrape --file <filename.txt>")
		}

		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		var addresses []string
		for scanner.Scan() {
			addresses = append(addresses, scanner.Text())
		}

		for _, address := range addresses {
			fmt.Println("Scraping address:", address)
			scrapeAddress(address)
			fmt.Println("------")
		}
	},
}

func scrapeAddress(url string) {
	c := colly.NewCollector()
	counter := 0

	c.OnHTML("div[class*='price']", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		if strings.Contains(text, "$") {
			counter++
			if counter == 1 {
				fmt.Println("Estimated Price:", text)
			} else if counter == 2 {
				// Splitting text on "•" to separate price details and timestamp
				parts := strings.Split(text, "•")
				if len(parts) > 1 {
					fmt.Println("Details:", strings.TrimSpace(parts[0]))
				} else {
					fmt.Println("Details:", text)
				}
			}
		}
	})

	// Set up error handling
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start the web scraping for the given address
	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
}


func init() {
	rootCmd.AddCommand(scrapeCmd)
	scrapeCmd.Flags().StringVarP(&filename, "file", "f", "", "Filename containing addresses")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scrapeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scrapeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
