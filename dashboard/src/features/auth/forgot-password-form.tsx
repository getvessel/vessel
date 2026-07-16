import { Link } from '@tanstack/react-router';
import { Loader2, Mail } from 'lucide-react';
import type { SyntheticEvent } from 'react';
import { useState } from 'react';

import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';

import { useForgotPassword } from '#/hooks/useAuth';

export const ForgotPasswordForm = () => {
  const [email, setEmail] = useState('');
  const { mutate, isPending, isSuccess } = useForgotPassword();

  const handleSubmit = (e: SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!email) return;
    mutate(email);
  };

  if (isSuccess) {
    return (
      <div className="space-y-4 text-center">
        <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
          <Mail className="h-6 w-6 text-primary" />
        </div>
        <h3 className="font-medium text-xl">Check your email</h3>
        <p className="text-sm text-white/55">
          If an account with that email exists, we've sent you instructions to reset your password.
        </p>
        <div className="mt-8">
          <Link to="/signin" className="font-medium text-primary text-sm hover:underline">
            Back to sign in
          </Link>
        </div>
      </div>
    );
  }

  return (
    <>
      <form onSubmit={handleSubmit} className="space-y-5">
        <div className="space-y-2">
          <Label htmlFor="email" className="font-medium text-sm">
            Email
          </Label>
          <div className="relative">
            <div className="absolute top-2.5 left-3 text-white/35">
              <Mail className="h-4 w-4" />
            </div>
            <Input
              id="email"
              type="email"
              placeholder="name@example.com"
              autoComplete="email"
              className="h-10 border-white/10 bg-white/[0.055] pl-9 text-white placeholder:text-white/35"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={isPending}
            />
          </div>
        </div>

        <Button
          type="submit"
          className="mt-2 h-10 w-full font-medium"
          disabled={isPending || !email}
        >
          {isPending && <Loader2 className="size-4 animate-spin" />}
          {isPending ? 'Sending link' : 'Send reset link'}
        </Button>
      </form>

      <div className="mt-8 text-center text-sm">
        <span className="text-white/50">Remember your password? </span>
        <Link to="/signin" className="font-medium text-primary hover:underline">
          Sign in
        </Link>
      </div>
    </>
  );
};
