'use client';

import { useState, Fragment } from 'react';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { StarIcon } from '@heroicons/react/24/solid';
import { StarIcon as StarOutlineIcon } from '@heroicons/react/24/outline';
import {
  Dialog,
  DialogPanel,
  DialogTitle,
  Transition,
  TransitionChild,
} from '@headlessui/react';
import { XCircleIcon } from '@heroicons/react/24/outline';
import { useForm } from 'react-hook-form';
import { OrderItemModel } from '@/app/lib/definitions';
import { ImageUploader } from '@/app/components/FormFields';
import { toast } from 'react-toastify';

interface OrderItemRatingProps {
  orderItemId: string;
  ratingModel?: OrderItemModel['rating'];
}

type ReviewFormData = {
  headline: string;
  comment: string;
  imageUrls: string[];
};

export default function OrderItemRating({
  orderItemId,
  ratingModel,
}: OrderItemRatingProps) {
  // Form handling with react-hook-form
  const {
    register,
    handleSubmit: handleFormSubmit,
    formState: { errors },
    setValue,
  } = useForm<ReviewFormData>({
    defaultValues: {
      headline: ratingModel?.title ?? '',
      comment: ratingModel?.content ?? '',
      imageUrls: ratingModel?.imageUrl ? [ratingModel.imageUrl] : [],
    },
  });

  const [rating, setRating] = useState<number | null>(
    ratingModel?.rating ?? null
  );
  const [existingRating, setExistingRating] = useState<number | null>(
    ratingModel?.rating ?? null
  );
  const [hoveredRating, setHoveredRating] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isOpen, setIsOpen] = useState(false);
  const [imageFiles, setImageFiles] = useState<(File & { preview: string })[]>(
    []
  );
  const [uploadingImages, setUploadingImages] = useState(false);

  const closeModal = () => {
    setIsOpen(false);
  };

  const openModal = () => {
    setIsOpen(true);
  };

  // Handle image upload from the ImageUploader component
  const handleImageUpload = (files: (File & { preview: string })[]) => {
    setImageFiles([...imageFiles, ...files]);

    // You can also update form values with the preview URLs
    const currentUrls = imageFiles.map((file) => file.preview);
    const newUrls = files.map((file) => file.preview);
    setValue('imageUrls', [...currentUrls, ...newUrls]);
  };

  // Remove an image by index
  const handleRemoveImage = (indexToRemove: number) => {
    setImageFiles(imageFiles.filter((_, index) => index !== indexToRemove));
    setValue(
      'imageUrls',
      imageFiles
        .filter((_, index) => index !== indexToRemove)
        .map((file) => file.preview)
    );
  };

  const onSubmit = async (data: ReviewFormData) => {
    if (!rating) {
      setError('Please select a star rating');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      // Create a FormData object for the multipart/form-data request
      const formData = new FormData();

      // Add all image files to the FormData object
      if (imageFiles.length > 0) {
        setUploadingImages(true);
        imageFiles.forEach((file) => {
          formData.append('files', file);
        });
      }

      // Add other form fields
      formData.append('orderItemId', orderItemId);
      formData.append('rating', rating.toString());
      formData.append('title', data.headline.trim());
      formData.append('content', data.comment.trim());

      // Submit the review with image URL if available
      const response = await clientSideFetch(PUBLIC_API_PATHS.RATING, {
        method: 'POST',
        body: formData,
      });

      if (response.error) {
        setError(response.error.details || 'Failed to submit rating');
      } else {
        setExistingRating(rating);
        toast.success('Your review has been submitted successfully!');
        closeModal();
      }
    } catch (err) {
      setError('An error occurred while submitting your rating');
    } finally {
      setLoading(false);
      setUploadingImages(false);
    }
  };

  // If there's already a rating, just show the rating summary
  // Otherwise, provide a button to create a new rating
  if (existingRating) {
    return (
      <div className='p-2 bg-gray-50 rounded border border-gray-200 inline-flex items-center text-sm'>
        <div className='flex items-center'>
          <span className='text-gray-700 mr-1'>Your Rating:</span>
          <div className='flex'>
            {[1, 2, 3, 4, 5].map((star) => (
              <StarIcon
                key={star}
                className={`w-4 h-4 ${
                  star <= existingRating ? 'text-yellow-400' : 'text-gray-300'
                }`}
              />
            ))}
          </div>
          <span className='ml-1 text-xs font-medium text-gray-600'>
            ({existingRating.toFixed(1)})
          </span>
        </div>
      </div>
    );
  }

  // Button to open the rating dialog (only for creating new ratings)
  return (
    <>
      <button
        onClick={openModal}
        className='px-4 py-2 rounded-md text-sm font-medium transition-colors bg-indigo-600 text-white hover:bg-indigo-700'
      >
        Rate This Product
      </button>

      <Transition appear show={isOpen} as={Fragment}>
        <Dialog as='div' className='relative z-50' onClose={closeModal}>
          <TransitionChild
            as={Fragment}
            enter='ease-out duration-300'
            enterFrom='opacity-0'
            enterTo='opacity-100'
            leave='ease-in duration-200'
            leaveFrom='opacity-100'
            leaveTo='opacity-0'
          >
            <div className='fixed inset-0 bg-black bg-opacity-25' />
          </TransitionChild>

          <div className='fixed inset-0 overflow-y-auto'>
            <div className='flex min-h-full items-center justify-center p-4 text-center'>
              <TransitionChild
                as={Fragment}
                enter='ease-out duration-300'
                enterFrom='opacity-0 scale-95'
                enterTo='opacity-100 scale-100'
                leave='ease-in duration-200'
                leaveFrom='opacity-100 scale-100'
                leaveTo='opacity-0 scale-95'
              >
                <DialogPanel className='w-full max-w-2xl transform overflow-hidden rounded-2xl bg-white p-8 text-left align-middle shadow-xl transition-all'>
                  {/* Header with product rating title */}
                  <div className='flex justify-between items-center mb-6 border-b pb-4'>
                    <DialogTitle
                      as='h3'
                      className='text-xl font-semibold text-gray-900'
                    >
                      Share Your Experience
                    </DialogTitle>
                    <button
                      type='button'
                      className='text-gray-400 hover:text-gray-500 transition-colors'
                      onClick={closeModal}
                    >
                      <XCircleIcon className='h-7 w-7' aria-hidden='true' />
                    </button>
                  </div>

                  <form onSubmit={handleFormSubmit(onSubmit)}>
                    <div className='mt-2'>
                      {/* Rating section with bigger stars */}
                      <div className='bg-gray-50 p-4 rounded-lg mb-6'>
                        <h4 className='text-base font-medium text-gray-800 mb-3'>
                          How would you rate this product?
                        </h4>
                        <div className='flex items-center mb-2'>
                          <div className='flex'>
                            {[1, 2, 3, 4, 5].map((star) => {
                              const isActive =
                                (hoveredRating || rating || 0) >= star;

                              return (
                                <div
                                  key={star}
                                  className='cursor-pointer p-1'
                                  onMouseEnter={() => setHoveredRating(star)}
                                  onMouseLeave={() => setHoveredRating(null)}
                                  onClick={() => setRating(star)}
                                >
                                  {isActive ? (
                                    <StarIcon className='w-8 h-8 text-yellow-400' />
                                  ) : (
                                    <StarOutlineIcon className='w-8 h-8 text-gray-400' />
                                  )}
                                </div>
                              );
                            })}
                          </div>
                          <span className='ml-4 text-sm font-medium text-gray-700'>
                            {rating ? (
                              <span className='px-3 py-1 bg-yellow-100 border border-orange-900 shadow-sm text-orange-600 rounded-full'>
                                {rating} star{rating !== 1 ? 's' : ''}
                              </span>
                            ) : (
                              'Select a rating'
                            )}
                          </span>
                        </div>
                        <p className='text-xs text-gray-500 mt-1'>
                          Click on a star to set your rating. 5 stars being the
                          best experience.
                        </p>
                      </div>

                      {/* Headline field */}
                      <div className='mb-5'>
                        <label
                          htmlFor='headline'
                          className='block text-sm font-medium text-gray-700 mb-2'
                        >
                          Review Title <span className='text-red-500'>*</span>
                        </label>
                        <input
                          type='text'
                          id='headline'
                          {...register('headline', {
                            required: 'Review title is required',
                          })}
                          placeholder='Summarize your experience in a few words'
                          className={`w-full px-4 py-2.5 text-sm border ${
                            errors.headline
                              ? 'border-red-500'
                              : 'border-gray-300'
                          } rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500`}
                        />
                        {errors.headline && (
                          <p className='mt-1 text-xs text-red-600'>
                            {errors.headline.message}
                          </p>
                        )}
                      </div>

                      {/* Review comment field */}
                      <div className='mb-5'>
                        <label
                          htmlFor='comment'
                          className='block text-sm font-medium text-gray-700 mb-2'
                        >
                          Your Review <span className='text-red-500'>*</span>
                        </label>
                        <textarea
                          id='comment'
                          {...register('comment', {
                            required: 'Review content is required',
                          })}
                          placeholder='Share your thoughts about this product. What did you like or dislike?'
                          className={`w-full px-4 py-3 text-sm border ${
                            errors.comment
                              ? 'border-red-500'
                              : 'border-gray-300'
                          } rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500`}
                          rows={5}
                        />
                        {errors.comment && (
                          <p className='mt-1 text-xs text-red-600'>
                            {errors.comment.message}
                          </p>
                        )}
                        <p className='mt-1 text-xs text-gray-500'>
                          Your review helps other shoppers make better
                          decisions.
                        </p>
                      </div>

                      {/* Image upload section using the reusable ImageUploader component */}
                      <div className='mb-6'>
                        <ImageUploader
                          label='Add Photos'
                          multiple={true}
                          onUpload={handleImageUpload}
                        />

                        {/* Preview of uploaded images */}
                        {imageFiles.length > 0 && (
                          <div className='mt-4'>
                            <h4 className='text-sm font-medium text-gray-700 mb-2'>
                              Uploaded Images
                            </h4>
                            <div className='flex flex-wrap gap-4'>
                              {imageFiles.map((file, index) => (
                                <div key={index} className='relative'>
                                  <img
                                    src={file.preview}
                                    alt={`Preview ${index + 1}`}
                                    className='w-32 h-32 object-cover border border-gray-300 rounded-md'
                                  />
                                  <button
                                    type='button'
                                    onClick={() => handleRemoveImage(index)}
                                    className='absolute top-1 right-1 bg-white rounded-full p-1 shadow-md hover:bg-gray-100'
                                  >
                                    <XCircleIcon className='h-5 w-5 text-gray-600' />
                                  </button>
                                </div>
                              ))}
                            </div>
                          </div>
                        )}

                        <p className='mt-2 text-xs text-gray-500'>
                          Add photos to help other shoppers visualize your
                          experience. Max size: 5MB per image.
                        </p>
                      </div>

                      {error && (
                        <div className='mb-5 p-4 bg-red-50 border border-red-200 rounded-md'>
                          <p className='text-sm text-red-600 font-medium'>
                            {error}
                          </p>
                        </div>
                      )}

                      <p className='text-xs text-gray-500 mb-5'>
                        By submitting, you agree to our Review Guidelines.
                        Reviews are moderated and appear if they comply with our
                        guidelines.
                      </p>
                    </div>

                    <div className='mt-6 flex justify-end gap-3 border-t pt-4'>
                      <button
                        type='button'
                        className='px-5 py-2.5 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-gray-500'
                        onClick={closeModal}
                      >
                        Cancel
                      </button>
                      <button
                        type='submit'
                        disabled={!rating || loading}
                        className={`px-5 py-2.5 text-sm font-medium text-white rounded-md focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-indigo-500 ${
                          !rating || loading
                            ? 'bg-gray-300 cursor-not-allowed'
                            : 'bg-indigo-600 hover:bg-indigo-700'
                        }`}
                      >
                        {loading || uploadingImages
                          ? 'Submitting...'
                          : 'Submit Review'}
                      </button>
                    </div>
                  </form>
                </DialogPanel>
              </TransitionChild>
            </div>
          </div>
        </Dialog>
      </Transition>
    </>
  );
}
