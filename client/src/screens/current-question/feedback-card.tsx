import { clsx } from 'clsx';
import { useForm } from 'react-hook-form';
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
  onNext: () => Promise<unknown>;
}

export function QuestionFeedbackCard({ question, feedback, onNext }: Props): JSX.Element {
  const { handleSubmit, formState } = useForm();

  const isCorrect = feedback.selectedChoiceId === feedback.correctChoiceId;
  const correctChoice = question.choices.find(({ id }) => id === feedback.correctChoiceId);

  return (
    <form
      className="card bg-neutral shadow-xl"
      onSubmit={handleSubmit(() => onNext())}
    >
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
              <div className="divider" />
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
            type="submit"
            disabled={formState.isSubmitting || !formState.isValid}
            className="btn btn-primary btn-block sm:btn-wide"
          >
            {formState.isSubmitting ? <InlineSpinner /> : null}
            Next
          </button>
        </div>
      </div>
    </form>
  );
}
