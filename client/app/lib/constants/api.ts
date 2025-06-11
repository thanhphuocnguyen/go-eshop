const _apiPaths = {
  LOGIN: 'auth/login',
  REGISTER: 'auth/register',
  REFRESH_TOKEN: 'auth/refresh-token',
  FORGOT_PASSWORD: 'auth/forgot-password',
  RESET_PASSWORD: 'auth/reset-password',
  ADVERTISE_CATEGORIES: 'homepage',
  CATEGORIES: 'categories',
  CATEGORY: 'categories/:slug',
  COLLECTIONS: 'collections',
  COLLECTION: 'collections/:slug',
  BRANDS: 'brands',
  BRAND: 'brands/:slug',
  // products
  PRODUCTS: 'products',
  PRODUCT_DETAIL: 'products/:id',
  UPLOAD_PRODUCT_IMAGES: 'images/product/:id',
  PRODUCT_RATINGS: 'products/:id/ratings',
  // users
  GET_ME: 'user/me',
  USER_ADDRESSES: 'user/addresses',
  USER_ADDRESS: 'user/addresses/:id',
  USER_ADDRESS_DEFAULT: 'user/addresses/:id/default',
  // cart
  CART: 'cart',
  CART_ITEM: 'cart/item',
  CART_DISCOUNT_LIST: 'cart/discounts',
  CHECKOUT: 'cart/checkout',
  UPDATE_CART_ITEM_QUANTITY: 'cart/item/:id/quantity',
  // orders
  ORDERS: 'orders',
  ORDER: 'orders/:id',
  ORDER_RATINGS: 'orders/:id/ratings',
  CANCEL_ORDER: 'orders/:id/cancel',
  REFUND_ORDER: 'orders/:id/refund',
  CONFIRM_RECEIVED_ORDER: 'orders/:id/confirm-received',
  // payments
  PAYMENTS: 'payments',
  PAYMENT: 'payments/:id',
  // ratings
  RATING: 'ratings',
  RATING_VOTE: 'ratings/:id/helpful',
} as const;

const _adminPaths = {
  UPLOADS: 'uploads',
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
  PRODUCT_IMAGES_UPLOAD: 'images/products/:id',
  USERS: 'users',
  USER: 'users/:id',
  UPDATE_USER_ROLE: 'users/:id/role',
  ORDERS: 'orders',
  ORDER_DETAIL: 'orders/:id',
  ORDER_DETAIL_STATUS: 'orders/:id/status',
  ORDER_DETAIL_TRACKING: 'orders/:id/tracking',
  RATINGS: 'ratings',
  RATING_DETAIL: 'ratings/:id',
  RATING_APPROVE: 'ratings/:id/approve',
  RATING_REJECT: 'ratings/:id/ban',
  DISCOUNTS: 'discounts',
  DISCOUNT: 'discounts/:id',
  DISCOUNT_PRODUCTS: 'discounts/:id/products',
  DISCOUNT_CATEGORIES: 'discounts/:id/categories',
  DISCOUNT_USERS: 'discounts/:id/users',
  DISCOUNT_ADD_USERS: 'discounts/:id/users/add',
  DISCOUNT_REMOVE_USERS: 'discounts/:id/users/remove',
  DISCOUNT_ADD_PRODUCTS: 'discounts/:id/products/add',
  DISCOUNT_REMOVE_PRODUCTS: 'discounts/:id/products/remove',
  DISCOUNT_ADD_CATEGORIES: 'discounts/:id/categories/add',
  DISCOUNT_REMOVE_CATEGORIES: 'discounts/:id/categories/remove',
  DISCOUNT_ADD_ALL_PRODUCTS: 'discounts/:id/products/add-all',
  DISCOUNT_REMOVE_ALL_PRODUCTS: 'discounts/:id/products/remove-all',
  DISCOUNT_ADD_ALL_CATEGORIES: 'discounts/:id/categories/add-all',
  DISCOUNT_REMOVE_ALL_CATEGORIES: 'discounts/:id/categories/remove-all',
  DISCOUNT_ADD_ALL_USERS: 'discounts/:id/users/add-all',
  DISCOUNT_REMOVE_ALL_USERS: 'discounts/:id/users/remove-all',
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
