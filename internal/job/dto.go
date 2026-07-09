package job

type CreateJobRequest struct {
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`
	Schedule  string `json:"schedule"`
	Command   string `json:"command"`
}

type UpdateJobRequest struct {
	Name     *string `json:"name,omitempty"`
	Schedule *string `json:"schedule,omitempty"`
	Command  *string `json:"command,omitempty"`
	Status   *string `json:"status,omitempty"`
}
