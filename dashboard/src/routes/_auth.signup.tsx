import { createFileRoute, Link, Navigate } from '@tanstack/react-router';
import { AuthPageHeader } from '#/features/auth/auth-page-header';
import { OAuthButtons } from '#/features/auth/o-auth-buttons';
import { RegisterForm } from '#/features/auth/register-form';
import { useGetPublicSettings, useGetSetupStatus } from '#/hooks/useSettings';

export const Route = createFileRoute('/_auth/signup')({
  component: RegisterPage,
});

function RegisterPage() {
  const { data: publicSettings } = useGetPublicSettings();
  const { data: setupStatus, isLoading } = useGetSetupStatus();
  const registrationEnabled = publicSettings?.data?.registrationEnabled ?? true;

  if (!isLoading && setupStatus?.data?.setupRequired) {
    return <Navigate to="/setup" replace />;
  }

  if (!isLoading && !registrationEnabled) {
    return <Navigate to="/signin" replace />;
  }

  return (
    <div className="fade-in slide-in-from-bottom-4 animate-in duration-500">
      <AuthPageHeader title="Create an account" description="Join this Vessl instance." />

      <OAuthButtons />
      <RegisterForm />

      <p className="mt-8 text-center text-sm text-white/50">
        Already have an account?{' '}
        <Link to="/signin" className="font-semibold text-primary hover:underline">
          Sign in
        </Link>
      </p>
    </div>
  );
}
