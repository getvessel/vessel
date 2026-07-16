import { zodResolver } from '@hookform/resolvers/zod';
import { Link } from '@tanstack/react-router';
import { Eye, EyeOff, Loader2, Lock, Mail } from 'lucide-react';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useLogin } from '#/hooks/useAuth';

const loginSchema = z.object({
  email: z.email('Please enter a valid email address'),
  password: z.string().min(1, 'Password is required'),
});

type LoginSchema = z.infer<typeof loginSchema>;

export const LoginForm = () => {
  const { mutateAsync: login, isPending } = useLogin();
  const [showPassword, setShowPassword] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginSchema>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  });

  const onSubmit = async (data: LoginSchema) => {
    try {
      await login(data);
    } catch {}
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
      <div className="space-y-2">
        <Label htmlFor="email" className="font-medium text-sm text-white/80">
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
            aria-invalid={!!errors.email}
            className="h-10 border-white/10 bg-white/[0.055] pl-9 text-white placeholder:text-white/35"
            {...register('email')}
          />
        </div>
        {errors.email && <p className="text-[13px] text-destructive">{errors.email.message}</p>}
      </div>

      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <Label htmlFor="password" className="font-medium text-sm text-white/80">
            Password
          </Label>
          <Link to="/forgot-password" className="font-medium text-primary text-sm hover:underline">
            Forgot password?
          </Link>
        </div>
        <div className="relative">
          <div className="absolute top-2.5 left-3 text-white/35">
            <Lock className="h-4 w-4" />
          </div>
          <Input
            id="password"
            type={showPassword ? 'text' : 'password'}
            autoComplete="current-password"
            aria-invalid={!!errors.password}
            className="h-10 border-white/10 bg-white/[0.055] pr-10 pl-9 text-white"
            {...register('password')}
          />
          <button
            type="button"
            aria-label={showPassword ? 'Hide password' : 'Show password'}
            className="absolute top-2.5 right-3 text-white/35 transition-colors hover:text-white"
            onClick={() => setShowPassword(!showPassword)}
          >
            {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
          </button>
        </div>
        {errors.password && (
          <p className="text-[13px] text-destructive">{errors.password.message}</p>
        )}
      </div>

      <Button type="submit" className="mt-2 h-10 w-full font-medium" disabled={isPending}>
        {isPending && <Loader2 className="size-4 animate-spin" />}
        {isPending ? 'Signing in' : 'Sign in'}
      </Button>
    </form>
  );
};
