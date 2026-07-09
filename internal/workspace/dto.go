package workspace

type CreateWorkspaceRequest struct {
	Name            string `json:"name"`
	AvatarURL       string `json:"avatarUrl,omitempty"`
	PreferredRegion string `json:"preferredRegion,omitempty"`
}

type UpdateWorkspaceRequest struct {
	Name            string `json:"name,omitempty"`
	AvatarURL       string `json:"avatarUrl,omitempty"`
	PreferredRegion string `json:"preferredRegion,omitempty"`
}

type CreateTrustedDomainRequest struct {
	Domain string `json:"domain"`
	Role   string `json:"role"`
}

type CreateSSHKeyRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}
