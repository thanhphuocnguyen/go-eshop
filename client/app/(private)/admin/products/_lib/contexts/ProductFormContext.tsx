import React from 'react';

// Define types for variant images with variant assignments
export type VariantImage = {
  id?: string; // For existing images from the backend
  file: File;
  preview: string;
  variantIds: string[]; // IDs of variants this image is assigned to
};

const ProductDetailFormContext = React.createContext<{
  productImages: File[];
  setProductImages: React.Dispatch<React.SetStateAction<File[]>>;
  variantImages: VariantImage[];
  setVariantImages: React.Dispatch<React.SetStateAction<VariantImage[]>>;
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
  const [productImages, setProductImages] = React.useState<File[]>([]);
  const [variantImages, setVariantImages] = React.useState<
    VariantImage[]
  >([]);
  
  return (
    <ProductDetailFormContext.Provider
      value={{
        variantImages,
        setVariantImages,
        productImages,
        setProductImages,
      }}
    >
      {children}
    </ProductDetailFormContext.Provider>
  );
};
