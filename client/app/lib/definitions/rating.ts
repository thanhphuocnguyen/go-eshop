export type RatingModel = {
  id: string;
  user_id: string;
  name: string;
  rating: number;
  review_title: string;
  review_content: string;
  verified_purchase: boolean;
  helpful_votes: number;
  unhelpful_votes: number;
  created_at?: string;
  images: {
    id: string;
    url: string;
  }[];
};
