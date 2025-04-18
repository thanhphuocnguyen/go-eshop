import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';

export default function NotFound() {
  return (
    <div className='not-found-bg'>
      <h3 className='mt-32'>404</h3>
      <h1 className='text-4xl font-bold'>Page not found</h1>
      <p className='text-lg my-4 text-white/60'>
        Sorry, we couldn’t find the page you’re looking for.
      </p>

      <Link className='mt-6 hover:underline items-center gap-2 flex' href='/'>
        <ArrowLeftIcon className='size-5' />
        Back to home
      </Link>
    </div>
  );
}
