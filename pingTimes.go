package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

// ping google.com
// 64 bytes from lax17s44-in-x0e.1e100.net (2607:f8b0:4007:80f::200e): icmp_seq=3796 ttl=58 time=18.7 ms
// 64 bytes from lax30s03-in-x0e.1e100.net (2607:f8b0:4007:80f::200e): icmp_seq=3797 ttl=58 time=11.1 ms
// 64 bytes from lax17s44-in-x0e.1e100.net (2607:f8b0:4007:80f::200e): icmp_seq=3798 ttl=58 time=26.6 ms
// 64 bytes from lax17s44-in-x0e.1e100.net (2607:f8b0:4007:80f::200e): icmp_seq=3799 ttl=58 time=16.3 ms
// 64 bytes from lax30s03-in-x0e.1e100.net (2607:f8b0:4007:80f::200e): icmp_seq=3800 ttl=58 time=22.6 ms
// 64 bytes from lax30s03-in-x0e.1e100.net (2607:f8b0:4007:80f::200e): icmp_seq=3801 ttl=58 time=13.6 ms
// 64 bytes from lax30s03-in-x0e.1e100.net (2607:f8b0:4007:80f::200e): icmp_seq=3802 ttl=58 time=11.8 ms
// 64 bytes from lax17s44-in-x0e.1e100.net (2607:f8b0:4007:80f::200e): icmp_seq=3803 ttl=58 time=15.2 ms
// 64 bytes from lax30s03-in-x0e.1e100.net (2607:f8b0:4007:80f::200e): icmp_seq=3804 ttl=58 time=10.9 ms
// 64 bytes from lax17s44-in-x0e.1e100.net (2607:f8b0:4007:80f::200e): icmp_seq=3805 ttl=58 time=18.6 ms
// 64 bytes from lax30s03-in-x0e.1e100.net (2607:f8b0:4007:80f::200e): icmp_seq=3806 ttl=58 time=11.8 ms
// 64 bytes from lax30s03-in-x0e.1e100.net (2607:f8b0:4007:80f::200e): icmp_seq=3807 ttl=58 time=22.8 ms
// 64 bytes from lax17s44-in-x0e.1e100.net (2607:f8b0:4007:80f::200e): icmp_seq=3808 ttl=58 time=35.0 ms

const (
	filenameCst = "sample.txt"
	tsvnameCst  = "sample.csv"
)

var (
	// Passed by "go build -ldflags" for the show version
	commit string
	date   string
)

func main() {

	version := flag.Bool("version", false, "version")
	filename := flag.String("filename", filenameCst, "filename")
	tsvname := flag.String("tsvname", tsvnameCst, "tsvname")

	flag.Parse()

	if *version {
		fmt.Println("commit:", commit, "\tdate(UTC):", date)
		os.Exit(0)
	}

	in, err := os.Open(*filename)
	if err != nil {
		log.Fatal("failed to open:", *filename)
	}
	defer in.Close()

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	re1 := regexp.MustCompile(`icmp_seq=(\d+) ttl=(\d+) time=([\d.]+)`)

	out, e := os.Create(*tsvname)
	if e != nil {
		log.Fatal("failed to open:", *tsvname)
	}
	defer out.Close()

	writer := csv.NewWriter(out)
	writer.Comma = '\t'

	for scanner.Scan() {

		//parts := strings.Split(scanner.Text(), "=")
		//log.Println(parts)
		//2024/03/04 16:30:50 [64 bytes from lax30s03-in-x0e.1e100.net (2607:f8b0:4007:80f::200e): icmp_seq 3857 ttl 58 time 15.7 ms]
		//2024/03/04 16:30:50 [64 bytes from lax17s44-in-x0e.1e100.net (2607:f8b0:4007:80f::200e): icmp_seq 3858 ttl 58 time 25.2 ms]

		result := re1.FindStringSubmatch(scanner.Text())
		log.Print(result)

		record := []string{}
		for k, v := range result {
			log.Printf("%d. %s\n", k, v)
			if k > 0 {
				record = append(record, v)
			}
		}

		//record := []string{string(result[1][1]), string(result[2][1]), string(result[3][1])}

		if err := writer.Write(record); err != nil {
			log.Fatalln("error writing record to tsv:", err)
		}
	}

	writer.Flush()

}
