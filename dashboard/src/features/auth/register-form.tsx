import { zodResolver } from '@hookform/resolvers/zod';
import { Eye, EyeOff, Lock, Mail, User } from 'lucide-react';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useRegister } from '#/hooks/useAuth';

const registerSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.email('Please enter a valid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters long'),
});

type RegisterSchema = z.infer<typeof registerSchema>;

export const RegisterForm = () => {
  const { mutateAsync: registerUser, isPending } = useRegister();
  const [showPassword, setShowPassword] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterSchema>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      name: '',
      email: '',
      password: '',
    },
  });

  const onSubmit = async (data: RegisterSchema) => {
    try {
      await registerUser(data);
    } catch {}
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
      <div className="space-y-2">
        <Label htmlFor="name" className="font-medium text-sm">
          Full Name
        </Label>
        <div className="relative">
          <div className="absolute top-3.5 left-3 text-muted-foreground">
            <User className="h-5 w-5" />
          </div>
          <Input
            id="name"
            type="text"
            placeholder="John Doe"
            className="h-12 bg-background pl-10 text-base"
            {...register('name')}
          />
        </div>
        {errors.name && <p className="text-[13px] text-destructive">{errors.name.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="email" className="font-medium text-sm">
          Email
        </Label>
        <div className="relative">
          <div className="absolute top-3.5 left-3 text-muted-foreground">
            <Mail className="h-5 w-5" />
          </div>
          <Input
            id="email"
            type="email"
            placeholder="name@example.com"
            className="h-12 bg-background pl-10 text-base"
            {...register('email')}
          />
        </div>
        {errors.email && <p className="text-[13px] text-destructive">{errors.email.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="password" className="font-medium text-sm">
          Password
        </Label>
        <div className="relative">
          <div className="absolute top-3.5 left-3 text-muted-foreground">
            <Lock className="h-5 w-5" />
          </div>
          <Input
            id="password"
            type={showPassword ? 'text' : 'password'}
            className="h-12 bg-background pr-12 pl-10 text-base"
            {...register('password')}
          />
          <button
            type="button"
            className="absolute top-3.5 right-3 text-muted-foreground transition-colors hover:text-foreground"
            onClick={() => setShowPassword(!showPassword)}
          >
            {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
          </button>
        </div>
        {errors.password ? (
          <p className="text-[13px] text-destructive">{errors.password.message}</p>
        ) : (
          <p className="pt-1 text-muted-foreground text-xs">Must be at least 8 characters long.</p>
        )}
      </div>

      <Button type="submit" className="mt-2 h-12 w-full font-medium text-base" disabled={isPending}>
        {isPending ? 'Creating account...' : 'Create Account'}
      </Button>
    </form>
  );
};
