'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { Dialog, DialogBackdrop } from '@headlessui/react';
import { toast } from 'react-toastify';

interface UserDetails {
  id: string;
  email: string;
  fullname: string;
  role: string;
  locked: boolean;
}

export default function AdminActions({ user }: { user: UserDetails }) {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [showResetModal, setShowResetModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [showPromoteModal, setShowPromoteModal] = useState(false);

  // Toggle user locked status
  const toggleLockStatus = async () => {
    setLoading(true);
    const { error } = await clientSideFetch<boolean>(
      `${ADMIN_API_PATHS.USERS}/${user.id}/lock`,
      {
        method: user.locked ? 'DELETE' : 'POST',
      }
    );

    if (error) {
      throw new Error(error.details || 'Failed to update user status');
    }

    toast.success(
      user.locked
        ? 'User account has been unlocked'
        : 'User account has been locked'
    );
    router.refresh();

    setLoading(false);
  };

  // Send password reset email
  const sendPasswordReset = async () => {
    setLoading(true);
    try {
      const response = await clientSideFetch<boolean>(
        `${ADMIN_API_PATHS.USERS}/${user.id}/reset-password`,
        {
          method: 'POST',
        }
      );

      if (response.error) {
        throw new Error(
          response.error.details || 'Failed to send password reset'
        );
      }

      toast.success('Password reset link has been sent to the user');
      setShowResetModal(false);
    } catch (error: any) {
      toast.error(
        error.message || 'An error occurred while sending the reset link'
      );
    } finally {
      setLoading(false);
    }
  };

  // Delete user account
  const deleteUser = async () => {
    setLoading(true);
    try {
      const response = await clientSideFetch<boolean>(
        `${ADMIN_API_PATHS.USERS}/${user.id}`,
        {
          method: 'DELETE',
        }
      );

      if (response.error) {
        throw new Error(response.error.details || 'Failed to delete user');
      }

      toast.success('User has been deleted successfully');
      setShowDeleteModal(false);
      router.push('/admin/users');
    } catch (error: any) {
      toast.error(error.message || 'An error occurred while deleting the user');
    } finally {
      setLoading(false);
    }
  };

  // Promote user to moderator
  const promoteToModerator = async () => {
    setLoading(true);
    try {
      const response = await clientSideFetch<boolean>(
        `${ADMIN_API_PATHS.USERS}/${user.id}/role`,
        {
          method: 'PUT',
          body: { role: 'moderator' },
        }
      );

      if (response.error) {
        throw new Error(response.error.details || 'Failed to update user role');
      }

      toast.success('User has been promoted to Moderator');
      setShowPromoteModal(false);
      router.refresh();
    } catch (error: any) {
      toast.error(
        error.message || 'An error occurred while updating the user role'
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className='space-y-4'>
      {/* Show promote button only if user is not already an admin or moderator */}
      {user.role !== 'admin' && user.role !== 'moderator' && (
        <button
          onClick={() => setShowPromoteModal(true)}
          disabled={loading}
          className='w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-purple-600 hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500'
        >
          Promote to Moderator
        </button>
      )}

      <button
        onClick={toggleLockStatus}
        disabled={loading}
        className={`w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white ${
          user.locked
            ? 'bg-green-600 hover:bg-green-700 focus:ring-green-500'
            : 'bg-red-600 hover:bg-red-700 focus:ring-red-500'
        } focus:outline-none focus:ring-2 focus:ring-offset-2`}
      >
        {loading
          ? 'Processing...'
          : user.locked
            ? 'Unlock Account'
            : 'Lock Account'}
      </button>

      <button
        onClick={() => setShowResetModal(true)}
        disabled={loading}
        className='w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
      >
        Send Password Reset
      </button>

      <button
        onClick={() => router.push(`/admin/users/${user.id}/send-email`)}
        className='w-full flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
      >
        Send Email
      </button>

      <div className='border-t border-gray-200 pt-4'>
        <button
          onClick={() => setShowDeleteModal(true)}
          className='w-full flex justify-center py-2 px-4 border border-gray-300 rounded-md text-sm font-medium text-red-600 hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500'
        >
          Delete Account
        </button>
      </div>

      {/* Promote User Modal */}
      <Dialog
        open={showPromoteModal}
        onClose={() => setShowPromoteModal(false)}
        className='fixed z-10 inset-0 overflow-y-auto'
      >
        <div className='flex items-center justify-center min-h-screen'>
          <DialogBackdrop className='fixed inset-0 bg-black opacity-30' />

          <div className='relative bg-white rounded-lg max-w-md w-full mx-auto p-6 shadow-xl'>
            <Dialog.Title className='text-lg font-medium text-gray-900'>
              Promote User to Moderator
            </Dialog.Title>
            <Dialog.Description className='mt-2 text-sm text-gray-500'>
              Are you sure you want to promote {user.fullname} to a Moderator
              role? This will give them additional permissions to manage content
              and users.
            </Dialog.Description>

            <div className='mt-6 flex space-x-4 justify-end'>
              <button
                type='button'
                className='inline-flex justify-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
                onClick={() => setShowPromoteModal(false)}
              >
                Cancel
              </button>
              <button
                type='button'
                className='inline-flex justify-center px-4 py-2 text-sm font-medium text-white bg-purple-600 border border-transparent rounded-md hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500'
                onClick={promoteToModerator}
                disabled={loading}
              >
                {loading ? 'Promoting...' : 'Promote User'}
              </button>
            </div>
          </div>
        </div>
      </Dialog>

      {/* Password Reset Confirmation Modal */}
      <Dialog
        open={showResetModal}
        onClose={() => setShowResetModal(false)}
        className='fixed z-10 inset-0 overflow-y-auto'
      >
        <div className='flex items-center justify-center min-h-screen'>
          <DialogBackdrop className='fixed inset-0 bg-black opacity-30' />

          <div className='relative bg-white rounded-lg max-w-md w-full mx-auto p-6 shadow-xl'>
            <Dialog.Title className='text-lg font-medium text-gray-900'>
              Send Password Reset
            </Dialog.Title>
            <Dialog.Description className='mt-2 text-sm text-gray-500'>
              Are you sure you want to send a password reset email to{' '}
              {user.email}?
            </Dialog.Description>

            <div className='mt-6 flex space-x-4 justify-end'>
              <button
                type='button'
                className='inline-flex justify-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
                onClick={() => setShowResetModal(false)}
              >
                Cancel
              </button>
              <button
                type='button'
                className='inline-flex justify-center px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
                onClick={sendPasswordReset}
                disabled={loading}
              >
                {loading ? 'Sending...' : 'Send Reset Link'}
              </button>
            </div>
          </div>
        </div>
      </Dialog>

      {/* Delete User Confirmation Modal */}
      <Dialog
        open={showDeleteModal}
        onClose={() => setShowDeleteModal(false)}
        className='fixed z-10 inset-0 overflow-y-auto'
      >
        <div className='flex items-center justify-center min-h-screen'>
          <DialogBackdrop className='fixed inset-0 bg-black opacity-30' />

          <div className='relative bg-white rounded-lg max-w-md w-full mx-auto p-6 shadow-xl'>
            <Dialog.Title className='text-lg font-medium text-gray-900'>
              Delete User Account
            </Dialog.Title>
            <Dialog.Description className='mt-2 text-sm text-gray-500'>
              Are you sure you want to delete this user? This action cannot be
              undone and will remove all data associated with {user.fullname}.
            </Dialog.Description>

            <div className='mt-6 flex space-x-4 justify-end'>
              <button
                type='button'
                className='inline-flex justify-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
                onClick={() => setShowDeleteModal(false)}
              >
                Cancel
              </button>
              <button
                type='button'
                className='inline-flex justify-center px-4 py-2 text-sm font-medium text-white bg-red-600 border border-transparent rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500'
                onClick={deleteUser}
                disabled={loading}
              >
                {loading ? 'Deleting...' : 'Delete User'}
              </button>
            </div>
          </div>
        </div>
      </Dialog>
    </div>
  );
}
