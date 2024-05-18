import { clsx } from 'clsx';
import { UnansweredQuestion } from '@/generated/api';

export interface Feedback {
  selectedChoiceId: string;
  correctChoiceId: string;
  moreInfo: string;
}

interface Props {
  question: UnansweredQuestion;
  feedback: Feedback;
  isLoadingNext: boolean;
  onNext: () => unknown;
}

export function QuestionFeedbackCard({
  question,
  feedback,
  isLoadingNext,
  onNext,
}: Props): JSX.Element {
  const isCorrect = feedback.selectedChoiceId === feedback.correctChoiceId;
  const correctChoice = question.choices.find(({ id }) => id === feedback.correctChoiceId);

  return (
    <div className="card bg-neutral shadow-xl">
      <div className="card-body">
        <h2 className={clsx('card-title', isCorrect ? 'text-success' : 'text-error')}>
          {isCorrect ? 'Correct!' : 'Wrong'}
        </h2>

        <div>
          {isCorrect ? null : (
            <>
              <p>
                Correct choice was <span className="font-bold">{correctChoice?.choice}</span>.
              </p>
              <div className="divider"></div>
            </>
          )}
          <p>{feedback.moreInfo}</p>
        </div>

        <div className="card-actions justify-end">
          <button onClick={onNext} className="btn btn-primary" disabled={isLoadingNext}>
            {isLoadingNext ? <span className="loading loading-spinner"></span> : null}
            Next
          </button>
        </div>
      </div>
    </div>
  );
}
