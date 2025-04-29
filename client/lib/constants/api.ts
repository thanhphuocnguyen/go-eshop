const _apiPaths = {
  LOGIN: '/auth/login',
  REGISTER: '/auth/register',
  REFRESH_TOKEN: '/auth/refresh-token',
  FORGOT_PASSWORD: '/auth/forgot-password',
  RESET_PASSWORD: '/auth/reset-password',
  CATEGORIES: '/categories',
  CATEGORY: '/categories/:id',
  CATEGORY_PRODUCTS: '/categories/:id/products',
  PRODUCTS: '/products',
  PRODUCT_DETAIL: '/products/:id',
  PRODUCT_VARIANTS: '/products/:id/variants',
  LOGOUT: '/auth/logout',
  USER: '/user',
  USER_ADDRESSES: '/user/addresses',
  USER_ADDRESS: '/user/addresses/:id',
  USER_ADDRESS_DEFAULT: '/user/addresses/:id/default',
  CART: '/cart',
  CART_ITEM: '/cart/item',
  CART_ITEM_DETAIL: '/cart/item/:id',
  CART_ITEM_QUANTITY: '/cart/item/:id/quantity',
  CART_CLEAR: '/cart/clear',
  CHECKOUT: '/cart/checkout',
  ORDER: '/order',
  ORDER_ITEM: '/order/:id',
  PAYMENTS: '/payments',
  PAYMENT_DETAIL: '/payments/:id',
  ATTRIBUTES: '/attributes',
  ATTRIBUTE: '/attributes/:id',
  BRANDS: '/brands',
  BRAND: '/brands/:id',
  BRAND_PRODUCTS: '/brands/:id/products',
  UPDATE_BRANDS: '/brands/:id',
  COLLECTIONS: '/collections',
  COLLECTION: '/collections/:id',
  COLLECTION_PRODUCTS: '/collections/:id/products',
  PRODUCT_IMAGES_UPLOAD: '/images/product/:id',
  PRODUCT_VARIANT_IMAGE_UPLOAD: '/images/product-variant/:id',
} as const;

const attachBasePath = {
  get(target: typeof _apiPaths, prop: keyof typeof target) {
    const BasePath = process.env.NEXT_PUBLIC_API_BASE_URL;
    const path = `${BasePath}${target[prop]}`;
    return path;
  },
};

export const API_PATHS = new Proxy(_apiPaths, attachBasePath);
