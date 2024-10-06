import { Auth } from '@supabase/auth-ui-react';
import { ThemeSupa } from '@supabase/auth-ui-shared';
import { supabase } from '@/supabase';
import { NapoleonicHatIcon } from '@/icons';
import { Credits, WithSkybox } from '@/layout';

export function Login(): JSX.Element {
  return (
    <WithSkybox>
      <div className="h-full flex flex-col items-center justify-center">
        <div className="card bg-neutral shadow-xl flex flex-col items-center space-y-2 p-4 mx-2 opacity-90">
          <NapoleonicHatIcon height={160} width={200} />

          <h1 className="text-2xl font-bold tracking-tight">Quizory</h1>

          <h2 className="text-xl text-center">
            LLM-generated questions to challenge your knowledge!
          </h2>

          <div className="mx-auto mt-4">
            <Auth
              supabaseClient={supabase}
              providers={['google']}
              appearance={{ theme: ThemeSupa }}
              onlyThirdPartyProviders
            />
          </div>

          <Credits />
        </div>
      </div>
    </WithSkybox>
  );
}
