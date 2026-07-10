package billing

type StripeService struct{}

func NewStripeService() *StripeService {
	return &StripeService{}
}

func (s *StripeService) CreateCheckoutSession(teamID string) (string, error) {
	// Implementation for Stripe
	return "chk_stub", nil
}
