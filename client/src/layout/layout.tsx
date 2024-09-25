import { Outlet } from 'react-router-dom';
import { Navbar } from './navbar';
import { Footer } from './footer';
import { WithSkybox } from './skybox';

export function Layout(): JSX.Element {
  return (
    <WithSkybox>
      <Navbar />
      <div className="container m-auto">
        <div className="m-4">
          <Outlet />
        </div>
      </div>
      <Footer />
    </WithSkybox>
  );
}
