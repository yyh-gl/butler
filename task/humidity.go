package task

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	slackWebHook "github.com/ashwanthkumar/slack-go-webhook"
)

const (
	lowerLimit float64 = 40
	upperLimit float64 = 60
)

type humidityNotification struct {
	currentHumidity float64
}

func NewHumidityNotification() *humidityNotification {
	return &humidityNotification{}
}

func (hn *humidityNotification) Action(ctx context.Context) (err error) {
	hn.currentHumidity, err = hn.fetchHumidity(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (hn *humidityNotification) fetchHumidity(ctx context.Context) (float64, error) {
	type response []struct {
		HumidityOffset float64 `json:"humidity_offset"`
		NewestEvents   struct {
			Humidity struct {
				Val       float64   `json:"val"`
				CreatedAt time.Time `json:"created_at"`
			} `json:"hu"`
		} `json:"newest_events"`
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.nature.global/1/devices", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("REMO_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var r response
	if err := json.Unmarshal(b, &r); err != nil {
		return 0, err
	}

	return r[0].NewestEvents.Humidity.Val + r[0].HumidityOffset, nil
}

func (hn *humidityNotification) Notify(ctx context.Context) error {
	if lowerLimit <= hn.currentHumidity && hn.currentHumidity <= upperLimit {
		return nil
	}
	return hn.notify(ctx)
}

func (hn *humidityNotification) notify(ctx context.Context) error {
	t := "現在湿度：" + strconv.FormatFloat(hn.currentHumidity, 'f', 2, 64) + " %\n"
	if hn.currentHumidity < lowerLimit {
		t += "*加湿しましょう*"
	} else if upperLimit < hn.currentHumidity {
		t += "*除湿しましょう*"
	}

	p := slackWebHook.Payload{
		Username: "Humidity Notification",
		Channel:  os.Getenv("SLACK_CHANNEL_NAME_02"),
		Text:     t,
	}

	wh := os.Getenv("SLACK_WEBHOOK_URL_02")
	errList := slackWebHook.Send(wh, "", p)
	if len(errList) != 0 {
		msg := ""
		for _, e := range errList {
			if msg != "" {
				msg += " & "
			}
			msg += e.Error()
		}
		return errors.New(msg)
	}
	return nil
}

func (hn *humidityNotification) Rest(ctx context.Context) error {
	time.Sleep(1 * time.Hour)
	return nil
}
