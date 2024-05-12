import { QueryClientProvider, SessionProvider } from '@/providers';
import { EnsureLoggedIn } from '@/screens/login';
import { Questions } from '@/screens/questions';

export default function App() {
  return (
    <SessionProvider>
      <QueryClientProvider>
        <EnsureLoggedIn>
          <Questions />
        </EnsureLoggedIn>
      </QueryClientProvider>
    </SessionProvider>
  );
}
