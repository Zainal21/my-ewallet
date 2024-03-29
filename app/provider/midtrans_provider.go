package provider

import (
	"context"
	"encoding/base64"
	"log"
	"time"

	"github.com/Zainal21/my-ewallet/pkg/config"
	"github.com/Zainal21/my-ewallet/pkg/httpclient"
	"github.com/Zainal21/my-ewallet/pkg/logger"
	"github.com/sony/gobreaker"
)

type midtrans struct {
	cb  CircuitBreaker
	cfg *config.Config
}

// CreateCharge implements Midtrans.
func (e *midtrans) CreateCharge(ctx context.Context, payload map[string]interface{}) ([]byte, error) {
	url := e.cfg.MidtransBaseUrl + "/v2/charge"
	return e.makeHTTPRequest(ctx, map[string]interface{}{
		"payment_type": "qris",
		"transaction_details": map[string]interface{}{
			"order_id":     payload["order_id"].(string),
			"gross_amount": payload["amount"].(int),
		},
	}, url, "POST")
}

func (e *midtrans) makeHTTPRequest(ctx context.Context, payload interface{}, url, method string) ([]byte, error) {
	reqOption := httpclient.RequestOptions{
		Payload: payload,
		URL:     url,
		Header: map[string]string{
			"Authorization":           "Basic " + base64.StdEncoding.EncodeToString([]byte(e.cfg.MidtransServerKey+":")),
			"Accept":                  "application/json",
			"Content-Type":            "application/json",
			"X-Override-Notification": e.cfg.MidtransCallbackUrl,
		},
		Method:        method,
		TimeoutSecond: 24,
		Context:       ctx,
	}

	result, err := e.cb.Execute(ctx, SendRequest, reqOption)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return result.(httpclient.Response).RawByte(), nil
}

func NewMidtransProvider(cnf *config.Config) Midtrans {
	var (
		goBreakerName       = "example-provider"
		goBreakerMaxRequest = 100
		goBreakerInterval   = 20 * time.Second
		goBreakerTimeout    = 5 * time.Second
	)

	settings := gobreaker.Settings{
		Name: goBreakerName,
		// MAXREQUESTS IS THE MAXIMUM NUMBER OF REQUESTS ALLOWED TO PASS THROUGH WHEN THE CIRCUITBREAKER IS HALF-OPEN
		MaxRequests: uint32(goBreakerMaxRequest),

		// INTERVAL IS THE CYCLIC PERIOD OF THE CLOSED STATE FOR CIRCUITBREAKER TO CLEAR THE INTERNAL COUNTS
		Interval: goBreakerInterval,

		// TIMEOUT IS THE PERIOD OF THE OPEN STATE, AFTER WHICH THE STATE OF CIRCUITBREAKER BECOMES HALF-OPEN.
		Timeout: goBreakerTimeout,
		OnStateChange: func(name string, from, to gobreaker.State) {
			var (
				lf = logger.NewFields(
					logger.EventName("ExampleProviderGoBreaker"),
				)
			)
			lf.Append(logger.Any("gobreaker.name", goBreakerName))
			lf.Append(logger.Any("gobreaker.maxRequest", goBreakerMaxRequest))
			lf.Append(logger.Any("gobreaker.timeout", goBreakerInterval))
			lf.Append(logger.Any("gobreaker.timeout", goBreakerTimeout))

			// TO DO CALLBACK GOBREAKER WITH CURRENT STATE
			switch to {
			case gobreaker.StateOpen:
				lf.Append(logger.String("gobreaker.state", "OPEN"))
				logger.ErrorWithContext(context.Background(), "success trigger gobreaker with OPEN state", lf...)
			case gobreaker.StateHalfOpen:
				lf.Append(logger.String("gobreaker.state", "HALF-OPEN"))
				logger.WarnWithContext(context.Background(), "success trigger gobreaker with HALF-OPEN state", lf...)
			case gobreaker.StateClosed:
				lf.Append(logger.String("gobreaker.state", "CLOSED"))
				logger.InfoWithContext(context.Background(), "success trigger gobreaker with CLOSED state", lf...)
			}

		},
		// READYTOTRIP IS CALCULATION MOVE TO OPEN STATE
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests > 2 && failureRatio >= 0.3
		},
	}
	cb := NewCircuitBreaker(cnf, settings)
	return &midtrans{
		cb:  cb,
		cfg: cnf,
	}
}
