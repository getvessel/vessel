import { Loader2, Lock } from 'lucide-react';
import { type SyntheticEvent, useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useResetPassword } from '#/hooks/useAuth';

interface ResetPasswordFormProps {
  token: string;
}

export const ResetPasswordForm = ({ token }: ResetPasswordFormProps) => {
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');

  const { mutate, isPending } = useResetPassword();

  const handleSubmit = (e: SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!newPassword || !confirmPassword) return;

    if (newPassword !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }
    setError('');

    mutate({ token, newPassword });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      <div className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="new-password" className="font-medium text-sm text-white/80">
            New password
          </Label>
          <div className="relative">
            <div className="absolute top-2.5 left-3 text-white/35">
              <Lock className="h-4 w-4" />
            </div>
            <Input
              id="new-password"
              type="password"
              placeholder="Enter new password"
              autoComplete="new-password"
              aria-invalid={!!error}
              className="h-10 border-white/10 bg-white/[0.055] pl-9 text-white placeholder:text-white/35"
              value={newPassword}
              onChange={(e) => {
                setNewPassword(e.target.value);
                setError('');
              }}
              required
              disabled={isPending}
            />
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="confirm-password" className="font-medium text-sm text-white/80">
            Confirm password
          </Label>
          <div className="relative">
            <div className="absolute top-2.5 left-3 text-white/35">
              <Lock className="h-4 w-4" />
            </div>
            <Input
              id="confirm-password"
              type="password"
              placeholder="Confirm new password"
              autoComplete="new-password"
              aria-invalid={!!error}
              className="h-10 border-white/10 bg-white/[0.055] pl-9 text-white placeholder:text-white/35"
              value={confirmPassword}
              onChange={(e) => {
                setConfirmPassword(e.target.value);
                setError('');
              }}
              required
              disabled={isPending}
            />
          </div>
        </div>

        {error && <p className="font-medium text-destructive text-sm">{error}</p>}
      </div>

      <Button
        type="submit"
        className="mt-2 h-10 w-full font-medium"
        disabled={isPending || !newPassword || !confirmPassword}
      >
        {isPending && <Loader2 className="size-4 animate-spin" />}
        {isPending ? 'Resetting password' : 'Reset password'}
      </Button>
    </form>
  );
};
