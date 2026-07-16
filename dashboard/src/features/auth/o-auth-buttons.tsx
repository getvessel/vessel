import { Button } from '#/components/ui/button';
import { useEnabledOAuthProviders } from '#/hooks/useOAuth';
import { oauthService } from '#/services/oauth';

export const OAuthButtons = () => {
  const { data: enabledProvidersData } = useEnabledOAuthProviders();
  const enabledProviders = (enabledProvidersData?.data || []).map((p) =>
    p.providerName.toLowerCase()
  );

  const handleOAuthLogin = (provider: string) => {
    oauthService.triggerOAuthLogin(provider);
  };

  if (!enabledProviders.length) return null;

  return (
    <>
      <div className="mb-6 flex flex-col space-y-3">
        {enabledProviders.map((provider) => (
          <Button
            key={provider}
            variant="outline"
            type="button"
            onClick={() => handleOAuthLogin(provider)}
            className="h-11 w-full border-border bg-background font-medium text-base capitalize hover:bg-muted/50"
          >
            Continue with {provider}
          </Button>
        ))}
      </div>
      <div className="relative mb-6">
        <div className="absolute inset-0 flex items-center">
          <span className="w-full border-border border-t" />
        </div>
        <div className="relative flex justify-center text-xs uppercase">
          <span className="bg-background px-2 font-semibold text-muted-foreground">
            Or continue with email
          </span>
        </div>
      </div>
    </>
  );
};
