'use client';

import { useState, useEffect, useCallback } from 'react';
import { UseFormSetValue } from 'react-hook-form';
import { CheckoutFormValues } from '../../_lib/definitions';
import { UserModel } from '@/app/lib/definitions';

interface UseAddressLogicProps {
  user: UserModel | undefined;
  setValue: UseFormSetValue<CheckoutFormValues>;
  reset: (values?: Partial<CheckoutFormValues>) => void;
}

export const useAddressLogic = ({ user, setValue, reset }: UseAddressLogicProps) => {
  const [isNewAddress, setIsNewAddress] = useState(true);
  const [selectedAddressId, setSelectedAddressId] = useState<string | null>(null);

  const handleAddressChange = useCallback((e: React.ChangeEvent<HTMLSelectElement>) => {
    const addressId = e.target.value;

    if (addressId === '-1') {
      setIsNewAddress(true);
      setSelectedAddressId(null);

      setValue('address', {
        city: '',
        district: '',
        phone: '',
        street: '',
        ward: '',
      });
    } else {
      setIsNewAddress(false);
      setSelectedAddressId(addressId);

      const selectedAddress = user?.addresses?.find(
        (addr) => addr.id === addressId
      );
      if (selectedAddress) {
        setValue('addressId', selectedAddress.id);
        setValue('address', {
          street: selectedAddress.street,
          city: selectedAddress.city,
          district: selectedAddress.district,
          ward: selectedAddress.ward,
          phone: selectedAddress.phone,
        });
      }
    }
  }, [user?.addresses, setValue]);

  const handleAddNewAddress = useCallback(() => {
    setIsNewAddress(true);
    setSelectedAddressId(null);
    setValue('address', {
      city: '',
      district: '',
      phone: '',
      street: '',
      ward: '',
    });
  }, [setValue]);

  useEffect(() => {
    if (user) {
      const defaultValue: Partial<CheckoutFormValues> = {
        email: user.email,
        fullname: user.fullname,
        paymentMethod: 'cod',
        discountCode: '',
        termsAccepted: false,
      };

      if (user.addresses && user.addresses.length > 0) {
        const defaultAddress =
          user.addresses.find((addr) => addr.default) || user.addresses[0];
        defaultValue.addressId = defaultAddress.id;
        defaultValue.address = {
          street: defaultAddress.street,
          city: defaultAddress.city,
          district: defaultAddress.district,
          phone: defaultAddress.phone,
          ward: defaultAddress.ward,
        };
        setIsNewAddress(false);
        setSelectedAddressId(defaultAddress.id);
      } else {
        setIsNewAddress(true);
        setSelectedAddressId(null);
      }

      reset(defaultValue);
    }
  }, [user, reset, setValue]);

  return {
    isNewAddress,
    selectedAddressId,
    handleAddressChange,
    handleAddNewAddress,
  };
};
