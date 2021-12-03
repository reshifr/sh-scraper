all: cmd/main.go pkg/scraper.go pkg/validator.go
	go build -o scraping cmd/main.go
	@chmod +x scraping &> /dev/null

clean:
	@rm -rf scraping &> /dev/null
