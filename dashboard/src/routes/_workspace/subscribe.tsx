import { createFileRoute } from '@tanstack/react-router';
import { Check } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Tabs, TabsList, TabsTrigger } from '#/components/ui/tabs';

export const Route = createFileRoute('/_workspace/subscribe')({
  component: SubscriptionNewPage,
});

function SubscriptionNewPage() {
  return (
    <div className="flex-1 space-y-6 max-w-4xl mx-auto py-8 px-4">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight text-foreground">Subscriptions</h1>
      </div>

      <div className="flex flex-col items-center mt-8">
        <Tabs defaultValue="monthly" className="w-100 flex justify-center mb-10">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="monthly">Monthly</TabsTrigger>
            <TabsTrigger value="annually">Annually (save ~20%)</TabsTrigger>
          </TabsList>
        </Tabs>

        <Card className="w-full max-w-lg border-border bg-card shadow-sm">
          <CardHeader>
            <CardTitle className="text-2xl font-bold">Pay-as-you-go</CardTitle>
            <CardDescription className="text-muted-foreground mt-2">
              Dynamic pricing based on the number of servers you connect.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="flex items-baseline gap-2">
              <span className="text-4xl font-extrabold tracking-tight">$5</span>
              <span className="text-muted-foreground font-medium">/ mo base</span>
            </div>
            <p className="text-sm font-medium text-muted-foreground">
              + <span className="text-foreground font-bold">$3</span> per additional server, billed
              monthly (+VAT)
            </p>

            <Button size="lg" className="w-full font-semibold">
              Subscribe
            </Button>

            <div className="pt-6 space-y-3">
              {[
                'Connect unlimited servers',
                'Deploy unlimited applications per server',
                'Free email notifications',
                'Support by email',
                '+ All Upcoming Features',
              ].map((feature, i) => (
                <div key={i} className="flex items-center gap-3">
                  <Check className="h-4 w-4 text-primary" />
                  <span className="text-sm font-medium text-card-foreground">{feature}</span>
                </div>
              ))}
            </div>

            <div className="pt-6 border-t border-border mt-6">
              <p className="text-xs text-muted-foreground leading-relaxed">
                You need to bring your own servers from any cloud provider (Hetzner, DigitalOcean,
                AWS, etc.) or connect any device running a supported OS.
              </p>
              <p className="text-xs text-muted-foreground mt-3">
                Need official support for your self-hosted instance?{' '}
                <a href="mailto:support@vessl.dev" className="underline hover:text-foreground">
                  Contact Us
                </a>
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
