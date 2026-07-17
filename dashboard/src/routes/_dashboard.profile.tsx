import { createFileRoute } from '@tanstack/react-router';
import { Security2FASetup } from '#/features/profile/security-2fa-setup';
import {
  ProfileEmailForm,
  ProfileNameForm,
  ProfilePasswordForm,
} from '#/features/profile/user-profile-form';

export const Route = createFileRoute('/_dashboard/profile')({
  component: ProfilePage,
});

function ProfilePage() {
  return (
    <div className="space-y-6">
      <div className="flex flex-col justify-between gap-6 pb-2 md:flex-row md:items-start">
        <div className="flex-1 space-y-4">
          <div className="space-y-1">
            <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              ACCOUNT MANAGEMENT
            </p>
            <h1 className="font-bold text-3xl tracking-tight">Profile & Security</h1>
          </div>
          <p className="max-w-2xl text-muted-foreground text-sm leading-relaxed">
            Manage your personal profile and security preferences.
          </p>
        </div>
      </div>

      <div className="grid gap-6">
        <ProfileNameForm />
        <ProfileEmailForm />
        <ProfilePasswordForm />
        <Security2FASetup />
      </div>
    </div>
  );
}
