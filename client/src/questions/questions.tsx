import { useContext } from 'react';
import { useQuery } from 'react-query';
import { QueryClientContext, SessionContext } from '../providers';

export function Questions(): JSX.Element {
  const { logout } = useContext(SessionContext);
  const apiClient = useContext(QueryClientContext);

  useQuery('current-question', {
    queryFn: async () => {
      const q = await apiClient.getNextQuestion();
      console.log({ q });
      return q;
    },
  });

  return (
    <div>
      <button onClick={logout} className="btn btn-outline btn-primary">
        Sign Out
      </button>
    </div>
  );
}
