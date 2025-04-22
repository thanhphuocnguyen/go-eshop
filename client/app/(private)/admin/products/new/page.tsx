'use client';

import React from 'react';
import { ProductDetailForm } from '../_components/ProductDetailForm';
import { ProductDetailFormProvider } from '../_lib/contexts/ProductFormContext';

const Page: React.FC = () => {
  return (
    <ProductDetailFormProvider>
      <ProductDetailForm />
    </ProductDetailFormProvider>
  );
};

export default Page;
