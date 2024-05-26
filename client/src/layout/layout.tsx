import { Outlet } from 'react-router-dom';
import { Navbar } from './navbar';
import { Footer } from './footer';

export function Layout(): JSX.Element {
  return (
    <div className="bg-primary h-screen flex flex-col">
      <Navbar />
      <div className="container m-auto">
        <div className="m-4">
          <Outlet />
        </div>
      </div>
      <Footer />
    </div>
  );
}
