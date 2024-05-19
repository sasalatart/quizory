import { createBrowserRouter, Outlet } from 'react-router-dom';
import { QueryClientProvider, SessionProvider } from '@/providers';
import { Layout } from '@/layout';
import { Login } from '@/screens/login';
import { Question } from '@/screens/current-question';
import { AnswersLog } from '@/screens/answers-log';
import { ROUTES } from './routes';

export const router = createBrowserRouter([
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
        path: ROUTES.login,
        element: <Login />,
      },
      {
        path: ROUTES.questions.root,
        element: <Layout />,
        children: [
          {
            path: ROUTES.questions.current,
            element: <Question />,
          },
          {
            path: ROUTES.questions.log,
            element: <AnswersLog />,
          },
        ],
      },
    ],
  },
]);
