import { useContext, useEffect, useRef, useState } from 'react';
import { useQuery } from 'react-query';
import { useSearchParams } from 'react-router-dom';
import { type AnswersLogItem } from '@/generated/api';
import { CenteredSpinner, Pagination } from '@/layout';
import { QueryClientContext, SessionContext } from '@/providers';
import { LogItem } from './log-item';
import { LogItemModal } from './log-item-modal';

const PAGE_KEY = 'page';
const PAGE_SIZE = 10;

function useAnswersLog(pageNumber: number) {
  const { answersApi } = useContext(QueryClientContext);
  const { session } = useContext(SessionContext);

  return useQuery(`page-${pageNumber}-log`, {
    queryFn: () =>
      answersApi.getAnswersLog({
        page: pageNumber,
        pageSize: PAGE_SIZE,
        userId: session!.user.id,
      }),
    refetchOnWindowFocus: false,
  });
}

export function AnswersLog(): JSX.Element {
  const [searchParams] = useSearchParams();
  const [selectedLogItem, setSelectedLogItem] = useState<AnswersLogItem | undefined>();

  const currentPage = searchParams.get(PAGE_KEY) ? parseInt(searchParams.get(PAGE_KEY)!) : 0;
  const nextPage = currentPage + 1;

  const { data: currentPageLog, isLoading: isLoadingCurrentPage } = useAnswersLog(currentPage);
  const { data: nextPageLog } = useAnswersLog(nextPage);

  const modalRef = useRef<HTMLDialogElement>(null);
  useEffect(() => {
    if (selectedLogItem) {
      modalRef.current?.showModal();
    } else {
      modalRef.current?.close();
    }
  }, [selectedLogItem]);

  return (
    <div className="card bg-neutral shadow-xl">
      <div className="card-body">
        <h1 className="card-title">Your previous answers</h1>

        {isLoadingCurrentPage ? (
          <CenteredSpinner />
        ) : (
          (currentPageLog ?? []).map((logItem) => (
            <LogItem
              key={logItem.id}
              logItem={logItem}
              onClick={() => setSelectedLogItem(logItem)}
            />
          ))
        )}

        <Pagination
          currentPage={currentPage}
          hasNextPage={!!nextPageLog && nextPageLog.length > 0}
          queryParamName={PAGE_KEY}
        />

        <LogItemModal
          logItem={selectedLogItem}
          onClose={() => setSelectedLogItem(undefined)}
          modalRef={modalRef}
        />
      </div>
    </div>
  );
}
