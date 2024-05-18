import { createBrowserRouter, Outlet, RouterProvider } from 'react-router-dom';
import { QueryClientProvider, SessionProvider } from '@/providers';
import { Layout } from '@/layout';
import { Login } from '@/screens/login';
import { Question } from '@/screens/next-question';

const router = createBrowserRouter([
  {
    path: '/',
    element: (
      <>
        <SessionProvider>
          <QueryClientProvider>
            <Outlet />
          </QueryClientProvider>
        </SessionProvider>
      </>
    ),
    children: [
      {
        path: '',
        element: <Login />,
      },
      {
        path: 'questions',
        element: <Layout />,
        children: [
          {
            path: 'next',
            element: <Question />,
          },
        ],
      },
    ],
  },
]);

export default function App() {
  return <RouterProvider router={router} />;
}
