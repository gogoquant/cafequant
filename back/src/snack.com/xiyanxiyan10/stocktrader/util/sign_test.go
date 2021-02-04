package util

import (
	"net/url"
	"testing"
)

const (
	PrivateKey = "MIIEogIBAAKCAQEAxdg8LfjAIINluvMytwIdsewrumHEp7q3n4FzN7UpGakPJPAxgXrgjFY88G1tG+RchiKy5TXLA7k4fAH35glxeXZgdkmDhCZ0H2hM75eeC3cnX3Pdu0UNxDvt2WlIcmIsAHJDmZjROyt1Q4TN6JhzAu4POXRO8a3+t9pdzYioU8I9gesiE9YmyhOguuhjMzuS1q1necMEIqQIMUdNoDSBpRyjMQZ0HtdzkjtGtf9uP89wQ1a9yaKcylJ35kA9IZ/Hghs1bhkIzLeiOW7eTbrEC4cy9xGAxibQ8MjKSm6JFlm+C+bzVYV2ipCteOi5V3sMROHsL37e+fMP9kbaS7Hk7wIDAQABAoIBAACFspr4diFf12vn6nFbOxLWKcNjMK60qnlsUQ6LluEvdg/F5ouN9HvKWnzT/R6+upPMEabTPoby/TgulSXxTnBgpJ6LUSKPK21NzC4xu0QSe3MgDizJYODsu5MAWSWcJruVkaIdKig61CNqfVSo2lzengGr0e2HZQ29MNQzESavcmJ8cEm2ui/ddr5KdOiN3nBlfWBZfo038964YdHUQm76PHN7QcQNSXSLnoUN3SCFRXuio03FaVcuOo/1+3R5aqgN4JQ9s3/Sy7UgbvuZsf6+zVonOMknt8b+we4IihAZMR/K3MsnLMz8XmcqdXeDjtDPbxy30rLk6epk+A615NkCgYEA+DhP9xAVfqOttAYpucHPV8EVsIvrx/8KmhELJumncRtSwg9IqjPcylkAW1D7MgeLUT+uLn/paZ5XGScj8RsoIq+gSUQ5an+OmxkRJRITQ9XyVNO15g5ouUFaLVFkcJ2Co4unqpZByCHCLxkTmAyjHWtxr7dX7FkFmn7GqsErqG0CgYEAzAu3njYKysmwfhYdEn2XJ3RnN+gc7iLaaLDhFh9/W6ffcRZV1sdi+WES3o9cwNSIL24F4BHp3fLog3XLxiumzY8ZzOnUKhppWhu309GWV5qBfIoXoBQ5A2Sm9JauNR+tVzHW1U1eYrQGJkbhgv9WtJf1h7Zl9gGvlwul/bQHoUsCgYBEogVyUe8vkgBwm5ej9jPnlsrxgu7R4PJEgVvtCYQz4RMz91fnP+nXxV404aJjRfS+pXX7A4E9o/t/R/RHMXQaiyctuwCJMvXyaq7z6hiLlDeqPtO35doNB0Xw6+VywgqiP/Y/U8aimLsBnNRvIWdkthW8OVzFTCQhgNZb1ofEzQKBgB3t9hAJ40ldjjrgaYFF1L8fzugfbubrS9ghYdLJ6fd6x0aiPRMVCgqEV603oCZUxmkWnVwBpKk+sSZfR/WYf44VWHZ7Mfi/CQcDm9JBIulUq3umEdMUREygHfEwjPsT22w3zkhZYefeeixxJzD83S3+QDCY65nLI4NnXQC6wIfDAoGASnAgfyD4unueYzVFe6/wa9NZmgmygcFn04ujRmqLDduD6OvejrYChoHnL6UC32xQ/g+Ti9nipV+J0S2X1JEnqU4bXi01uXl9Uzog2j1Fs6gZSwt+kdbbJseBVJ2t104SdKLInCLObJ4m3Fxew0YMJ9/L+pbkjqkI1myVot3T1hk="
	PublicKey  = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxdg8LfjAIINluvMytwIdsewrumHEp7q3n4FzN7UpGakPJPAxgXrgjFY88G1tG+RchiKy5TXLA7k4fAH35glxeXZgdkmDhCZ0H2hM75eeC3cnX3Pdu0UNxDvt2WlIcmIsAHJDmZjROyt1Q4TN6JhzAu4POXRO8a3+t9pdzYioU8I9gesiE9YmyhOguuhjMzuS1q1necMEIqQIMUdNoDSBpRyjMQZ0HtdzkjtGtf9uP89wQ1a9yaKcylJ35kA9IZ/Hghs1bhkIzLeiOW7eTbrEC4cy9xGAxibQ8MjKSm6JFlm+C+bzVYV2ipCteOi5V3sMROHsL37e+fMP9kbaS7Hk7wIDAQAB"
)

func TestSign(t *testing.T) {
	var p = url.Values{}
	start := "2019-07-05 17:18:32"
	random := "123456"
	p.Add("username", "wei.bai")
	p.Add("time-stamp", start)
	p.Add("random", random)
	//p.Add("token", "aaaaaaaa")
	sign, err := SinData(p, PrivateKey)
	if err != nil {
		t.Log(err)
	}
	t.Logf("sign=%s", sign)
	var data = url.Values{}
	data.Add("username", "wei.bai")
	data.Add("time-stamp", start)
	data.Add("random", random)
	//data.Add("token", "aaaaaaaa")
	data.Add("sign", sign)
	data.Add("sign_type", KSIGNTYPERSA2)
	ok, err := VerifySign(data, PublicKey)
	if err != nil {
		t.Log(err)
	}
	t.Logf("OK:%t", ok)

}
