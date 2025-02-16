'use client';
import { useEffect } from 'react';
import { Bounce, toast } from 'react-toastify';

export default function ErrorPage({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error(error);
    toast.warn('Wow so easy!', {
      position: 'top-center',
      autoClose: 5000,
      theme: 'colored',
      transition: Bounce,
    });
  }, []);

  return (
    <div>
      <h2>Something went wrong!</h2>
      <button onClick={() => reset()}>Try again</button>
    </div>
  );
}
