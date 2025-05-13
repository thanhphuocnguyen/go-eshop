'use client';

import { useState, useEffect, Fragment } from 'react';
import { apiFetch } from '@/lib/apis/api';
import { PUBLIC_API_PATHS } from '@/lib/constants/api';
import { StarIcon } from '@heroicons/react/24/solid';
import { StarIcon as StarOutlineIcon } from '@heroicons/react/24/outline';
import {
  Dialog,
  DialogPanel,
  DialogTitle,
  Transition,
  TransitionChild,
} from '@headlessui/react';
import { XCircleIcon, PhotoIcon } from '@heroicons/react/24/outline';
import { useForm } from 'react-hook-form';

interface OrderItemRatingProps {
  orderId: string;
  productId: string;
}

type ReviewFormData = {
  headline: string;
  comment: string;
  imageUrl: string | null;
};

export default function OrderItemRating({
  orderId,
  productId,
}: OrderItemRatingProps) {
  // Form handling with react-hook-form
  const {
    register,
    handleSubmit: handleFormSubmit,
    formState: { errors },
    setValue,
    watch,
  } = useForm<ReviewFormData>({
    defaultValues: {
      headline: '',
      comment: '',
      imageUrl: null,
    },
  });

  const [rating, setRating] = useState<number | null>(null);
  const [existingRating, setExistingRating] = useState<number | null>(null);
  const [hoveredRating, setHoveredRating] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isOpen, setIsOpen] = useState(false);
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [uploadingImage, setUploadingImage] = useState(false);

  // Check if user has already rated this product in this order
  useEffect(() => {
    const fetchExistingRating = async () => {
      try {
        const response = await apiFetch(
          `${PUBLIC_API_PATHS.ORDER_ITEM.replace(':id', orderId)}/ratings/${productId}`,
          {}
        );

        if (response.data) {
          setExistingRating(response.data.rating);
          setRating(response.data.rating);
          setValue('headline', response.data.headline || '');
          setValue('comment', response.data.comment || '');

          if (response.data.imageUrl) {
            setValue('imageUrl', response.data.imageUrl);
            setPreviewUrl(response.data.imageUrl);
          }

          setSubmitted(true);
        }
      } catch (err) {
        // No existing rating found, that's okay
      }
    };

    fetchExistingRating();
  }, [orderId, productId, setValue]);

  const closeModal = () => {
    setIsOpen(false);
  };

  const openModal = () => {
    setIsOpen(true);
  };

  // Handle image upload
  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Check if file is an image
    if (!file.type.match('image.*')) {
      setError('Please select an image file (png, jpg, jpeg)');
      return;
    }

    // Check file size (limit to 5MB)
    if (file.size > 5 * 1024 * 1024) {
      setError('Image size should not exceed 5MB');
      return;
    }

    setImageFile(file);

    // Create preview URL
    const reader = new FileReader();
    reader.onloadend = () => {
      setPreviewUrl(reader.result as string);
    };
    reader.readAsDataURL(file);
  };

  // Remove uploaded image
  const removeImage = () => {
    setImageFile(null);
    setPreviewUrl(null);
    setValue('imageUrl', null);
  };

  const onSubmit = async (formData: ReviewFormData) => {
    if (!rating) {
      setError('Please select a star rating');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      let imageUrl = watch('imageUrl');

      // Upload the image first if there's a new file
      if (imageFile) {
        setUploadingImage(true);

        // Create FormData to upload the image
        const formData = new FormData();
        formData.append('image', imageFile);

        // Upload image to your server endpoint
        try {
          const uploadResponse = await apiFetch(
            `${PUBLIC_API_PATHS.UPLOAD_IMAGE}`,
            {
              method: 'POST',
              body: formData,
              isFormData: true,
            }
          );

          if (uploadResponse.error) {
            throw new Error(
              uploadResponse.error.details || 'Failed to upload image'
            );
          }

          imageUrl = uploadResponse.data.url;
          setValue('imageUrl', imageUrl);
        } catch (err) {
          setError('Failed to upload image. Please try again.');
          setLoading(false);
          setUploadingImage(false);
          return;
        } finally {
          setUploadingImage(false);
        }
      }

      // Submit the review with image URL if available
      const response = await apiFetch(
        `${PUBLIC_API_PATHS.ORDER_ITEM.replace(':id', orderId)}/rate`,
        {
          method: 'POST',
          body: {
            product_id: productId,
            rating,
            headline: formData.headline.trim(),
            comment: formData.comment.trim(),
            imageUrl: imageUrl || undefined,
          },
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
          <div className='flex items-center'>
            <span className='mr-2'>Your Rating: </span>
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
          </div>
        ) : (
          'Rate This Product'
        )}
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
                      {existingRating
                        ? 'Update Your Product Review'
                        : 'Share Your Experience'}
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
                              <span className='px-3 py-1 bg-yellow-100 text-yellow-800 rounded-full'>
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

                      {/* Image upload section */}
                      <div className='mb-6'>
                        <label
                          htmlFor='image'
                          className='block text-sm font-medium text-gray-700 mb-2'
                        >
                          Add Photo{' '}
                          <span className='text-gray-500'>(Optional)</span>
                        </label>

                        {previewUrl ? (
                          <div className='relative rounded-md overflow-hidden mb-3'>
                            <img
                              src={previewUrl}
                              alt='Preview'
                              className='w-32 h-32 object-cover border border-gray-300 rounded-md'
                            />
                            <button
                              type='button'
                              onClick={removeImage}
                              className='absolute top-1 right-1 bg-white rounded-full p-1 shadow-md hover:bg-gray-100'
                            >
                              <XCircleIcon className='h-5 w-5 text-gray-600' />
                            </button>
                          </div>
                        ) : (
                          <div className='mb-3'>
                            <label
                              htmlFor='image-upload'
                              className='flex flex-col items-center justify-center w-32 h-32 border-2 border-gray-300 border-dashed rounded-md cursor-pointer bg-gray-50 hover:bg-gray-100 transition-colors'
                            >
                              <div className='flex flex-col items-center justify-center'>
                                <PhotoIcon className='h-10 w-10 text-gray-400 mb-1' />
                                <p className='text-xs text-gray-500'>
                                  Click to upload
                                </p>
                              </div>
                              <input
                                id='image-upload'
                                type='file'
                                accept='image/*'
                                onChange={handleImageChange}
                                className='sr-only'
                              />
                            </label>
                          </div>
                        )}
                        <p className='text-xs text-gray-500'>
                          Add a photo to help other shoppers visualize your
                          experience. Max size: 5MB.
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
                        {loading || uploadingImage
                          ? 'Submitting...'
                          : existingRating
                            ? 'Update Review'
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
