package sampstar

import "testing"

func TestStart(t *testing.T) {
	server := NewSampstarServer(&SampstarOptions{
		Port: 7777,
	})

	server.Start()
}
