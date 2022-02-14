package mcloud

// Order -
type Order struct {
	ID    int         `json:"id,omitempty"`
	Items []OrderItem `json:"items,omitempty"`
}

// OrderItem -
type OrderItem struct {
	Coffee   Coffee `json:"coffee"`
	Quantity int    `json:"quantity"`
}

// Coffee -
type Coffee struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Teaser      string       `json:"teaser"`
	Description string       `json:"description"`
	Price       float64      `json:"price"`
	Image       string       `json:"image"`
	Ingredient  []Ingredient `json:"ingredients"`
}

// Ingredient -
type Ingredient struct {
	ID       int    `json:"ingredient_id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
}

type SshKey struct {
	Name       string `json:"name"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type SshKeyResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	SshKey  SshKey `json:"SshKey"`
}

type ServerPoolHcloud struct {
	Name          string `json:"name"`
	InstanceType  string `json:"instance_type"`
	InstanceCount int    `json:"instance_count"`
}

type ServerPoolHcloudRequest struct {
	Name          string `json:"name"`
	InstanceType  string `json:"instance_type"`
	InstanceCount int    `json:"instance_count"`
	RunSetup      bool   `json:"run_setup"`
}

type ServerPoolHcloudResponse struct {
	Success          bool             `json:"success"`
	Error            string           `json:"error"`
	ServerPoolHcloud ServerPoolHcloud `json:"ServerPoolHcloud"`
	Task             Task             `json:"Task"`
}

type Task struct {
	Id      int    `json:"id"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Created string `json:"created"`
}

type TaskResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Task    Task   `json:"Task"`
}
