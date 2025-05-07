const _apiPaths = {
  LOGIN: 'auth/login',
  REGISTER: 'auth/register',
  REFRESH_TOKEN: 'auth/refresh-token',
  FORGOT_PASSWORD: 'auth/forgot-password',
  RESET_PASSWORD: 'auth/reset-password',
  HOME_PAGE_DATA: 'homepage',
  CATEGORIES: 'categories',
  CATEGORY: 'categories/:slug',
  PRODUCTS: 'products',
  PRODUCT_DETAIL: 'products/:id',
  LOGOUT: 'auth/logout',
  USER: 'user',
  USER_ADDRESSES: 'user/addresses',
  USER_ADDRESS: 'user/addresses/:id',
  USER_ADDRESS_DEFAULT: 'user/addresses/:id/default',
  CART: 'cart',
  CART_ITEM: 'cart/item',
  CART_ITEM_QUANTITY: 'cart/item/:id/quantity',
  CART_CLEAR: 'cart/clear',
  CHECKOUT: 'cart/checkout',
  ORDERS: 'orders',
  ORDER_ITEM: 'orders/:id',
  PAYMENTS: 'payments',
  PAYMENT_DETAIL: 'payments/:id',
  BRANDS: 'brands',
  BRAND: 'brands/:slug',
  BRAND_PRODUCTS: 'brands/:slug/products',
  COLLECTIONS: 'collections',
  COLLECTION: 'collections/:slug',
  COLLECTION_PRODUCTS: 'collections/:slug/products',
  PRODUCT_IMAGES_UPLOAD: 'images/product/:id',
} as const;

const _adminPaths = {
  ATTRIBUTES: 'attributes',
  ATTRIBUTE: 'attributes/:id',
  BRANDS: 'brands',
  BRAND: 'brands/:id',
  BRAND_PRODUCTS: 'brands/:id/products',
  CATEGORIES: 'categories',
  PRODUCTS: 'products',
  PRODUCT_DETAIL: 'products/:id',
  PRODUCT_VARIANTS: 'products/:id/variants',
  CATEGORY: 'categories/:id',
  CATEGORY_PRODUCTS: 'categories/:id',
  COLLECTIONS: 'collections',
  COLLECTION: 'collections/:id',
  COLLECTION_PRODUCTS: 'collections/:id/products',
  PRODUCT_IMAGES_UPLOAD: 'images/products/:id',
  USERS: 'users',
  USER: 'users/:id',
  UPDATE_USER_ROLE: 'users/:id/role',
} as const;

const attachBasePath = {
  get(target: typeof _apiPaths, prop: keyof typeof target) {
    const BasePath = process.env.NEXT_PUBLIC_API_BASE_URL;
    const path = `${BasePath}/${target[prop]}`;
    return path;
  },
};

const attachBasePathAdmin = {
  get(target: typeof _adminPaths, prop: keyof typeof target) {
    const BasePath = process.env.NEXT_PUBLIC_API_BASE_URL;
    const path = `${BasePath}/admin/${target[prop]}`;
    return path;
  },
};

export const PUBLIC_API_PATHS = new Proxy(_apiPaths, attachBasePath);
export const ADMIN_API_PATHS = new Proxy(_adminPaths, attachBasePathAdmin);
