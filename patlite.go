package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func send_patlite(ipAddress string, alertStatus string) error {
	url := fmt.Sprintf("http://%s/api/control?alert=%s", ipAddress, alertStatus)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(byteArray) != "Success." {
		return fmt.Errorf("Failed to send alert to patlite.")
	}

	return nil
}

func send_patlites(alertStatus string) {
	go func() {
		for _, patlite := range patlites {
			send_patlite(patlite, alertStatus)
		}
	}()
}
