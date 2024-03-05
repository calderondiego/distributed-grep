package client

import (
	"sync"
	"testing"
)

// Define a test struct to store test cases. Each test case includes a name, the number of concurrent clients,
// a query to be executed on servers, and the expected output for the query.
var tests = []struct {
	name       string
	concurrent int
	query      string
	expected   string
}{
	{"count lines", 1, "grep -c '' testdata/apache.log", "1000"},
	{"count ips", 1, "grep -c -E '((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)' testdata/apache.log", "1000"},
	{"count http 1.x", 1, "grep -c -E 'HTTP\\/1\\.[0-9]' testdata/apache.log", "657"},
	{"count specific user", 1, "grep -c 'lindgren4315' testdata/apache.log", "1"},
	{"count POST requests", 1, "grep -c 'POST /' testdata/apache.log", "171"},

	{"count lines (multi)", 5, "grep -c '' testdata/apache.log", "1000"},
	{"count ips (multi)", 5, "grep -c -E '((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)' testdata/apache.log", "1000"},
	{"count http 1.x (multi)", 5, "grep -c -E 'HTTP\\/1\\.[0-9]' testdata/apache.log", "657"},
	{"count specific user (multi)", 5, "grep -c 'lindgren4315' testdata/apache.log", "1"},
	{"count POST requests (multi)", 5, "grep -c 'POST /' testdata/apache.log", "171"},
}

func TestQueryServers(t *testing.T) {
	servers := GetServersFromFile("testdata/servers.txt")
	if servers == nil {
		t.Errorf("Error reading servers from file")
		return
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(test.concurrent)

			// Spawn multiple goroutines to simulate concurrent clients.
			for i := 0; i < test.concurrent; i++ {
				go func(clientId int) {
					outputs := QueryServers(servers, test.query)
					for _, output := range outputs {
						if output.Output != test.expected {
							t.Errorf("Client %d output from '%s': '%s' not equal to expected '%s'",
								clientId, output.Address, output.Output, test.expected)
						}
					}
					wg.Done()
				}(i)
			}
			wg.Wait()
		})
	}
}

var benchmarks = []struct {
	name  string
	query string
}{
	{"count single ip", "grep -c '130\\.192\\.18\\.[0-9]{1,3}' vm*.log"},
	{"count lines", "grep -c '' vm*.log"},
	{"count POST requests", "grep -c 'POST /' vm*.log"},
	{"count ips", "grep -c -E '((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)' vm*.log"},
}

func BenchmarkStartClient(b *testing.B) {
	for _, benchmark := range benchmarks {
		b.Run(benchmark.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				StartClient("testdata/servers.txt", benchmark.query)
			}
		})
	}
}
