import { createFileRoute } from '@tanstack/react-router';
import { AuthPageHeader } from '#/features/auth/auth-page-header';
import { ForgotPasswordForm } from '#/features/auth/forgot-password-form';

export const Route = createFileRoute('/_auth/forgot-password')({
  component: ForgotPasswordPage,
});

import { AlertCircle } from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from '#/components/ui/alert';
import { useGetPublicSettings } from '#/hooks/useSettings';

function ForgotPasswordPage() {
  const { data, isLoading } = useGetPublicSettings();
  const emailEnabled = data?.data?.emailEnabled;

  return (
    <div className="fade-in slide-in-from-bottom-4 animate-in duration-500">
      <AuthPageHeader title="Reset your password" description="Receive a recovery link by email." />

      {!isLoading && emailEnabled === false ? (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Email not configured</AlertTitle>
          <AlertDescription>
            Your team is yet to set or enable email. Please contact your administrator.
          </AlertDescription>
        </Alert>
      ) : (
        <ForgotPasswordForm />
      )}
    </div>
  );
}
