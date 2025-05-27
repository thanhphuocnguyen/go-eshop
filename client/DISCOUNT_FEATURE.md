# Discount Feature Implementation

## Overview
Added a comprehensive discount system to the checkout page with loading UI and discount selection functionality.

## Features Added

### 1. Loading UI for Discounts
- Shows a spinner with "Loading discounts..." message while fetching discount data
- Graceful error handling for discount fetch failures

### 2. Discount List Box
- Displays all available discounts in a scrollable container
- Each discount shows:
  - Discount code (clickable to apply)
  - Discount value (percentage or fixed amount)
  - Description
  - Expiration date
- Click on any discount to automatically apply it

### 3. Discount Application
- Text input field for manual discount code entry
- "Apply" button to apply the entered discount code
- Auto-application when selecting from the discount list
- Real-time validation of discount codes

### 4. Applied Discount Display
- Shows applied discount in the order summary
- Displays discount code and amount
- Remove button (✕) to remove the applied discount
- Updates total price calculation

### 5. Form Integration
- Added `discount_code` field to checkout form schema
- Discount information is sent to the backend during checkout
- Proper form validation and error handling

## Technical Implementation

### Schema Updates
```typescript
// Added to CheckoutFormSchema
discount_code: z.string().optional(),
```

### State Management
```typescript
const [appliedDiscount, setAppliedDiscount] = useState<DiscountModel | null>(null);
const [discountAmount, setDiscountAmount] = useState(0);
```

### Discount Calculation
- Percentage discounts: `(subtotal * discountValue) / 100`
- Fixed amount discounts: `Math.min(discountValue, subtotal)`

### Auto-application Effect
- Monitors discount code changes
- Automatically applies valid discounts
- Removes invalid or cleared discount codes

## UI Components

### Loading State
```jsx
{discountsLoading ? (
  <div className='flex items-center justify-center p-4 bg-gray-50 rounded-md'>
    <div className='animate-spin rounded-full h-5 w-5 border-b-2 border-indigo-600 mr-2'></div>
    <span className='text-sm text-gray-500'>Loading discounts...</span>
  </div>
) : ...}
```

### Discount List
- Scrollable container with max height
- Hover effects and smooth transitions
- Visual badges for discount types
- Click-to-apply functionality

### Order Summary Integration
- Shows applied discount with remove option
- Updates final total calculation
- Color-coded (green) for applied discounts

## Backend Integration
- Sends applied discount code in checkout request
- Handles discount validation server-side
- Updates order total with applied discounts

## Error Handling
- Toast notifications for success/error states
- Graceful handling of invalid discount codes
- Loading states for better UX

## Styling
- Consistent with existing design system
- Responsive layout
- Accessible color contrasts
- Smooth animations and transitions
