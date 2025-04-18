import React from 'react';

export interface ProductImageFile {
  variantID: number | null;
  image: File | null;
}
const ProductDetailFormContext = React.createContext<{
  productImages: ProductImageFile[];
  setProductImages: React.Dispatch<React.SetStateAction<ProductImageFile[]>>;
  productVariantImages: ProductImageFile[];
  setProductVariantImages: React.Dispatch<
    React.SetStateAction<ProductImageFile[]>
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
  const [productImages, setProductImages] = React.useState<ProductImageFile[]>(
    []
  );
  const [productVariantImages, setProductVariantImages] = React.useState<
    ProductImageFile[]
  >([]);
  return (
    <ProductDetailFormContext.Provider
      value={{
        productVariantImages,
        setProductVariantImages,
        productImages,
        setProductImages,
      }}
    >
      {children}
    </ProductDetailFormContext.Provider>
  );
};
