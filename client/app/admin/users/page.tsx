'use client';

import { useState, useEffect } from 'react';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { UserModel } from '@/app/lib/definitions';
import { Button } from '@headlessui/react';
import Link from 'next/link';
import Loading from '@/app/loading';
import { toast } from 'react-toastify';
import { MagnifyingGlassIcon } from '@heroicons/react/24/outline';
import { ConfirmDialog } from '@/app/components/Common/Dialogs/ConfirmDialog';
import dayjs from 'dayjs';

export default function Page() {
  const [users, setUsers] = useState<UserModel[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchInput, setSearchInput] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [selectedUser, setSelectedUser] = useState<UserModel | null>(null);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);

  // Pagination settings
  const pageSize = 10;

  const fetchUsers = async (page = 1, search = '') => {
    setLoading(true);
    try {
      // Using the admin users API endpoint
      const response = await clientSideFetch<UserModel[]>(
        `${ADMIN_API_PATHS.USERS}?page=${page}&pageSize=${pageSize}${search ? `&search=${search}` : ''}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );

      if (response.error) {
        toast.error('Failed to load users: ' + response.error.details);
        return;
      }

      setUsers(response.data || []);

      // If pagination data is available in the response
      if (response.pagination) {
        setTotalPages(Math.ceil(response.pagination.total / pageSize));
      }
    } catch (error) {
      console.error('Error fetching users:', error);
      toast.error('An error occurred while loading users');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers(currentPage, searchInput);
  }, [currentPage]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setCurrentPage(1); // Reset to first page on new search
    fetchUsers(1, searchInput);
  };

  const handleDeleteConfirm = async () => {
    if (!selectedUser) return;

    try {
      const response = await clientSideFetch<boolean>(
        ADMIN_API_PATHS.USER.replace(':id', selectedUser.id),
        {
          method: 'DELETE',
        }
      );

      if (response.error) {
        toast.error('Failed to delete user: ' + response.error.details);
        return;
      }

      toast.success('User deleted successfully');
      setUsers(users.filter((user) => user.id !== selectedUser.id));
      setShowDeleteDialog(false);
      setSelectedUser(null);
    } catch (error) {
      console.error('Error deleting user:', error);
      toast.error('An error occurred while deleting the user');
    }
  };

  const handlePageChange = (newPage: number) => {
    if (newPage >= 1 && newPage <= totalPages) {
      setCurrentPage(newPage);
    }
  };

  if (loading && users.length === 0) {
    return <Loading />;
  }

  return (
    <div className='container mx-auto px-4'>
      <div className='flex justify-between items-center mb-6'>
        <h2 className='text-2xl font-semibold text-primary'>
          Users Management
        </h2>
      </div>

      {/* Search Filter */}
      <div className='mb-6'>
        <form onSubmit={handleSearch} className='flex gap-2'>
          <div className='relative flex-grow'>
            <div className='absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none'>
              <MagnifyingGlassIcon className='w-5 h-5 text-gray-500' />
            </div>
            <input
              type='text'
              className='bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full pl-10 p-2.5'
              placeholder='Search users by name, email, or username...'
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
            />
          </div>
          <Button
            type='submit'
            className='px-4 py-2.5 bg-blue-500 text-white rounded-lg hover:bg-blue-600'
          >
            Search
          </Button>
        </form>
      </div>

      {/* Users Table */}
      <div className='relative overflow-x-auto shadow-md sm:rounded-lg'>
        <table className='w-full text-sm text-left text-gray-500'>
          <thead className='text-xs text-gray-700 uppercase bg-gray-50'>
            <tr>
              <th scope='col' className='px-6 py-3'>
                ID
              </th>
              <th scope='col' className='px-6 py-3'>
                Name
              </th>
              <th scope='col' className='px-6 py-3'>
                Username
              </th>
              <th scope='col' className='px-6 py-3'>
                Email
              </th>
              <th scope='col' className='px-6 py-3'>
                Role
              </th>
              <th scope='col' className='px-6 py-3'>
                Verified
              </th>
              <th scope='col' className='px-6 py-3'>
                Created
              </th>
              <th scope='col' className='px-6 py-3'>
                Actions
              </th>
            </tr>
          </thead>
          <tbody>
            {users.length > 0 ? (
              users.map((user) => (
                <tr
                  key={user.id}
                  className='bg-white border-b hover:bg-gray-50'
                >
                  <td className='px-6 py-4 font-medium text-gray-900 whitespace-nowrap'>
                    {user.id.substring(0, 8)}...
                  </td>
                  <td className='px-6 py-4'>{user.fullname}</td>
                  <td className='px-6 py-4'>{user.username}</td>
                  <td className='px-6 py-4'>{user.email}</td>
                  <td className='px-6 py-4'>
                    <span
                      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                        user.role === 'admin'
                          ? 'bg-purple-100 text-purple-800'
                          : user.role === 'moderator'
                            ? 'bg-blue-100 text-blue-800'
                            : 'bg-green-100 text-green-800'
                      }`}
                    >
                      {user.role}
                    </span>
                  </td>
                  <td className='px-6 py-4'>
                    <div className='flex space-x-1'>
                      <span
                        className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${
                          user.verifiedEmail
                            ? 'bg-green-100 text-green-800'
                            : 'bg-yellow-100 text-amber-800'
                        }`}
                      >
                        {user.verifiedEmail ? 'Email ✓' : 'Email ✗'}
                      </span>
                      <span
                        className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${
                          user.verifiedPhone
                            ? 'bg-green-100 text-green-800'
                            : 'bg-yellow-100 text-amber-500'
                        }`}
                      >
                        {user.verifiedPhone ? 'Phone ✓' : 'Phone ✗'}
                      </span>
                    </div>
                  </td>
                  <td className='px-6 py-4'>
                    {user.createdAt
                      ? dayjs(user.createdAt).format('YYYY/MM/DD')
                      : 'Unknown'}
                  </td>
                  <td className='px-6 py-4'>
                    <div className='flex space-x-2'>
                      <Link
                        href={`/admin/users/${user.id}`}
                        className='font-medium text-blue-600 hover:underline mr-2'
                      >
                        Edit
                      </Link>
                      <button
                        onClick={() => {
                          setSelectedUser(user);
                          setShowDeleteDialog(true);
                        }}
                        className='font-medium text-red-600 hover:underline'
                        // Prevent deleting admin users or your own account
                        disabled={user.role === 'admin'}
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))
            ) : (
              <tr className='odd:bg-white even:bg-gray-50 border-b'>
                <td colSpan={8} className='px-6 py-4 text-center'>
                  {loading ? 'Loading users...' : 'No users found'}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className='flex justify-center mt-6'>
          <nav aria-label='Page navigation'>
            <ul className='inline-flex items-center -space-x-px'>
              <li>
                <button
                  onClick={() => handlePageChange(currentPage - 1)}
                  disabled={currentPage === 1}
                  className='block px-3 py-2 ml-0 leading-tight text-gray-500 bg-white border border-gray-300 rounded-l-lg hover:bg-gray-100 hover:text-gray-700 disabled:opacity-50'
                >
                  Previous
                </button>
              </li>
              {[...Array(totalPages)].map((_, i) => (
                <li key={i}>
                  <button
                    onClick={() => handlePageChange(i + 1)}
                    className={`px-3 py-2 leading-tight border border-gray-300 
                      ${
                        currentPage === i + 1
                          ? 'text-blue-600 bg-blue-50 hover:bg-blue-100 hover:text-blue-700'
                          : 'text-gray-500 bg-white hover:bg-gray-100 hover:text-gray-700'
                      }`}
                  >
                    {i + 1}
                  </button>
                </li>
              ))}
              <li>
                <button
                  onClick={() => handlePageChange(currentPage + 1)}
                  disabled={currentPage === totalPages}
                  className='block px-3 py-2 leading-tight text-gray-500 bg-white border border-gray-300 rounded-r-lg hover:bg-gray-100 hover:text-gray-700 disabled:opacity-50'
                >
                  Next
                </button>
              </li>
            </ul>
          </nav>
        </div>
      )}

      {/* Confirmation Dialog */}
      <ConfirmDialog
        open={showDeleteDialog}
        title='Delete User'
        message={`Are you sure you want to delete the user "${selectedUser?.fullname || selectedUser?.username}"? This action cannot be undone.`}
        onClose={() => {
          setShowDeleteDialog(false);
          setSelectedUser(null);
        }}
        onConfirm={handleDeleteConfirm}
        confirmStyle='btn-danger'
      />
    </div>
  );
}
