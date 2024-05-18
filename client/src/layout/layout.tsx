import { PropsWithChildren } from 'react';
import { Navbar } from './navbar';
import { Footer } from './footer';

export function Layout({ children }: PropsWithChildren): JSX.Element {
  return (
    <div className="bg-primary h-screen flex flex-col">
      <Navbar />
      <div className="container m-auto">{children}</div>
      <Footer />
    </div>
  );
}
