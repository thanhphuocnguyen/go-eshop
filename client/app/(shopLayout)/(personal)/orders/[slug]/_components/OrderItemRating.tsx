'use client';

import { useState, useEffect, Fragment } from 'react';
import { apiFetch } from '@/lib/apis/api';
import { PUBLIC_API_PATHS } from '@/lib/constants/api';
import { StarIcon } from '@heroicons/react/24/solid';
import { StarIcon as StarOutlineIcon } from '@heroicons/react/24/outline';
import { Dialog, Transition } from '@headlessui/react';
import { XCircleIcon } from '@heroicons/react/24/outline';

interface OrderItemRatingProps {
  orderId: string;
  productId: string;
}

export default function OrderItemRating({ orderId, productId }: OrderItemRatingProps) {
  const [rating, setRating] = useState<number | null>(null);
  const [existingRating, setExistingRating] = useState<number | null>(null);
  const [hoveredRating, setHoveredRating] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const [comment, setComment] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [isOpen, setIsOpen] = useState(false);

  // Check if user has already rated this product in this order
  useEffect(() => {
    const fetchExistingRating = async () => {
      try {
        const response = await apiFetch(`${PUBLIC_API_PATHS.ORDER_ITEM.replace(':id', orderId)}/ratings/${productId}`, {});
        
        if (response.data) {
          setExistingRating(response.data.rating);
          setRating(response.data.rating);
          setComment(response.data.comment || '');
          setSubmitted(true);
        }
      } catch (err) {
        // No existing rating found, that's okay
      }
    };

    fetchExistingRating();
  }, [orderId, productId]);

  const closeModal = () => {
    setIsOpen(false);
  };

  const openModal = () => {
    setIsOpen(true);
  };

  const handleSubmit = async () => {
    if (!rating) return;
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await apiFetch(
        `${PUBLIC_API_PATHS.ORDER_ITEM.replace(':id', orderId)}/rate`,
        {
          method: 'POST',
          body: {
            product_id: productId,
            rating,
            comment: comment.trim() || undefined
          }
        }
      );
      
      if (response.error) {
        setError(response.error.details || 'Failed to submit rating');
      } else {
        setSubmitted(true);
        setExistingRating(rating);
        closeModal();
      }
    } catch (err) {
      setError('An error occurred while submitting your rating');
    } finally {
      setLoading(false);
    }
  };

  const renderStars = (currentRating: number | null, interactive: boolean = false) => {
    return (
      <div className="flex">
        {[1, 2, 3, 4, 5].map((star) => {
          const isActive = (hoveredRating || currentRating || 0) >= star;
          
          if (!interactive) {
            return (
              <StarIcon 
                key={star}
                className={`w-5 h-5 ${isActive ? 'text-yellow-400' : 'text-gray-300'}`}
              />
            );
          }
          
          return (
            <div
              key={star}
              className="cursor-pointer"
              onMouseEnter={() => setHoveredRating(star)}
              onMouseLeave={() => setHoveredRating(null)}
              onClick={() => setRating(star)}
            >
              {isActive ? (
                <StarIcon className="w-5 h-5 text-yellow-400" />
              ) : (
                <StarOutlineIcon className="w-5 h-5 text-gray-400" />
              )}
            </div>
          );
        })}
      </div>
    );
  };

  // Button to open the rating dialog
  return (
    <>
      <button 
        onClick={openModal}
        className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
          existingRating 
            ? 'bg-indigo-100 text-indigo-700 hover:bg-indigo-200'
            : 'bg-indigo-600 text-white hover:bg-indigo-700'
        }`}
      >
        {existingRating ? (
          <div className="flex items-center">
            <span className="mr-2">Your Rating: </span>
            {renderStars(existingRating)}
          </div>
        ) : (
          'Rate This Product'
        )}
      </button>

      <Transition appear show={isOpen} as={Fragment}>
        <Dialog as="div" className="relative z-50" onClose={closeModal}>
          <Transition.Child
            as={Fragment}
            enter="ease-out duration-300"
            enterFrom="opacity-0"
            enterTo="opacity-100"
            leave="ease-in duration-200"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <div className="fixed inset-0 bg-black bg-opacity-25" />
          </Transition.Child>

          <div className="fixed inset-0 overflow-y-auto">
            <div className="flex min-h-full items-center justify-center p-4 text-center">
              <Transition.Child
                as={Fragment}
                enter="ease-out duration-300"
                enterFrom="opacity-0 scale-95"
                enterTo="opacity-100 scale-100"
                leave="ease-in duration-200"
                leaveFrom="opacity-100 scale-100"
                leaveTo="opacity-0 scale-95"
              >
                <Dialog.Panel className="w-full max-w-md transform overflow-hidden rounded-2xl bg-white p-6 text-left align-middle shadow-xl transition-all">
                  <div className="flex justify-between items-center mb-4">
                    <Dialog.Title as="h3" className="text-lg font-medium leading-6 text-gray-900">
                      {existingRating ? 'Update Your Rating' : 'Rate This Product'}
                    </Dialog.Title>
                    <button
                      type="button"
                      className="text-gray-400 hover:text-gray-500"
                      onClick={closeModal}
                    >
                      <XCircleIcon className="h-6 w-6" aria-hidden="true" />
                    </button>
                  </div>

                  <div className="mt-2">
                    <div className="flex items-center mb-4">
                      <span className="text-sm font-medium text-gray-700 mr-3">Your Rating:</span>
                      <div className="flex">
                        {renderStars(rating, true)}
                      </div>
                      <span className="ml-2 text-sm text-gray-500">
                        {rating ? `${rating} star${rating !== 1 ? 's' : ''}` : 'Select rating'}
                      </span>
                    </div>
                    
                    <div className="mb-4">
                      <label htmlFor="comment" className="block text-sm font-medium text-gray-700 mb-2">
                        Review (Optional)
                      </label>
                      <textarea
                        id="comment"
                        value={comment}
                        onChange={(e) => setComment(e.target.value)}
                        placeholder="Share your thoughts about this product..."
                        className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                        rows={4}
                      />
                    </div>
                    
                    {error && (
                      <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
                        <p className="text-sm text-red-600">{error}</p>
                      </div>
                    )}
                  </div>

                  <div className="mt-6 flex justify-end gap-3">
                    <button
                      type="button"
                      className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-gray-500"
                      onClick={closeModal}
                    >
                      Cancel
                    </button>
                    <button
                      type="button"
                      disabled={!rating || loading}
                      className={`px-4 py-2 text-sm font-medium text-white rounded-md focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-indigo-500 ${
                        !rating || loading 
                          ? 'bg-gray-300 cursor-not-allowed' 
                          : 'bg-indigo-600 hover:bg-indigo-700'
                      }`}
                      onClick={handleSubmit}
                    >
                      {loading ? 'Submitting...' : 'Submit Rating'}
                    </button>
                  </div>
                </Dialog.Panel>
              </Transition.Child>
            </div>
          </div>
        </Dialog>
      </Transition>
    </>
  );
}