import { zodResolver } from '@hookform/resolvers/zod';
import { Eye, EyeOff, Loader2, Lock, Mail, User } from 'lucide-react';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useSetup } from '#/hooks/useAuth';

const setupSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.email('Please enter a valid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters long'),
});

type SetupSchema = z.infer<typeof setupSchema>;

export const SetupForm = () => {
  const { mutateAsync: setupUser, isPending } = useSetup();
  const [showPassword, setShowPassword] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SetupSchema>({
    resolver: zodResolver(setupSchema),
    defaultValues: {
      name: '',
      email: '',
      password: '',
    },
  });

  const onSubmit = async (data: SetupSchema) => {
    try {
      await setupUser(data);
    } catch {}
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
      <div className="space-y-2">
        <Label htmlFor="name" className="font-medium text-sm text-white/80">
          Owner name
        </Label>
        <div className="relative">
          <div className="absolute top-2.5 left-3 text-white/35">
            <User className="h-4 w-4" />
          </div>
          <Input
            id="name"
            type="text"
            placeholder="John Doe"
            autoComplete="name"
            aria-invalid={!!errors.name}
            className="h-10 border-white/10 bg-white/[0.055] pl-9 text-white placeholder:text-white/35"
            {...register('name')}
          />
        </div>
        {errors.name && <p className="text-[13px] text-destructive">{errors.name.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="email" className="font-medium text-sm text-white/80">
          Owner email
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
        <Label htmlFor="password" className="font-medium text-sm text-white/80">
          Owner password
        </Label>
        <div className="relative">
          <div className="absolute top-2.5 left-3 text-white/35">
            <Lock className="h-4 w-4" />
          </div>
          <Input
            id="password"
            type={showPassword ? 'text' : 'password'}
            autoComplete="new-password"
            aria-invalid={!!errors.password}
            className="h-10 border-white/10 bg-white/[0.055] pr-10 pl-9 text-white"
            {...register('password')}
          />
          <button
            type="button"
            aria-label={showPassword ? 'Hide password' : 'Show password'}
            className="absolute top-2.5 right-3 text-white/35 hover:text-white"
            onClick={() => setShowPassword(!showPassword)}
          >
            {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
          </button>
        </div>
        {errors.password && (
          <p className="text-[13px] text-destructive">{errors.password.message}</p>
        )}
      </div>

      <Button type="submit" className="h-10 w-full font-medium" disabled={isPending}>
        {isPending && <Loader2 className="size-4 animate-spin" />}
        {isPending ? 'Creating owner account' : 'Complete setup'}
      </Button>
    </form>
  );
};
