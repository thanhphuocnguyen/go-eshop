export type UserModel = {
  id: string;
  role: string;
  username: string;
  fullname: string;
  email: string;
  phone: string;
  createdAt?: Date;
  verifiedEmail?: boolean;
  verifiedPhone?: boolean;
  updatedAt?: Date;
  passwordChangedAt?: Date;
  addresses?: AddressModel[];
};

export type AddressModel = {
  id: string;
  phone: string;
  street: string;
  ward?: string;
  district: string;
  city: string;
  default: boolean;
};
