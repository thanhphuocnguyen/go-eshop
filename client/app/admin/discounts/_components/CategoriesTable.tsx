import React from 'react';

type Category = {
  id: string;
  name: string;
};

type CategoriesTableProps = {
  categories?: Category[];
};

export const CategoriesTable: React.FC<CategoriesTableProps> = ({
  categories,
}) => {
  if (!categories || categories.length === 0) {
    return (
      <div className='bg-gray-50 rounded-md p-8 text-center'>
        <p className='text-gray-500'>
          No specific categories selected for this discount.
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
              Category Name
            </th>
            <th
              scope='col'
              className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
            >
              Category ID
            </th>
          </tr>
        </thead>
        <tbody className='bg-white divide-y divide-gray-200'>
          {categories.map((category) => (
            <tr key={category.id}>
              <td className='px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900'>
                {category.name}
              </td>
              <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-500'>
                {category.id}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
