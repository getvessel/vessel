import { zodResolver } from '@hookform/resolvers/zod';
import { Lock, Mail, User } from 'lucide-react';
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
        <Label htmlFor="name" className="text-sm font-medium">
          Owner Full Name
        </Label>
        <div className="relative">
          <div className="absolute left-3 top-3.5 text-muted-foreground">
            <User className="h-5 w-5" />
          </div>
          <Input
            id="name"
            type="text"
            placeholder="John Doe"
            className="h-12 pl-10 bg-background text-base"
            {...register('name')}
          />
        </div>
        {errors.name && <p className="text-[13px] text-destructive">{errors.name.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="email" className="text-sm font-medium">
          Owner Email
        </Label>
        <div className="relative">
          <div className="absolute left-3 top-3.5 text-muted-foreground">
            <Mail className="h-5 w-5" />
          </div>
          <Input
            id="email"
            type="email"
            placeholder="name@example.com"
            className="h-12 pl-10 bg-background text-base"
            {...register('email')}
          />
        </div>
        {errors.email && <p className="text-[13px] text-destructive">{errors.email.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="password" className="text-sm font-medium">
          Owner Password
        </Label>
        <div className="relative">
          <div className="absolute left-3 top-3.5 text-muted-foreground">
            <Lock className="h-5 w-5" />
          </div>
          <Input
            id="password"
            type="password"
            className="h-12 pl-10 bg-background text-base"
            {...register('password')}
          />
        </div>
        {errors.password && (
          <p className="text-[13px] text-destructive">{errors.password.message}</p>
        )}
      </div>

      <Button type="submit" className="w-full h-12 text-base" disabled={isPending}>
        {isPending ? 'Creating Owner Account...' : 'Complete Setup'}
      </Button>
    </form>
  );
};
