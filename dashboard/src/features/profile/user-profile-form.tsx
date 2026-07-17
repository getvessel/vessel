import { Loader2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { InputOTP, InputOTPGroup, InputOTPSlot } from '#/components/ui/input-otp';
import { Label } from '#/components/ui/label';
import {
  useChangePassword,
  useGetProfile,
  useRequestEmailChange,
  useUpdateProfile,
  useVerifyEmailChange,
} from '#/hooks/useProfile';

export function ProfileNameForm() {
  const { data: profile, isLoading } = useGetProfile();
  const updateProfile = useUpdateProfile();
  const [name, setName] = useState('');

  // Sync state once profile loads
  if (!isLoading && profile?.data && !name && profile.data.name) {
    setName(profile.data.name);
  }

  const handleSave = () => {
    updateProfile.mutate(
      { name },
      {
        onSuccess: () => toast.success('Profile name updated!'),
        onError: (err) => toast.error(err.message),
      }
    );
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Profile Name</CardTitle>
        <CardDescription>Update your display name.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="name">Full Name</Label>
          <Input
            id="name"
            placeholder="John Doe"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>
        <Button
          onClick={handleSave}
          disabled={isLoading || updateProfile.isPending || name === profile?.data?.name}
        >
          {updateProfile.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          Save Name
        </Button>
      </CardContent>
    </Card>
  );
}

export function ProfileEmailForm() {
  const { data: profile, isLoading } = useGetProfile();
  const requestEmailChange = useRequestEmailChange();
  const verifyEmailChange = useVerifyEmailChange();

  const [email, setEmail] = useState('');
  const [otpOpen, setOtpOpen] = useState(false);
  const [otp, setOtp] = useState('');

  // Sync state once profile loads
  if (!isLoading && profile?.data && !email && profile.data.email) {
    setEmail(profile.data.email);
  }

  const handleRequest = () => {
    requestEmailChange.mutate(
      { newEmail: email },
      {
        onSuccess: () => {
          setOtpOpen(true);
          toast.success('Verification code sent to your new email.');
        },
        onError: (err) => toast.error(err.message),
      }
    );
  };

  const handleVerify = () => {
    verifyEmailChange.mutate(
      { otp },
      {
        onSuccess: () => {
          setOtpOpen(false);
          setOtp('');
          toast.success('Email updated successfully!');
        },
        onError: (err) => toast.error(err.message),
      }
    );
  };

  return (
    <>
      <Card>
        <CardHeader>
          <CardTitle>Email Address</CardTitle>
          <CardDescription>Change the email address associated with your account.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </div>
          <Button
            onClick={handleRequest}
            disabled={isLoading || requestEmailChange.isPending || email === profile?.data?.email}
          >
            {requestEmailChange.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Request Email Change
          </Button>
        </CardContent>
      </Card>

      <Dialog open={otpOpen} onOpenChange={setOtpOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Verify Email Change</DialogTitle>
            <DialogDescription>
              We've sent a 6-digit verification code to your new email. Enter it below.
            </DialogDescription>
          </DialogHeader>
          <div className="flex justify-center py-6">
            <InputOTP maxLength={6} value={otp} onChange={setOtp}>
              <InputOTPGroup>
                <InputOTPSlot index={0} />
                <InputOTPSlot index={1} />
                <InputOTPSlot index={2} />
                <InputOTPSlot index={3} />
                <InputOTPSlot index={4} />
                <InputOTPSlot index={5} />
              </InputOTPGroup>
            </InputOTP>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setOtpOpen(false)}>
              Cancel
            </Button>
            <Button
              onClick={handleVerify}
              disabled={otp.length !== 6 || verifyEmailChange.isPending}
            >
              {verifyEmailChange.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Verify & Update
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}

export function ProfilePasswordForm() {
  const changePassword = useChangePassword();
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');

  const handleSave = () => {
    changePassword.mutate(
      { oldPassword, newPassword },
      {
        onSuccess: () => {
          toast.success('Password updated successfully!');
          setOldPassword('');
          setNewPassword('');
        },
        onError: (err) => toast.error(err.message),
      }
    );
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Change Password</CardTitle>
        <CardDescription>Update your account password.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="oldPassword">Current Password</Label>
          <Input
            id="oldPassword"
            type="password"
            value={oldPassword}
            onChange={(e) => setOldPassword(e.target.value)}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="newPassword">New Password</Label>
          <Input
            id="newPassword"
            type="password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
          />
        </div>
        <Button
          onClick={handleSave}
          disabled={changePassword.isPending || !oldPassword || !newPassword}
        >
          {changePassword.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          Update Password
        </Button>
      </CardContent>
    </Card>
  );
}
