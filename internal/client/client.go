package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

type GrepOutput struct {
	Machine int
	Address string
	Output  string
}

// QueryServers performs a grep query on a list of remote servers concurrently.
// It returns a sorted list of GrepOutput containing results from all servers.
func QueryServers(servers []string, query string) []GrepOutput {
	var wg sync.WaitGroup
	wg.Add(len(servers))
	outputChannel := make(chan GrepOutput, len(servers))

	for i := 0; i < len(servers); i++ {
		go func(machineId int, ip string) {
			connection, err := net.DialTimeout("tcp", ip, 5*time.Second)

			if err != nil {
				outputChannel <- GrepOutput{Machine: machineId, Address: ip, Output: ""}
				return
			}

			defer connection.Close()

			connection.SetDeadline(time.Now().Add(5 * time.Second))
			_, err = connection.Write([]byte(query))
			if err != nil {
				outputChannel <- GrepOutput{Machine: machineId, Address: ip, Output: ""}
				return
			}

			buffer := make([]byte, 2048)
			mLen, err := connection.Read(buffer)
			if err != nil {
				outputChannel <- GrepOutput{Machine: machineId, Address: ip, Output: ""}
			} else {
				outputChannel <- GrepOutput{Machine: machineId, Address: ip, Output: string(buffer[:mLen])}
			}
		}(i, servers[i])
	}

	sortedOutput := make([]GrepOutput, 0, len(servers))
	go func() {
		for grepOutput := range outputChannel {
			sortedOutput = append(sortedOutput, grepOutput)
			wg.Done()
		}
	}()

	wg.Wait()

	sort.SliceStable(sortedOutput, func(i, j int) bool {
		return sortedOutput[i].Machine < sortedOutput[j].Machine
	})
	return sortedOutput
}

func GetServersFromFile(filepath string) []string {
	servers := []string{}
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var address = scanner.Text()
		servers = append(servers, address)
	}
	return servers
}

func StartClient(serversFile string, query string) {
	servers := GetServersFromFile(serversFile)
	if servers == nil {
		return
	}

	sortedOutput := QueryServers(servers, query)
	for _, output := range sortedOutput {
		if output.Output != "" { // don't print failed machines
			fmt.Printf("%s | %s\n", output.Address, output.Output)
		}
	}
}
