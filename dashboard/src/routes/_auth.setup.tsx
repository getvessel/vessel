import { createFileRoute, Navigate } from '@tanstack/react-router';
import { AuthPageHeader } from '#/features/auth/auth-page-header';
import { SetupForm } from '#/features/auth/setup-form';
import { useGetSetupStatus } from '#/hooks/useSettings';

export const Route = createFileRoute('/_auth/setup')({
  component: SetupPage,
});

function SetupPage() {
  const { data: setupStatus, isLoading } = useGetSetupStatus();

  if (!isLoading && !setupStatus?.data?.setupRequired) {
    return <Navigate to="/signin" replace />;
  }

  return (
    <div className="fade-in slide-in-from-bottom-4 animate-in duration-500">
      <AuthPageHeader
        title="Set up Vessl"
        description="Create the initial owner account for this instance."
      />

      <SetupForm />
    </div>
  );
}
