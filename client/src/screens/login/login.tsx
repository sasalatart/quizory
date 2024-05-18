import { PropsWithChildren, useContext } from 'react';
import { Auth } from '@supabase/auth-ui-react';
import { ThemeSupa } from '@supabase/auth-ui-shared';
import { supabase } from '@/supabase';
import { GITHUB_USER_LINK, GITHUB_REPO_LINK } from '@/config';
import { NapoleonicHatIcon } from '@/icons';
import { SessionContext } from '@/providers';

export function Login(): JSX.Element {
  return (
    <div className="h-screen flex flex-col justify-center">
      <div className="flex flex-col items-center space-y-2">
        <NapoleonicHatIcon height={160} width={200} />
        <h1 className="text-2xl font-bold tracking-tight">Quizory</h1>
        <h2 className="text-xl text-center">AI-generated questions to challenge your knowledge!</h2>
      </div>
      <div className="mx-auto mt-4">
        <Auth
          supabaseClient={supabase}
          providers={['google']}
          appearance={{ theme: ThemeSupa }}
          onlyThirdPartyProviders
        />
      </div>

      <div className="flex flex-col items-center mt-4">
        <p>
          Created by Sebastián Salata R-T{' '}
          <a href={GITHUB_USER_LINK} target="_blank" className="link">
            (sasalatart)
          </a>
        </p>
        <p>
          See code on{' '}
          <a href={GITHUB_REPO_LINK} target="_blank" className="link">
            GitHub
          </a>
        </p>
      </div>
    </div>
  );
}

export function EnsureLoggedIn({ children }: PropsWithChildren): JSX.Element {
  const { session } = useContext(SessionContext);
  return session ? <>{children}</> : <Login />;
}
