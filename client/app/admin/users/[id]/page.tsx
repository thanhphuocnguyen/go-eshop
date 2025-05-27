import { Metadata } from 'next';
import { notFound } from 'next/navigation';
import { serverSideFetch } from '@/app/lib/api/apiServer';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import UserProfileForm from '../_components/UserProfileForm';
import dayjs from 'dayjs';
import AdminActions from '../_components/AdminActions';

// Define the user type
interface UserDetails {
  id: string;
  email: string;
  fullname: string;
  phone: string;
  role: string;
  createdAt: string;
  updatedAt: string;
  locked: boolean;
  avatarUrl: string | null;
  avatarImageId: string | null;
  // Additional user properties
  lastLoginAt: string | null;
  addresses: UserAddress[];
}

interface UserAddress {
  id: string;
  street: string;
  city: string;
  state: string;
  postalCode: string;
  country: string;
  isDefault: boolean;
}

export async function generateMetadata({
  params,
}: {
  params: { id: string };
}): Promise<Metadata> {
  // Fetch user data
  const user = await getUserDetails(params.id);

  if (!user) {
    return {
      title: 'User Not Found',
      description: 'The requested user profile could not be found.',
    };
  }

  return {
    title: `${user.fullname} | User Profile`,
    description: `Admin view of user profile for ${user.email}`,
  };
}

async function getUserDetails(userId: string): Promise<UserDetails | null> {
  try {
    const { data, error } = await serverSideFetch<UserDetails>(
      `${ADMIN_API_PATHS.USERS}/${userId}`
    );

    if (error || !data) {
      console.error('Error fetching user details:', error);
      return null;
    }

    return data;
  } catch (error) {
    console.error('Error fetching user:', error);
    return null;
  }
}

export default async function UserProfilePage({
  params,
}: {
  params: { id: string };
}) {
  const user = await getUserDetails(params.id);
  if (!user) {
    notFound();
  }

  return (
    <div className='container mx-auto px-4 py-6'>
      <div className='flex justify-between items-center mb-8'>
        <div>
          <h1 className='text-2xl font-bold text-gray-900'>User Profile</h1>
          <p className='text-gray-600 mt-1'>
            Manage and view details for {user.fullname}
          </p>
        </div>
        <div className='flex items-center space-x-2'>
          <span
            className={`px-3 py-1 text-xs font-medium rounded-full ${
              user.role === 'admin'
                ? 'bg-purple-100 text-purple-800'
                : user.role === 'staff'
                  ? 'bg-blue-100 text-blue-800'
                  : 'bg-green-100 text-green-800'
            }`}
          >
            {user.role.charAt(0).toUpperCase() + user.role.slice(1)}
          </span>
          {user.locked && (
            <span className='bg-red-100 text-red-800 px-3 py-1 text-xs font-medium rounded-full'>
              Account Locked
            </span>
          )}
        </div>
      </div>

      <div className='grid grid-cols-1 lg:grid-cols-3 gap-8'>
        {/* Left Column - User Info */}
        <div className='lg:col-span-2'>
          <div className='bg-white shadow rounded-lg overflow-hidden'>
            <div className='px-6 py-5 border-b border-gray-200'>
              <h2 className='text-lg font-medium text-gray-900'>
                User Information
              </h2>
            </div>
            <div className='p-6'>
              <UserProfileForm user={user} />
            </div>
          </div>

          <div className='bg-white shadow rounded-lg overflow-hidden mt-8'>
            <div className='px-6 py-5 border-b border-gray-200'>
              <h2 className='text-lg font-medium text-gray-900'>
                Order History
              </h2>
            </div>
          </div>
        </div>

        {/* Right Column - Admin Actions & User Stats */}
        <div>
          <div className='bg-white shadow rounded-lg overflow-hidden'>
            <div className='px-6 py-5 border-b border-gray-200'>
              <h2 className='text-lg font-medium text-gray-900'>
                Admin Actions
              </h2>
            </div>
            <div className='p-6'>
              <AdminActions user={user} />
            </div>
          </div>

          <div className='bg-white shadow rounded-lg overflow-hidden mt-8'>
            <div className='px-6 py-5 border-b border-gray-200'>
              <h2 className='text-lg font-medium text-gray-900'>
                User Statistics
              </h2>
            </div>
            <div className='p-6'>
              <dl className='divide-y divide-gray-200'>
                <div className='py-3 flex justify-between'>
                  <dt className='text-sm font-medium text-gray-500'>
                    Member Since
                  </dt>
                  <dd className='text-sm text-gray-900'>
                    {dayjs(user.createdAt).format('MMM D, YYYY')}
                  </dd>
                </div>
                <div className='py-3 flex justify-between'>
                  <dt className='text-sm font-medium text-gray-500'>
                    Last Login
                  </dt>
                  <dd className='text-sm text-gray-900'>
                    {user.lastLoginAt
                      ? dayjs(user.lastLoginAt).format('MMM d, yyyy h:mm a')
                      : 'Never'}
                  </dd>
                </div>
                <div className='py-3 flex justify-between'>
                  <dt className='text-sm font-medium text-gray-500'>
                    Total Orders
                  </dt>
                  <dd className='text-sm text-gray-900'>
                  </dd>
                </div>
                <div className='py-3 flex justify-between'>
                  <dt className='text-sm font-medium text-gray-500'>
                    Default Shipping Address
                  </dt>
                  <dd className='text-sm text-gray-900'>
                    {user.addresses && user.addresses.find((a) => a.isDefault)
                      ? `${user.addresses.find((a) => a.isDefault)?.city}, ${user.addresses.find((a) => a.isDefault)?.country}`
                      : 'None'}
                  </dd>
                </div>
              </dl>

              {/* User addresses section */}
              <div className='mt-4'>
                <h3 className='text-sm font-medium text-gray-900 mb-3'>
                  Addresses
                </h3>
                {user.addresses && user.addresses.length > 0 ? (
                  <div className='space-y-3'>
                    {user.addresses.map((address) => (
                      <div
                        key={address.id}
                        className={`p-3 rounded-md text-sm ${
                          address.isDefault
                            ? 'bg-blue-50 border border-blue-200'
                            : 'bg-gray-50 border border-gray-200'
                        }`}
                      >
                        <div className='flex justify-between items-start'>
                          <div>
                            <p>{address.street}</p>
                            <p>
                              {address.city}, {address.state}{' '}
                              {address.postalCode}
                            </p>
                            <p>{address.country}</p>
                          </div>
                          {address.isDefault && (
                            <span className='bg-blue-100 text-blue-700 px-2.5 py-0.5 rounded-full text-xs font-medium'>
                              Default
                            </span>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className='text-sm text-gray-500'>No addresses on file</p>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
