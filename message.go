package main

type Envelope struct {
	Type  string `json:"type"`
	Token string `json:"token"`
	UID   string `json:"uid"`
}

type LoginRequest struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

type LoginResponse struct {
	Status  string `json:"status"`
	Token   string `json:"token"`
	Message string `json:"msg"`
}

type DragStartRequest struct {
	CID string
}

type DragStartResponse struct {
	Status string `json:"status"`
	Event  string `json:"event"`
	CID    string `json:"cid"`
	UID    string `json:"uid"`
}

type DragStartFailResponse struct {
	Status  string `json:"status"`
	Message string `json:"msg"`
}

type DragFinishRequest struct {
	CID string     `json:"cid"`
	Pos [2]float64 `json:"position"`
}

type DragFinishResponse struct {
	Status string     `json:"status"`
	Event  string     `json:"event"`
	CID    string     `json:"cid"`
	UID    string     `json:"uid"`
	Pos    [2]float64 `json:"position"`
}

type DragFinishFailResponse struct {
	Status  string `json:"status"`
	Message string `json:"msg"`
}
