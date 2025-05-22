import React from 'react';

type Product = {
  id: string;
  name: string;
  price: number;
};

type ProductsTableProps = {
  products?: Product[];
  discountType: string;
  discountValue: number;
};

export const ProductsTable: React.FC<ProductsTableProps> = ({
  products,
  discountType,
  discountValue,
}) => {
  if (!products || products.length === 0) {
    return (
      <div className='bg-gray-50 rounded-md p-8 text-center'>
        <p className='text-gray-500'>
          This discount applies to all products.
        </p>
      </div>
    );
  }

  return (
    <div className='overflow-x-auto'>
      <table className='min-w-full divide-y divide-gray-200'>
        <thead className='bg-gray-50'>
          <tr>
            <th
              scope='col'
              className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
            >
              Product
            </th>
            <th
              scope='col'
              className='px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider'
            >
              Price
            </th>
            <th
              scope='col'
              className='px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider'
            >
              Discounted Price
            </th>
          </tr>
        </thead>
        <tbody className='bg-white divide-y divide-gray-200'>
          {products.map((product) => (
            <tr key={product.id}>
              <td className='px-6 py-4 whitespace-nowrap'>
                <div className='text-sm font-medium text-gray-900'>
                  {product.name}
                </div>
                <div className='text-xs text-gray-500'>
                  {product.id}
                </div>
              </td>
              <td className='px-6 py-4 whitespace-nowrap text-right text-sm text-gray-500'>
                ${product.price.toFixed(2)}
              </td>
              <td className='px-6 py-4 whitespace-nowrap text-right text-sm text-gray-900'>
                {discountType === 'percentage'
                  ? `$${(product.price * (1 - discountValue / 100)).toFixed(2)}`
                  : `$${Math.max(0, product.price - discountValue).toFixed(2)}`}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
