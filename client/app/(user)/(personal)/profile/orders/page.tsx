'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function OrdersProfileRedirect() {
  const router = useRouter();
  
  useEffect(() => {
    router.replace('/orders');
  }, [router]);

  return (
    <div className="flex justify-center items-center p-8">
      <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-indigo-600"></div>
    </div>
  );
}
