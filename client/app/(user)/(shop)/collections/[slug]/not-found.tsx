import Link from 'next/link';

export default function CollectionNotFound() {
  return (
    <div className="container mx-auto px-4 py-16 flex flex-col items-center justify-center min-h-[60vh]">
      <h2 className="text-3xl font-bold text-gray-800 mb-4">Collection Not Found</h2>
      <p className="text-gray-600 mb-8 text-center max-w-md">
        We couldn't find the collection you're looking for. It may have been removed or renamed.
      </p>
      <div className="flex gap-4">
        <Link
          href="/collections"
          className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          View All Collections
        </Link>
        <Link
          href="/"
          className="px-6 py-3 bg-gray-200 text-gray-800 rounded-lg hover:bg-gray-300 transition-colors"
        >
          Go Home
        </Link>
      </div>
    </div>
  );
}
