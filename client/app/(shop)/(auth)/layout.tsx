import Image from 'next/image';

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className='bg-gray-50 h-full p-10'>
      <div className='container m-auto flex gap-20'>
        <div className='w-1/2 flex flex-col'>{children}</div>
        <div className='w-1/2'>
          <Image
            src={'/images/illustration.svg'}
            height={2000}
            width={1000}
            alt='Register Image'
          />
        </div>
      </div>
    </div>
  );
}
