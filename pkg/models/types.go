package models

type Deployment struct {
	Name             string `json:"name"`
	ApplicationGroup string `json:"applicationGroup"`
	RunningPodsCount int32  `json:"runningPodsCount"`
}

type Deployments []*Deployment

type Pod struct {
	Name   string            `json:"name"`
	Status string            `json:"status"`
	Labels map[string]string `json:"labels"`
}

type Pods []*Pod

type Client struct {
	Name   string
	Events chan *SseMessage
}
type SseMessage struct {
	Message string `json:"message"`
}
