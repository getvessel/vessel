import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { AuthPageHeader } from '#/features/auth/auth-page-header';
import { ResetPasswordForm } from '#/features/auth/reset-password-form';

export const Route = createFileRoute('/_auth/reset-password')({
  validateSearch: (search: Record<string, unknown>) => {
    return {
      token: (search.token as string) || '',
    };
  },
  component: ResetPasswordPage,
});

function ResetPasswordPage() {
  const { token } = Route.useSearch();
  const navigate = useNavigate();

  if (!token) {
    return (
      <div className="fade-in slide-in-from-bottom-4 animate-in text-center duration-500">
        <h1 className="mb-4 font-semibold text-2xl text-white tracking-tight">Invalid Request</h1>
        <p className="mb-8 text-sm text-white/55">
          The password reset token is missing. Please check your email link again.
        </p>
        <button
          type="button"
          onClick={() => navigate({ to: '/signin' })}
          className="font-medium text-primary text-sm hover:underline"
        >
          Return to sign in
        </button>
      </div>
    );
  }

  return (
    <div className="fade-in slide-in-from-bottom-4 animate-in duration-500">
      <AuthPageHeader
        title="Create new password"
        description="Use a password that is different from previous ones."
      />

      <ResetPasswordForm token={token} />
    </div>
  );
}
