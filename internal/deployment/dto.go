package deployment

type TriggerDeploymentRequest struct {
	Branch *string `json:"branch,omitempty"`
}
