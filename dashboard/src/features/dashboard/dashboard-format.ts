export function formatMb(value: number) {
  if (value >= 1024) {
    return `${(value / 1024).toFixed(1)} GB`;
  }

  return `${Math.round(value)} MB`;
}

export function formatGb(value: number) {
  return `${value.toFixed(1)} GB`;
}

export function formatDuration(seconds: number) {
  const days = Math.floor(seconds / 86_400);
  const hours = Math.floor((seconds % 86_400) / 3_600);
  const minutes = Math.floor((seconds % 3_600) / 60);

  if (days > 0) {
    return `${days}d ${hours}h`;
  }

  if (hours > 0) {
    return `${hours}h ${minutes}m`;
  }

  return `${minutes}m`;
}
