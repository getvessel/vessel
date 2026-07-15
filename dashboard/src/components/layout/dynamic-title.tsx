import { useEffect } from 'react';
import { useGetPublicSettings } from '#/hooks/useSettings';

export function DynamicTitle() {
  const { data } = useGetPublicSettings();
  const siteName = data?.data?.siteName;

  useEffect(() => {
    if (siteName) {
      document.title = `${siteName} Dashboard`;
    } else {
      document.title = 'Vessl Dashboard';
    }
  }, [siteName]);

  return null;
}
