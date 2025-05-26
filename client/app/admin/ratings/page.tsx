'use client';

import { useState, useEffect } from 'react';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { Button } from '@headlessui/react';
import Loading from '@/app/loading';
import { toast } from 'react-toastify';
import {
  MagnifyingGlassIcon,
  CheckCircleIcon,
  XCircleIcon,
  TrashIcon,
} from '@heroicons/react/24/outline';
import { ConfirmDialog } from '@/app/components/Common/Dialogs/ConfirmDialog';
import dayjs from 'dayjs';

// Extended rating model for admin page with additional fields
type AdminRatingModel = {
  id: string;
  productId: string;
  productName?: string;
  productImage?: string;
  userId: string;
  userName: string;
  orderItemId: string;
  rating: number;
  reviewTitle: string;
  reviewContent: string;
  verifiedPurchase: boolean;
  isVisible: boolean;
  isApproved: boolean;
  helpfulVotes: number;
  unhelpfulVotes: number;
  createdAt: string;
  updatedAt: string;
};

export default function Page() {
  const [ratings, setRatings] = useState<AdminRatingModel[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchInput, setSearchInput] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [selectedRating, setSelectedRating] = useState<AdminRatingModel | null>(
    null
  );
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [approvalStatus, setApprovalStatus] = useState<
    'all' | 'pending' | 'approved' | 'rejected'
  >('pending');

  // Pagination settings
  const pageSize = 10;

  const fetchRatings = async (
    page = 1,
    search = '',
    status = approvalStatus
  ) => {
    setLoading(true);
    try {
      // Using the admin ratings API endpoint
      const response = await clientSideFetch<AdminRatingModel[]>(
        `${ADMIN_API_PATHS.RATINGS}?page=${page}&pageSize=${pageSize}${search ? `&search=${search}` : ''}${status !== 'all' ? `&status=${status}` : ''}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );

      if (response.error) {
        toast.error('Failed to load ratings: ' + response.error.details);
        return;
      }

      setRatings(response.data || []);

      // If pagination data is available in the response
      if (response.pagination) {
        setTotalPages(Math.ceil(response.pagination.total / pageSize));
      }
    } catch (error) {
      console.error('Error fetching ratings:', error);
      toast.error('An error occurred while loading ratings');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRatings(currentPage, searchInput, approvalStatus);
  }, [currentPage, approvalStatus]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setCurrentPage(1); // Reset to first page on new search
    fetchRatings(1, searchInput, approvalStatus);
  };

  const handleDelete = (rating: AdminRatingModel) => {
    setSelectedRating(rating);
    setShowDeleteDialog(true);
  };

  const handleDeleteConfirm = async () => {
    if (!selectedRating) return;

    try {
      const response = await clientSideFetch<boolean>(
        ADMIN_API_PATHS.RATING_DETAIL.replace(':id', selectedRating.id),
        {
          method: 'DELETE',
        }
      );

      if (response.error) {
        toast.error('Failed to delete rating: ' + response.error.details);
        return;
      }

      toast.success('Rating deleted successfully');
      setRatings(ratings.filter((rating) => rating.id !== selectedRating.id));
      setShowDeleteDialog(false);
      setSelectedRating(null);
    } catch (error) {
      console.error('Error deleting rating:', error);
      toast.error('An error occurred while deleting the rating');
    }
  };

  const handleApprove = async (rating: AdminRatingModel) => {
    try {
      const response = await clientSideFetch<boolean>(
        ADMIN_API_PATHS.RATING_APPROVE.replace(':id', rating.id),
        {
          method: 'PUT',
        }
      );

      if (response.error) {
        toast.error('Failed to approve rating: ' + response.error.details);
        return;
      }

      toast.success('Rating approved successfully');
      // Update the rating in the list
      setRatings(
        ratings.map((r) =>
          r.id === rating.id ? { ...r, isApproved: true } : r
        )
      );
    } catch (error) {
      console.error('Error approving rating:', error);
      toast.error('An error occurred while approving the rating');
    }
  };

  const handleReject = async (rating: AdminRatingModel) => {
    try {
      const response = await clientSideFetch<boolean>(
        ADMIN_API_PATHS.RATING_REJECT.replace(':id', rating.id),
        {
          method: 'PUT',
        }
      );

      if (response.error) {
        toast.error('Failed to reject rating: ' + response.error.details);
        return;
      }

      toast.success('Rating rejected successfully');
      // Update the rating in the list
      setRatings(
        ratings.map((r) =>
          r.id === rating.id ? { ...r, isVisible: false } : r
        )
      );
    } catch (error) {
      console.error('Error rejecting rating:', error);
      toast.error('An error occurred while rejecting the rating');
    }
  };

  const handlePageChange = (newPage: number) => {
    if (newPage >= 1 && newPage <= totalPages) {
      setCurrentPage(newPage);
    }
  };

  const handleStatusChange = (
    status: 'all' | 'pending' | 'approved' | 'rejected'
  ) => {
    setApprovalStatus(status);
    setCurrentPage(1);
  };

  // Function to render star rating
  const renderStars = (rating: number) => {
    const stars = [];
    const fullStars = Math.floor(rating);
    const hasHalfStar = rating % 1 >= 0.5;

    for (let i = 0; i < fullStars; i++) {
      stars.push(
        <span key={`full-${i}`} className='text-yellow-400'>
          ★
        </span>
      );
    }

    if (hasHalfStar) {
      stars.push(
        <span key='half' className='text-yellow-400'>
          ★
        </span>
      );
    }

    const remainingStars = 5 - fullStars - (hasHalfStar ? 1 : 0);
    for (let i = 0; i < remainingStars; i++) {
      stars.push(
        <span key={`empty-${i}`} className='text-gray-300'>
          ★
        </span>
      );
    }

    return <div className='flex'>{stars}</div>;
  };

  if (loading && ratings.length === 0) {
    return <Loading />;
  }

  return (
    <div className='container mx-auto px-4'>
      <div className='flex justify-between items-center mb-6'>
        <h2 className='text-2xl font-semibold text-primary'>
          Ratings Management
        </h2>
      </div>

      {/* Filters */}
      <div className='mb-6'>
        <div className='flex mb-4'>
          <button
            className={`px-4 py-2 mr-2 rounded-md ${
              approvalStatus === 'all'
                ? 'bg-blue-500 text-white'
                : 'bg-gray-200'
            }`}
            onClick={() => handleStatusChange('all')}
          >
            All
          </button>
          <button
            className={`px-4 py-2 mr-2 rounded-md ${
              approvalStatus === 'pending'
                ? 'bg-blue-500 text-white'
                : 'bg-gray-200'
            }`}
            onClick={() => handleStatusChange('pending')}
          >
            Pending
          </button>
          <button
            className={`px-4 py-2 mr-2 rounded-md ${
              approvalStatus === 'approved'
                ? 'bg-blue-500 text-white'
                : 'bg-gray-200'
            }`}
            onClick={() => handleStatusChange('approved')}
          >
            Approved
          </button>
          <button
            className={`px-4 py-2 mr-2 rounded-md ${
              approvalStatus === 'rejected'
                ? 'bg-blue-500 text-white'
                : 'bg-gray-200'
            }`}
            onClick={() => handleStatusChange('rejected')}
          >
            Rejected
          </button>
        </div>

        <form onSubmit={handleSearch} className='flex gap-2'>
          <div className='relative flex-grow'>
            <div className='absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none'>
              <MagnifyingGlassIcon className='w-5 h-5 text-gray-500' />
            </div>
            <input
              type='text'
              className='bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full pl-10 p-2.5'
              placeholder='Search ratings by product name, user, or review content...'
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
            />
          </div>
          <Button
            type='submit'
            className='px-4 py-2.5 bg-blue-500 text-white rounded-lg hover:bg-blue-600'
          >
            Search
          </Button>
        </form>
      </div>

      {/* Ratings Table */}
      <div className='relative overflow-x-auto shadow-md sm:rounded-lg'>
        <table className='w-full text-sm text-left text-gray-500'>
          <thead className='text-xs text-gray-700 uppercase bg-gray-50'>
            <tr>
              <th scope='col' className='px-6 py-3'>
                Product
              </th>
              <th scope='col' className='px-6 py-3'>
                User
              </th>
              <th scope='col' className='px-6 py-3'>
                Rating
              </th>
              <th scope='col' className='px-6 py-3'>
                Review Title
              </th>
              <th scope='col' className='px-6 py-3'>
                Content
              </th>
              <th scope='col' className='px-6 py-3'>
                Status
              </th>
              <th scope='col' className='px-6 py-3'>
                Date
              </th>
              <th scope='col' className='px-6 py-3'>
                Actions
              </th>
            </tr>
          </thead>
          <tbody>
            {ratings.length > 0 ? (
              ratings.map((rating) => (
                <tr
                  key={rating.id}
                  className='bg-white border-b hover:bg-gray-50'
                >
                  <td className='px-6 py-4'>{rating.productName}</td>
                  <td className='px-6 py-4'>
                    {rating.userName ||
                      `User: ${rating.userId.substring(0, 8)}...`}
                  </td>
                  <td className='px-6 py-4'>{renderStars(rating.rating)}</td>
                  <td className='px-6 py-4'>
                    {rating.reviewTitle || 'No title'}
                  </td>
                  <td className='px-6 py-4 max-w-xs truncate'>
                    {rating.reviewContent || 'No content'}
                  </td>
                  <td className='px-6 py-4'>
                    {rating.isApproved ? (
                      <span className='px-2 py-1 bg-green-100 text-green-800 rounded-full text-xs'>
                        Approved
                      </span>
                    ) : !rating.isVisible ? (
                      <span className='px-2 py-1 bg-red-100 text-red-800 rounded-full text-xs'>
                        Rejected
                      </span>
                    ) : (
                      <span className='px-2 py-1 bg-yellow-100 text-yellow-800 rounded-full text-xs'>
                        Pending
                      </span>
                    )}
                  </td>
                  <td className='px-6 py-4'>
                    {dayjs(rating.createdAt).format('MMM D, YYYY')}
                  </td>
                  <td className='px-6 py-4 whitespace-nowrap'>
                    <div className='flex space-x-2'>
                      {!rating.isApproved && rating.isVisible && (
                        <button
                          onClick={() => handleApprove(rating)}
                          className='text-blue-500 hover:text-blue-700'
                          title='Approve'
                        >
                          <CheckCircleIcon className='h-5 w-5' />
                        </button>
                      )}
                      {!rating.isApproved && rating.isVisible && (
                        <button
                          onClick={() => handleReject(rating)}
                          className='text-red-500 hover:text-red-700'
                          title='Reject'
                        >
                          <XCircleIcon className='h-5 w-5' />
                        </button>
                      )}
                      <button
                        onClick={() => handleDelete(rating)}
                        className='text-red-500 hover:text-red-700'
                        title='Delete'
                      >
                        <TrashIcon className='h-5 w-5' />
                      </button>
                    </div>
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={8} className='px-6 py-4 text-center'>
                  No ratings found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className='flex justify-center mt-6'>
          <nav>
            <ul className='flex space-x-2'>
              <li>
                <button
                  onClick={() => handlePageChange(currentPage - 1)}
                  disabled={currentPage === 1}
                  className={`px-3 py-1 rounded-md ${
                    currentPage === 1
                      ? 'bg-gray-200 text-gray-500 cursor-not-allowed'
                      : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                  }`}
                >
                  Previous
                </button>
              </li>
              {[...Array(totalPages).keys()].map((page) => (
                <li key={page + 1}>
                  <button
                    onClick={() => handlePageChange(page + 1)}
                    className={`px-3 py-1 rounded-md ${
                      currentPage === page + 1
                        ? 'bg-blue-500 text-white'
                        : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                    }`}
                  >
                    {page + 1}
                  </button>
                </li>
              ))}
              <li>
                <button
                  onClick={() => handlePageChange(currentPage + 1)}
                  disabled={currentPage === totalPages}
                  className={`px-3 py-1 rounded-md ${
                    currentPage === totalPages
                      ? 'bg-gray-200 text-gray-500 cursor-not-allowed'
                      : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                  }`}
                >
                  Next
                </button>
              </li>
            </ul>
          </nav>
        </div>
      )}

      {/* Delete Confirmation Dialog */}
      <ConfirmDialog
        open={showDeleteDialog}
        title='Delete Rating'
        message='Are you sure you want to delete this rating? This action cannot be undone.'
        onConfirm={handleDeleteConfirm}
        onClose={() => {
          setShowDeleteDialog(false);
          setSelectedRating(null);
        }}
      />
    </div>
  );
}
