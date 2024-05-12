import { PropsWithChildren, createContext, useContext } from 'react';
import { QueryClientProvider as BaseReactQueryClientProvider, QueryClient } from 'react-query';
import { Configuration, DefaultApi } from '@/generated/api';
import { BASE_API_URL } from '../config';
import { SessionContext } from './session-provider';

const reactQueryClient = new QueryClient();

export const QueryClientContext = createContext<DefaultApi>(newApiClient(undefined));

function newApiClient(token: string | undefined): DefaultApi {
  return new DefaultApi(
    new Configuration({
      basePath: BASE_API_URL,
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    }),
  );
}

export function QueryClientProvider({ children }: PropsWithChildren): JSX.Element {
  const { session } = useContext(SessionContext);

  return (
    <QueryClientContext.Provider value={newApiClient(session?.access_token)}>
      <BaseReactQueryClientProvider client={reactQueryClient}>
        {children}
      </BaseReactQueryClientProvider>
    </QueryClientContext.Provider>
  );
}
