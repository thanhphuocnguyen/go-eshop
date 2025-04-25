export type UserModel = {
  id: string;
  role: string;
  username: string;
  fullname: string;
  email: string;
  phone: string;
  created_at?: Date;
  verified_email?: boolean;
  verified_phone?: boolean;
  updated_at?: Date;
  password_changed_at?: Date;
  addresses?: AddressModel[];
};

export type AddressModel = {
  id: number;
  phone: string;
  address: string;
  ward?: string;
  district: string;
  city: string;
  default: boolean;
};
