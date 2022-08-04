<!-- ABOUT THE PROJECT -->
## About The Project

Reconciliation App aims to reconcile and generate reconciliation reports from proxy data and source data.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With

* [![Go][Go]][Go-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/AdiPP/reconciliation.git
   cd ./reconciliation
   ```
2. Build the app
   ```sh
   go build ./cmd/console/main.go
   ```
3. Execute the app
   ```sh
   ./main
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

Reconciliation App uses a command line interface, so it must be run in the terminal (eg: bash).

To perform reconciliation and generate reports, please run the following command:
```sh
./main reconcile {proxyfilepath} {sourcefilepath} {destinationdir} {startdate} {enddate}
```

### Arguments

| Argument         | Type       | Required | Description |
| ---------------- | ---------- | -------- | ----------- |
| `proxyfilepath`  | pathtofile | yes      | Path to proxy file (eg: `./proxy.csv`).     |
| `sourcefilepath` | pathtofile | yes      | Path to source file (eg: `./source.csv`).   |
| `destinationdir` | pathttodir | yes      | Path to destination dir (eg: `./result/`).  |
| `startdate`      | date       | no       | Start date filter. Format in: `DD-MM-YYYY`. |
| `enddate`        | date       | no       | End date filter. Format in: `DD-MM-YYYY`.   |

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Project Structure

```bash
├── cmd
│   ├── console
│   │   ├── main.go
│   ├── data
│   │   ├── main.go
├── pkg
│   ├── command
│   │   ├── commands.go
│   ├── http
│   │   ├── rest
│   │   │   ├──  handlers.go
│   ├── exporting
│   │   ├── export.go
│   │   ├── proxy.go
│   │   ├── report.go
│   │   ├── service.go
│   │   ├── source.go
│   ├── storage
│   │   ├── memory
│   │   │   ├──  proxy.go
│   │   │   ├──  repository.go
│   │   │   ├──  source.go
├── resources
│   ├── proxy.csv
│   ├── source.csv
```

Reconciliation App implement Domain Driven Design (DDD) concept for its structure. The goal itself is to support `Good Structure Goals` presented by Kat Zien at GopherCon 2018, which is:

* Consistent.
* Easy to understand, navigate and reason about.
* Easy to change, loosely-coupled.
* Easy to test.
* "As simple as possible, but no simpler".
* Design reflects exactly how the software works.
* Structure reflects the design exactly.

Check her [Github](https://github.com/katzien) and [Talk](https://youtu.be/oL6JBUk6tj0)


## Test

To perform test, please run the following command:
```sh
go run ./pkg/exporting/
```

<!-- CONTACT -->
## Contact

Adi Putra - [@_adiputrap](https://twitter.com/@_adiputrap) - adiputrapermana@gmail.com

Project Link: [https://github.com/AdiPP/reconciliation](https://github.com/AdiPP/reconciliation)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[Go]: https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/