'use client';

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  console.log(error);
  return (
    <div className='flex items-center justify-center h-screen'>
      <h2 className='text-lg font-bold text-red-600'>Something went wrong!</h2>
      <button
        className='btn primary text-white'
        onClick={
          // Attempt to recover by trying to re-render the segment
          () => reset()
        }
      >
        Try again
      </button>
      <div className='text-red-500 text-4xl'>Error</div>
    </div>
  );
}
