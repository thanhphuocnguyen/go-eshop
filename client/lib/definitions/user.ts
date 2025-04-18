export type UserModel = {
  id: string;
  role: string;
  username: string;
  fullname: string;
  email?: string;
  created_at?: Date;
  verified_email?: boolean;
  verified_phone?: boolean;
  updated_at?: Date;
  password_changed_at?: Date;
  addresses?: AddressModel[];
};

export type AddressModel = {
  id: number;
  user_id: number;
  address: string;
  city: string;
  state: string;
  country: string;
  postal_code: string;
  created_at: Date;
  updated_at: Date;
};
