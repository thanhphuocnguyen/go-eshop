export type RatingModel = {
  id: string;
  userId: string;
  name: string;
  rating: number;
  reviewTitle: string;
  reviewContent: string;
  verifiedPurchase: boolean;
  helpfulVotes: number;
  unhelpfulVotes: number;
  createdAt?: string;
  images: {
    id: string;
    url: string;
  }[];
};
