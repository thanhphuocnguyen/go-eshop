'use client';

import { API_PATHS } from '@/lib/constants/api';
import { toast } from 'react-toastify';

import dayjs from 'dayjs';
import { Button } from '@headlessui/react';
import useSWR from 'swr';
import Link from 'next/link';
import { GeneralCategoryModel, GenericResponse } from '@/lib/definitions';
import { ReactNode } from 'react';
import { apiFetch } from '@/lib/api/api';

export default function Page() {
  const { data: brands, isLoading } = useSWR<GeneralCategoryModel[]>(
    API_PATHS.BRANDS,
    (url) =>
      apiFetch<GenericResponse<GeneralCategoryModel[]>>(url).then(
        (res) => res.data
      ),
    {
      onError: (error) => {
        toast(<div>{error.message}</div>, { type: 'error' });
      },
    }
  );
  if (isLoading) return <div>Loading...</div>;
  return (
    <div className='h-full'>
      <div className='flex justify-between pt-4 pb-8'>
        <h2 className='text-lg font-bold'>Brands</h2>
        <Button
          as={Link}
          href={'/admin/brands/new'}
          className='btn btn-lg btn-primary'
        >
          Add new
        </Button>
      </div>
      <TableContainer
        header={['Name', 'Description', 'Slug', 'Created At', 'Actions']}
      >
        {brands?.map((e) => (
          <TableRow key={e.id}>
            <TableCellHead>
              <Link
                className='text-blue-500 underline underline-offset-2'
                href={'/admin/brands/' + e.id}
              >
                {e.name}
              </Link>
            </TableCellHead>
            <TableCell>{e.description}</TableCell>
            <TableCell>{e.slug}</TableCell>
            <TableCell>{dayjs(e.created_at).format('YYYY/MM/DD')}</TableCell>
            <TableCell className='flex space-x-2'>
              <Button className='btn btn-danger'>Edit</Button>
            </TableCell>
          </TableRow>
        ))}
      </TableContainer>
    </div>
  );
}

// Table components
interface TableContainerProps {
  header: string[];
  children: ReactNode;
}

function TableContainer({ header, children }: TableContainerProps) {
  return (
    <div className='overflow-x-auto'>
      <table className='min-w-full divide-y divide-gray-200'>
        <thead className='bg-gray-50'>
          <tr>
            {header.map((item, index) => (
              <th
                key={index}
                scope='col'
                className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
              >
                {item}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className='bg-white divide-y divide-gray-200'>{children}</tbody>
      </table>
    </div>
  );
}

function TableRow({ children }: { children: ReactNode }) {
  return <tr>{children}</tr>;
}

function TableCellHead({ children }: { children: ReactNode }) {
  return (
    <td className='px-6 py-4 whitespace-nowrap font-medium text-sm'>
      {children}
    </td>
  );
}

function TableCell({
  children,
  className = '',
}: {
  children?: ReactNode;
  className?: string;
}) {
  return (
    <td
      className={`px-6 py-4 whitespace-nowrap text-sm text-gray-500 ${className}`}
    >
      {children}
    </td>
  );
}
