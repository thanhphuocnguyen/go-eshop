import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { GeneralCategoryModel } from '@/app/lib/definitions';
import Link from 'next/link';
import Image from 'next/image';
import { serverSideFetch } from '@/app/lib/api/apiServer';
import { ChevronLeftIcon, ChevronRightIcon } from '@heroicons/react/24/outline';

// Define the search params type
type SearchParams = {
  page?: string;
};

async function getCollections(page = 1, pageSize = 9) {
  // Using apiFetch utility with pagination
  const result = await serverSideFetch<GeneralCategoryModel[]>(
    `${PUBLIC_API_PATHS.COLLECTIONS}?page=${page}&pageSize=${pageSize}`
  );

  return {
    data: result.data || [],
    pagination: result.pagination || { total: 0, page, pageSize: pageSize },
  };
}

export default async function CollectionsPage({
  searchParams,
}: {
  searchParams: SearchParams;
}) {
  const currentPage = searchParams.page ? Number(searchParams.page) : 1;
  const pageSize = 9; // 3x3 grid

  // Server component with async data fetching
  const { data: collections, pagination } = await getCollections(
    currentPage,
    pageSize
  );

  // Calculate total pages
  const totalPages = Math.ceil(pagination.total / pageSize);

  return (
    <div className='container mx-auto px-4 py-8'>
      <div className='mb-8 text-center'>
        <h1 className='text-3xl font-bold'>Collections</h1>
        <p className='text-gray-500 mt-2'>Discover our curated collections</p>
      </div>

      {/* Collections Grid */}
      <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8'>
        {collections.map((collection) => (
          <Link
            href={`/collections/${collection.slug}`}
            key={collection.id}
            className='group block'
          >
            <div className='bg-white rounded-xl shadow-lg overflow-hidden hover:shadow-xl transition-all duration-300 transform group-hover:scale-[1.02]'>
              <div className='relative h-[42rem] w-full'>
                <Image
                  fill
                  alt={`${collection.name} collection`}
                  className='object-cover'
                  src={collection.imageUrl ?? '/images/product-placeholder.webp'}
                  sizes='(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw'
                />
                <div className='absolute inset-0 bg-gradient-to-t from-black/70 via-black/40 to-transparent opacity-60 group-hover:opacity-70 transition-opacity duration-300'></div>
                
                <div className='absolute bottom-0 left-0 right-0 p-6 text-left z-10'>
                  <h3 className='text-2xl font-bold text-white mb-2 group-hover:text-blue-200 transition-colors drop-shadow-lg'>
                    {collection.name}
                  </h3>
                  <p className='text-gray-200 line-clamp-2 group-hover:text-gray-100 transition-colors'>
                    {collection.description || 'Explore this collection'}
                  </p>
                  
                  <div className='mt-4 inline-flex items-center text-white opacity-0 group-hover:opacity-100 transition-all duration-300 transform translate-y-2 group-hover:translate-y-0'>
                    <span className='font-medium mr-2'>Explore Collection</span>
                    <svg className='w-5 h-5' fill='none' stroke='currentColor' viewBox='0 0 24 24' xmlns='http://www.w3.org/2000/svg'>
                      <path strokeLinecap='round' strokeLinejoin='round' strokeWidth='2' d='M14 5l7 7m0 0l-7 7m7-7H3'></path>
                    </svg>
                  </div>
                </div>
              </div>
            </div>
          </Link>
        ))}
      </div>

      {collections.length === 0 && (
        <div className='text-center py-12'>
          <p className='text-xl text-gray-600'>No collections found.</p>
        </div>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className='flex items-center justify-center space-x-2 mt-12'>
          <Link
            href={`/collections?page=${Math.max(1, currentPage - 1)}`}
            className={`flex items-center px-3 py-2 rounded-md ${
              currentPage === 1
                ? 'text-gray-400 cursor-not-allowed'
                : 'text-gray-700 hover:bg-gray-100'
            }`}
            aria-disabled={currentPage === 1}
          >
            <ChevronLeftIcon className='h-5 w-5 mr-1' />
            Previous
          </Link>

          <div className='flex items-center space-x-1'>
            {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
              <Link
                key={page}
                href={`/collections?page=${page}`}
                className={`px-3 py-2 rounded-md ${
                  currentPage === page
                    ? 'bg-blue-600 text-white'
                    : 'text-gray-700 hover:bg-gray-100'
                }`}
              >
                {page}
              </Link>
            ))}
          </div>

          <Link
            href={`/collections?page=${Math.min(totalPages, currentPage + 1)}`}
            className={`flex items-center px-3 py-2 rounded-md ${
              currentPage === totalPages
                ? 'text-gray-400 cursor-not-allowed'
                : 'text-gray-700 hover:bg-gray-100'
            }`}
            aria-disabled={currentPage === totalPages}
          >
            Next
            <ChevronRightIcon className='h-5 w-5 ml-1' />
          </Link>
        </div>
      )}
    </div>
  );
}
