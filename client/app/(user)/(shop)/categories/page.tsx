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

async function getCategories(page = 1, pageSize = 9) {
  // Using apiFetch utility with pagination
  const result = await serverSideFetch<GeneralCategoryModel[]>(
    `${PUBLIC_API_PATHS.CATEGORIES}?page=${page}&pageSize=${pageSize}`
  );

  return {
    data: result.data || [],
    pagination: result.pagination || { total: 0, page, pageSize: pageSize },
  };
}

export default async function CategoryPage({
  searchParams,
}: {
  searchParams: SearchParams;
}) {
  const currentPage = searchParams.page ? Number(searchParams.page) : 1;
  const pageSize = 9; // 3x3 grid

  // Server component with async data fetching
  const { data: categories, pagination } = await getCategories(
    currentPage,
    pageSize
  );

  // Calculate total pages
  const totalPages = Math.ceil(pagination.total / pageSize);

  return (
    <div className='container mx-auto px-4 py-8'>
      <div className='mb-8 text-center'>
        <h1 className='text-3xl font-bold'>Categories</h1>
        <p className='text-gray-500 mt-2'>Explore our product collections</p>
      </div>

      {/* Categories Grid */}
      <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'>
        {categories.map((category) => (
          <Link
            href={`/categories/${category.slug}`}
            key={category.id}
            className='group'
          >
            <div className='bg-white rounded-lg overflow-hidden shadow-md hover:shadow-lg transition-shadow duration-300 h-full flex flex-col'>
              <div className='relative h-96 w-full bg-gray-100'>
                {category.imageUrl ? (
                  <Image
                    src={category.imageUrl}
                    alt={category.name}
                    fill
                    sizes='(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw'
                    className='object-cover group-hover:scale-105 transition-transform duration-300'
                  />
                ) : (
                  <div className='w-full h-full flex items-center justify-center bg-gray-200'>
                    <span className='text-gray-400'>No image</span>
                  </div>
                )}
              </div>
              <div className='p-6 flex flex-col flex-grow'>
                <h2 className='text-xl font-bold text-gray-800 mb-3 group-hover:text-indigo-600 transition-colors'>
                  {category.name}
                </h2>
                <p className='text-gray-600 text-sm line-clamp-3 flex-grow'>
                  {category.description || 'No description available'}
                </p>
                <div className='mt-4 text-indigo-600 text-sm font-medium flex items-center opacity-0 group-hover:opacity-100 transition-opacity'>
                  <span>Browse products</span>
                  <svg
                    className='w-4 h-4 ml-1'
                    fill='none'
                    stroke='currentColor'
                    viewBox='0 0 24 24'
                    xmlns='http://www.w3.org/2000/svg'
                  >
                    <path
                      strokeLinecap='round'
                      strokeLinejoin='round'
                      strokeWidth={2}
                      d='M14 5l7 7m0 0l-7 7m7-7H3'
                    />
                  </svg>
                </div>
              </div>
            </div>
          </Link>
        ))}
      </div>

      {categories.length === 0 && (
        <div className='text-center py-12'>
          <p className='text-xl text-gray-600'>No categories found.</p>
        </div>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className='flex items-center justify-center space-x-2 mt-12'>
          <Link
            href={`/categories?page=${Math.max(1, currentPage - 1)}`}
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
                href={`/categories?page=${page}`}
                className={`px-3 py-2 rounded-md ${
                  currentPage === page
                    ? 'bg-indigo-600 text-white'
                    : 'text-gray-700 hover:bg-gray-100'
                }`}
              >
                {page}
              </Link>
            ))}
          </div>

          <Link
            href={`/categories?page=${Math.min(totalPages, currentPage + 1)}`}
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
