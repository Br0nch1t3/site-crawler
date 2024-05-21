# Site Crawler

A program written in go intended to automatically create a sitemap for your website by crawling every text link it can find.
This program makes use of goroutines to accelerate the crawling process.

## Build

```sh
cd crawler
go build .
```

## Usage

```sh
./crawler [OPTIONS] [URL]
```

| option | description | default |
|----------|:-------------:|:-------------:|
| -v, --verbose  |  Verbose output | false |
| -d, --depth value  |  Max depth to crawl to | none |
| -h, --help | Prints this help message and exit | none |
