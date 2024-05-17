import { QueryClientProvider, SessionProvider } from '@/providers';
import { EnsureLoggedIn } from '@/screens/login';
import { Question } from '@/screens/question';
import { Layout } from './layout';

export default function App() {
  return (
    <SessionProvider>
      <QueryClientProvider>
        <EnsureLoggedIn>
          <Layout>
            <Question />
          </Layout>
        </EnsureLoggedIn>
      </QueryClientProvider>
    </SessionProvider>
  );
}
