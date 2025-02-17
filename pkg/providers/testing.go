package providers

import (
	"context"
	"time"

	"glide/pkg/routers/latency"

	"glide/pkg/api/schemas"
)

type ResponseMock struct {
	Msg string
	Err *error
}

func (m *ResponseMock) Resp() *schemas.ChatResponse {
	return &schemas.ChatResponse{
		ID: "rsp0001",
		ModelResponse: schemas.ModelResponse{
			SystemID: map[string]string{
				"ID": "0001",
			},
			Message: schemas.ChatMessage{
				Content: m.Msg,
			},
		},
	}
}

type ProviderMock struct {
	idx              int
	responses        []ResponseMock
	supportStreaming bool
}

func NewProviderMock(responses []ResponseMock, supportStreaming bool) *ProviderMock {
	return &ProviderMock{
		idx:              0,
		responses:        responses,
		supportStreaming: supportStreaming,
	}
}

func (c *ProviderMock) Chat(_ context.Context, _ *schemas.ChatRequest) (*schemas.ChatResponse, error) {
	response := c.responses[c.idx]
	c.idx++

	if response.Err != nil {
		return nil, *response.Err
	}

	return response.Resp(), nil
}

func (c *ProviderMock) SupportChatStream() bool {
	return c.supportStreaming
}

func (c *ProviderMock) ChatStream(_ context.Context, _ *schemas.ChatRequest, _ chan<- schemas.ChatResponse) error {
	// TODO: implement
	return nil
}

func (c *ProviderMock) Provider() string {
	return "provider_mock"
}

type LangModelMock struct {
	modelID string
	healthy bool
	latency *latency.MovingAverage
	weight  int
}

func NewLangModelMock(ID string, healthy bool, avgLatency float64, weight int) *LangModelMock {
	movingAverage := latency.NewMovingAverage(0.06, 3)

	if avgLatency > 0.0 {
		movingAverage.Set(avgLatency)
	}

	return &LangModelMock{
		modelID: ID,
		healthy: healthy,
		latency: movingAverage,
		weight:  weight,
	}
}

func (m *LangModelMock) ID() string {
	return m.modelID
}

func (m *LangModelMock) Healthy() bool {
	return m.healthy
}

func (m *LangModelMock) Latency() *latency.MovingAverage {
	return m.latency
}

func (m *LangModelMock) LatencyUpdateInterval() *time.Duration {
	updateInterval := 30 * time.Second

	return &updateInterval
}

func (m *LangModelMock) Weight() int {
	return m.weight
}
