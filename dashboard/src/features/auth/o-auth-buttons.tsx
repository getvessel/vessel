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
      <div className="mb-5 flex flex-col space-y-2">
        {enabledProviders.map((provider) => (
          <Button
            key={provider}
            variant="outline"
            type="button"
            onClick={() => handleOAuthLogin(provider)}
            className="h-10 w-full border-white/10 bg-white/[0.055] font-medium text-sm text-white capitalize hover:bg-white/[0.08]"
          >
            Continue with {provider}
          </Button>
        ))}
      </div>
      <div className="relative mb-5">
        <div className="absolute inset-0 flex items-center">
          <span className="w-full border-white/10 border-t" />
        </div>
        <div className="relative flex justify-center text-xs uppercase">
          <span className="bg-[#111114] px-2 font-semibold text-white/40">
            Or continue with email
          </span>
        </div>
      </div>
    </>
  );
};
