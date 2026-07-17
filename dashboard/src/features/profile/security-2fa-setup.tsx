import { Loader2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
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
import { InputOTP, InputOTPGroup, InputOTPSlot } from '#/components/ui/input-otp';
import { useDisable2FA, useSetup2FA, useVerify2FA } from '#/hooks/useAuth';
import { useGetProfile } from '#/hooks/useProfile';

export function Security2FASetup() {
  const { data: profile, isLoading } = useGetProfile();

  const setup2FA = useSetup2FA();
  const verify2FA = useVerify2FA();
  const disable2FA = useDisable2FA();

  const [setupData, setSetupData] = useState<{ qrCodeUrl: string; secret: string } | null>(null);
  const [verifyOpen, setVerifyOpen] = useState(false);
  const [disableOpen, setDisableOpen] = useState(false);
  const [otp, setOtp] = useState('');

  const isEnabled = profile?.data?.totpEnabled;

  const handleEnableClick = () => {
    setup2FA.mutate(undefined, {
      onSuccess: (res) => {
        setSetupData(res.data);
        setOtp('');
        setVerifyOpen(true);
      },
      onError: (err) => toast.error(err.message),
    });
  };

  const handleDisableClick = () => {
    setOtp('');
    setDisableOpen(true);
  };

  const handleVerify = () => {
    verify2FA.mutate(
      { token: otp },
      {
        onSuccess: () => {
          setVerifyOpen(false);
          setSetupData(null);
          toast.success('Two-factor authentication enabled successfully.');
        },
        onError: (err) => toast.error(err.message),
      }
    );
  };

  const handleDisable = () => {
    disable2FA.mutate(
      { token: otp },
      {
        onSuccess: () => {
          setDisableOpen(false);
          toast.success('Two-factor authentication disabled.');
        },
        onError: (err) => toast.error(err.message),
      }
    );
  };

  if (isLoading) {
    return (
      <Card>
        <CardContent className="flex justify-center py-10">
          <Loader2 className="animate-spin text-muted-foreground" />
        </CardContent>
      </Card>
    );
  }

  return (
    <>
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Two-Factor Authentication</CardTitle>
              <CardDescription>Add an extra layer of security to your account.</CardDescription>
            </div>
            {isEnabled ? (
              <Badge
                variant="default"
                className="bg-green-500/10 text-green-500 hover:bg-green-500/20"
              >
                Enabled
              </Badge>
            ) : (
              <Badge variant="secondary">Disabled</Badge>
            )}
          </div>
        </CardHeader>
        <CardContent>
          <p className="mb-4 text-muted-foreground text-sm">
            {isEnabled
              ? 'Two-factor authentication is currently enabled. You will need to enter a code from your authenticator app when signing in.'
              : 'Protect your account from unauthorized access by requiring a second authentication method in addition to your password.'}
          </p>
          {isEnabled ? (
            <Button variant="destructive" onClick={handleDisableClick}>
              Disable 2FA
            </Button>
          ) : (
            <Button onClick={handleEnableClick} disabled={setup2FA.isPending}>
              {setup2FA.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Enable 2FA
            </Button>
          )}
        </CardContent>
      </Card>

      <Dialog open={verifyOpen} onOpenChange={setVerifyOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Setup Two-Factor Authentication</DialogTitle>
            <DialogDescription>
              Scan the QR code below with your authenticator app (like Google Authenticator or
              Authy), then enter the 6-digit code.
            </DialogDescription>
          </DialogHeader>
          <div className="flex flex-col items-center space-y-6 py-4">
            {setupData && (
              <div className="rounded-lg bg-white p-4">
                <img src={setupData.qrCodeUrl} alt="2FA QR Code" className="h-48 w-48" />
              </div>
            )}
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
            <Button variant="outline" onClick={() => setVerifyOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleVerify} disabled={otp.length !== 6 || verify2FA.isPending}>
              {verify2FA.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Verify & Enable
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={disableOpen} onOpenChange={setDisableOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Disable Two-Factor Authentication</DialogTitle>
            <DialogDescription>
              Enter a 6-digit code from your authenticator app to confirm you want to disable 2FA.
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
            <Button variant="outline" onClick={() => setDisableOpen(false)}>
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleDisable}
              disabled={otp.length !== 6 || disable2FA.isPending}
            >
              {disable2FA.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Disable
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
