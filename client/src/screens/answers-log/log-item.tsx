import { CheckCircleIcon, XCircleIcon } from '@heroicons/react/24/solid';
import { type AnswersLogItem } from '@/generated/api';

interface Props {
  logItem: AnswersLogItem;
  onClick: () => unknown;
}

export function LogItem({ logItem, onClick }: Props): JSX.Element {
  const correctChoice = logItem.question.choices.find((choice) => choice.isCorrect);
  const isCorrect = logItem.choiceId === correctChoice?.id;

  return (
    <div className="flex items-center space-x-2 w-full">
      {isCorrect ? (
        <CheckCircleIcon className="text-success w-8 h-8 flex-shrink-0" />
      ) : (
        <XCircleIcon className="text-error w-8 h-8 flex-shrink-0" />
      )}
      <button onClick={onClick} className="link link-hover truncate max-w-full">
        {logItem.question.question}
      </button>
    </div>
  );
}
