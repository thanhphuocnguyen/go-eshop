import React from 'react';

// Define types for variant images with variant assignments
export type VariantImage = {
  file: File;
  preview: string;
  role?: string | null; // Role of the image (e.g., 'thumbnail', 'gallery')
  variantIds: string[]; // IDs of variants this image is assigned to
};

const ProductDetailFormContext = React.createContext<{
  tempProductImages: VariantImage[];
  setTempProductImages: React.Dispatch<React.SetStateAction<VariantImage[]>>;
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

  return (
    <ProductDetailFormContext.Provider
      value={{
        tempProductImages,
        setTempProductImages,
      }}
    >
      {children}
    </ProductDetailFormContext.Provider>
  );
};
