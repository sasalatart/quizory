import { clsx } from 'clsx';
import { UnansweredQuestion } from '@/generated/api';
import { InlineSpinner } from '@/layout';

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
          {feedback.moreInfo.split('\n').map((line, index) => (
            <p key={index} className="my-4">
              {line}
            </p>
          ))}
        </div>

        <div className="card-actions justify-center">
          <button
            onClick={onNext}
            disabled={isLoadingNext}
            className="btn btn-primary btn-block sm:btn-wide"
          >
            {isLoadingNext ? <InlineSpinner /> : null}
            Next
          </button>
        </div>
      </div>
    </div>
  );
}
