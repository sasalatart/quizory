import { EnsureLoggedIn } from './auth';
import { QueryClientProvider, SessionProvider } from './providers';
import { Questions } from './questions';

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
