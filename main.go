package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func HelloWorld() string {
	return `Hello World`
}

var run = true

type Feed struct {
	CreatedAt string `json:"created_at"`
	Field1    string `json:"field1"`
}

func httpGraphics() []float64 {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.thingspeak.com/channels/9/feeds.json", nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return nil
	}
	defer resp.Body.Close()

	var data []float64

	var feedData struct {
		Feeds []Feed `json:"feeds"`
	}
	err = json.NewDecoder(resp.Body).Decode(&feedData)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return nil
	}

	for _, feed := range feedData.Feeds {
		value, err := strconv.ParseFloat(feed.Field1, 64)
		if err != nil {
			fmt.Println("Ошибка при преобразовании строки в число:", err)
			continue
		}
		data = append(data, value)
	}

	return data
}

func showGraphics(data []float64) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	bc := widgets.NewBarChart()
	bc.Title = "Bar Chart"
	bc.SetRect(5, 5, 70, 36)
	bc.Labels = []string{"Value"}
	bc.BarColors[0] = ui.ColorBlue

	pause := func() {
		run = !run
		if run {
			bc.Title = "Bar Chart"
		} else {
			bc.Title = "Bar Chart (Stopped)"
		}
		ui.Render(bc)
	}

	ui.Render(bc)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "s":
				pause()
			}
		case <-ticker:
			if run {
				bc.Data = data
				ui.Render(bc)
			}
		}
	}
}

func main() {
	data := httpGraphics()
	if data == nil {
		return
	}

	showGraphics(data)
}
