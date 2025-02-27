import Image from 'next/image';
import Link from 'next/link';

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div>
      <nav className='side-bar'>
        <div className='icon'>
          <Link href={'/'}>
            <Image src='/logo.png' alt='logo' width={200} height={200} />
          </Link>
        </div>
        <ul className='side-bar-list'>
          <li className='side-bar-item'>Dashboard</li>
          <li className='side-bar-item'>Users</li>
          <li className='side-bar-item'>Categories</li>
          <li className='side-bar-item'>Products</li>
          <li className='side-bar-item'>Collections</li>
          <li className='side-bar-item'>Brands</li>
          <li className='side-bar-item'>Attributes</li>
          <li className='side-bar-item'>Orders</li>
        </ul>
      </nav>
      <main>
        <section className='header'>
          <div>Search</div>
          <div>
            <div>Notifications</div>
            <hr />
            <div>Avatar</div>
            <div>Name</div>
            <span>chevron down</span>
          </div>
        </section>
        <section className='content-container'>{children}</section>
      </main>
    </div>
  );
}
