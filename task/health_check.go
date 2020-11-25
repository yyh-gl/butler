package task

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	slackWebHook "github.com/ashwanthkumar/slack-go-webhook"
)

type healthCheckNotification struct {
	currentStatusCode int
}

func NewHealthCheckNotification() *healthCheckNotification {
	return &healthCheckNotification{}
}

func (hcn *healthCheckNotification) Action(ctx context.Context) (err error) {
	if hcn.currentStatusCode, err = hcn.checkHealth(ctx); err != nil {
		return err
	}
	return nil
}

func (hcn *healthCheckNotification) checkHealth(ctx context.Context) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://super.hobigon.work/api/v1/health", nil)
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	return resp.StatusCode, nil
}

func (hcn *healthCheckNotification) Notify(ctx context.Context) error {
	if hcn.currentStatusCode == http.StatusOK {
		return nil
	}
	return hcn.notify(ctx)
}

func (hcn *healthCheckNotification) notify(ctx context.Context) error {
	p := slackWebHook.Payload{
		Username: "Health Check Notification",
		Channel:  os.Getenv("SLACK_CHANNEL_NAME_50"),
		Text:     fmt.Sprintf("ヘルスチェックに失敗しました（StatusCode: %d）", hcn.currentStatusCode),
	}

	wh := os.Getenv("SLACK_WEBHOOK_URL_50")
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

func (hcn *healthCheckNotification) Rest(ctx context.Context) error {
	time.Sleep(1 * time.Minute)
	return nil
}
