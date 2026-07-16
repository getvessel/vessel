export function AuthPageHeader({ title, description }: { title: string; description: string }) {
  return (
    <div className="mb-7 flex flex-col space-y-2">
      <div className="inline-flex w-fit items-center rounded-full border border-white/10 bg-white/[0.04] px-2.5 py-1 text-white/55 text-xs">
        Secure access
      </div>
      <h1 className="font-semibold text-2xl text-white tracking-tight">{title}</h1>
      <p className="text-sm text-white/55">{description}</p>
    </div>
  );
}
