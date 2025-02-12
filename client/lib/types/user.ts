export type UserModel = {
  email: string;
  fullname: string;
  username: string;
  created_at: Date;
  verified_email: boolean;
  verified_phone: boolean;
  role: string;
  updated_at: Date;
  password_changed_at: Date;
  addresses: AddressModel[];
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
