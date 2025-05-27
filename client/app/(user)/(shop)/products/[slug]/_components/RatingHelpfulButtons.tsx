'use client';

import { clientSideFetch } from '@/app/lib/api/apiClient';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { useState } from 'react';
import { toast } from 'react-toastify';
import {
  HandThumbUpIcon as ThumbUpIcon,
  HandThumbDownIcon as ThumbDownIcon,
} from '@heroicons/react/24/outline';
import { useUser } from '@/app/hooks';

interface RatingHelpfulButtonsProps {
  ratingId: string;
  userId: string;
  helpfulVotes: number;
  unhelpfulVotes?: number;
}

export const RatingHelpfulButtons = ({
  ratingId,
  userId,
  helpfulVotes,
  unhelpfulVotes = 0,
}: RatingHelpfulButtonsProps) => {
  const { user } = useUser();
  const [helpful, setHelpful] = useState(helpfulVotes);
  const [unhelpful, setUnhelpful] = useState(unhelpfulVotes);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [voted, setVoted] = useState<'helpful' | 'unhelpful' | null>(null);

  // Check if this is the current user's review
  const isOwnReview = user && user.id === userId;

  // Don't show voting buttons for the user's own ratings
  if (isOwnReview) {
    return (
      <div className='flex items-start gap-x-4 text-sm'>
        <span className='text-gray-500'>Helpful ({helpful})</span>
        <span className='text-gray-500'>Not Helpful ({unhelpful})</span>
      </div>
    );
  }

  const handleVote = async (isHelpful: boolean) => {
    if (!user || isSubmitting) {
      toast.info('You need to be logged in to vote');
      return;
    }

    // Additional frontend check to prevent voting on own reviews
    if (isOwnReview) {
      toast.info("You can't vote on your own reviews");
      return;
    }

    setIsSubmitting(true);

    try {
      const response = await clientSideFetch(
        PUBLIC_API_PATHS.RATING_VOTE.replaceAll(':id', ratingId),
        {
          method: 'POST',
          body: {
            helpful: isHelpful,
          },
        }
      );

      if (response.error) {
        if (response.error.details?.toLowerCase().includes('forbidden')) {
          toast.info("You can't vote on your own reviews");
        } else {
          toast.error(
            'Failed to vote: ' +
              (response.error.details || 'Please try again later')
          );
        }
      } else {
        // Update local state
        if (isHelpful) {
          setHelpful((prev) => prev + 1);
          setVoted('helpful');
        } else {
          setUnhelpful((prev) => prev + 1);
          setVoted('unhelpful');
        }
        toast.success('Thanks for your feedback!');
      }
    } catch (err) {
      toast.error('An error occurred. Please try again later.');
      console.error('Vote error:', err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className='flex items-center gap-x-4 text-sm'>
      <button
        onClick={() => handleVote(true)}
        disabled={voted !== null || isSubmitting}
        className={`flex items-center gap-1 ${
          voted === 'helpful'
            ? 'text-green-600'
            : 'text-gray-500 hover:text-gray-700'
        } ${isSubmitting ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
        aria-label='Mark as helpful'
      >
        <ThumbUpIcon className='h-4 w-4' />
        <span>Helpful ({helpful})</span>
      </button>

      <button
        onClick={() => handleVote(false)}
        disabled={voted !== null || isSubmitting}
        className={`flex items-center gap-1 ${
          voted === 'unhelpful'
            ? 'text-red-600'
            : 'text-gray-500 hover:text-gray-700'
        } ${isSubmitting ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
        aria-label='Mark as not helpful'
      >
        <ThumbDownIcon className='h-4 w-4' />
        <span>Not Helpful ({unhelpful})</span>
      </button>
    </div>
  );
};

export default RatingHelpfulButtons;
