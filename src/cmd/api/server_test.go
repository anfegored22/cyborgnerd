package api

import "testing"

func TestStartReportsListenFailure(t *testing.T) {
	s := portfolioTestServer(t)
	if err := s.Start("127.0.0.1:not-a-port"); err == nil {
		t.Fatal("Start must report listen failures so systemd can restart the process")
	}
}
