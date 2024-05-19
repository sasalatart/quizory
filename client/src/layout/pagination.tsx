import { useSearchParams } from 'react-router-dom';

interface Props {
  currentPage: number;
  hasNextPage: boolean;
  queryParamName?: string;
}

export function Pagination({ currentPage, hasNextPage, queryParamName }: Props): JSX.Element {
  const paramName = queryParamName ?? 'page';
  const [, setSearchParams] = useSearchParams();

  return (
    <div className="join grid grid-cols-2">
      <button
        onClick={() => setSearchParams({ [paramName]: `${currentPage - 1}` })}
        className="join-item btn btn-outline"
        disabled={currentPage == 0}
      >
        Previous page
      </button>
      <button
        onClick={() => setSearchParams({ [paramName]: `${currentPage + 1}` })}
        className="join-item btn btn-outline"
        disabled={!hasNextPage}
      >
        Next page
      </button>
    </div>
  );
}
