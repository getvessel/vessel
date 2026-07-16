import { Link, Navigate } from '@tanstack/react-router';
import { AuthPageHeader } from '#/features/auth/auth-page-header';
import { LoginForm } from '#/features/auth/login-form';
import { OAuthButtons } from '#/features/auth/o-auth-buttons';
import { useGetPublicSettings, useGetSetupStatus } from '#/hooks/useSettings';

export function LoginPage() {
  const { data: publicSettings } = useGetPublicSettings();
  const { data: setupStatus, isLoading } = useGetSetupStatus();
  const registrationEnabled = publicSettings?.data?.registrationEnabled ?? true;

  if (!isLoading && setupStatus?.data?.setupRequired) {
    return <Navigate to="/setup" replace />;
  }

  return (
    <div className="fade-in slide-in-from-bottom-4 animate-in duration-300">
      <AuthPageHeader title="Welcome back" description="Sign in to manage this Vessl instance." />

      <OAuthButtons />
      <LoginForm />

      {registrationEnabled && (
        <p className="mt-8 text-center text-sm text-white/50">
          Need an account?{' '}
          <Link to="/signup" className="font-semibold text-primary hover:underline">
            Sign up
          </Link>
        </p>
      )}
    </div>
  );
}
