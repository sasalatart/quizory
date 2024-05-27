import { PropsWithChildren, createContext, useContext, useMemo } from 'react';
import { QueryClientProvider as BaseReactQueryClientProvider, QueryClient } from 'react-query';
import { Configuration, AnswersApi, QuestionsApi } from '@/generated/api';
import { BASE_API_URL } from '../config';
import { SessionContext } from './session-provider';

const reactQueryClient = new QueryClient();

interface Context {
  answersApi: AnswersApi;
  questionsApi: QuestionsApi;
}

export const QueryClientContext = createContext<Context>(newContext(undefined));

function newApiConfig(token: string | undefined): Configuration {
  return new Configuration({
    basePath: BASE_API_URL,
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  });
}

function newContext(token: string | undefined): Context {
  return {
    answersApi: new AnswersApi(newApiConfig(token)),
    questionsApi: new QuestionsApi(newApiConfig(token)),
  };
}

export function QueryClientProvider({ children }: PropsWithChildren): JSX.Element {
  const { session } = useContext(SessionContext);

  const context = useMemo<Context>(
    () => newContext(session?.access_token),
    [session?.access_token],
  );

  return (
    <QueryClientContext.Provider value={context}>
      <BaseReactQueryClientProvider client={reactQueryClient}>
        {children}
      </BaseReactQueryClientProvider>
    </QueryClientContext.Provider>
  );
}
