import { AttributeDetailModel } from '@/lib/definitions';
import React, { useState } from 'react';

// Define types for variant images with variant assignments
export type VariantImage = {
  file: File;
  preview: string;
  variantIds: string[]; // IDs of variants this image is assigned to
};

const ProductDetailFormContext = React.createContext<{
  tempProductImages: VariantImage[];
  setTempProductImages: React.Dispatch<React.SetStateAction<VariantImage[]>>;
  selectedAttributes: AttributeDetailModel[];
  setSelectedAttributes: React.Dispatch<
    React.SetStateAction<AttributeDetailModel[]>
  >;
} | null>(null);

export const useProductDetailFormContext = () => {
  const context = React.useContext(ProductDetailFormContext);
  if (!context) {
    throw new Error(
      'useProductDetailFormContext must be used within a ProductDetailFormProvider'
    );
  }
  return context;
};

export const ProductDetailFormProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [tempProductImages, setTempProductImages] = React.useState<
    VariantImage[]
  >([]);
  const [selectedAttributes, setSelectedAttributes] = useState<
    AttributeDetailModel[]
  >([]);

  return (
    <ProductDetailFormContext.Provider
      value={{
        tempProductImages,
        setTempProductImages,
        selectedAttributes,
        setSelectedAttributes,
      }}
    >
      {children}
    </ProductDetailFormContext.Provider>
  );
};
