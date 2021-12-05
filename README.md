# Sh-Scraper
Scrape data from Unix & Linux Stack Exchange using The Onion Router as a web proxy.

## Dependencies
- The Onion Router (Tor)
- SQLite3

## Start
```shell
$ git clone https://github.com/reshifr/sh-scraper.git
$ cd sh-scraper
$ make
$ chmod +x initial
$ ./initial
```

## Run
```shell
$ sudo ./scraping
```

## Finish
The scraped data is stored in the `data/shl.db`. To resume the previous scraping result, simply rerun again.
