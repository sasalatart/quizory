import { clsx } from 'clsx';
import startCase from 'lodash/startCase';
import { Controller, useForm } from 'react-hook-form';
import { RemainingTopic, UnansweredQuestion } from '@/generated/api';
import { InlineSpinner } from '@/layout';

export interface Feedback {
  selectedChoiceId: string;
  correctChoiceId: string;
  moreInfo: string;
}

interface Props {
  question: UnansweredQuestion;
  feedback: Feedback;
  remainingTopics: RemainingTopic[];
  onNext: (nextTopic: string | undefined) => Promise<unknown>;
}

interface Form {
  nextTopic: string | undefined;
}

/**
 * Returns the default next topic based on the current question and the remaining topics: If there
 * are more questions for the current topic, it returns the current question's topic. If not, then
 * it falls back to the first available topic (if any).
 */
function getDefaultNextTopic(
  currentQuestion: { topic: string },
  remainingTopics: RemainingTopic[],
) {
  const currentTopicHasMoreQuestions = remainingTopics.some(
    ({ topic }) => topic === currentQuestion.topic,
  );
  if (currentTopicHasMoreQuestions) {
    return currentQuestion.topic;
  }

  return remainingTopics[0]?.topic;
}

export function QuestionFeedbackCard({
  question,
  feedback,
  remainingTopics,
  onNext,
}: Props): JSX.Element {
  const { control, handleSubmit, formState } = useForm<Form>();

  const isCorrect = feedback.selectedChoiceId === feedback.correctChoiceId;
  const correctChoice = question.choices.find(({ id }) => id === feedback.correctChoiceId);

  return (
    <form
      className="card bg-neutral shadow-xl"
      onSubmit={handleSubmit((data) => onNext(data.nextTopic))}
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

        <div className="divider" />

        <div className="card-actions justify-center sm:justify-end">
          {remainingTopics.length > 0 ? (
            <label className="form-control w-full sm:max-w-64">
              <Controller
                name="nextTopic"
                control={control}
                rules={{ required: true }}
                defaultValue={getDefaultNextTopic(question, remainingTopics) ?? ''}
                render={({ field }) => (
                  <select {...field} className="select select-bordered">
                    <option value="">Select next topic</option>
                    {remainingTopics.map(({ topic, amountOfQuestions }) => (
                      <option key={topic} value={topic}>
                        {startCase(topic)} ({amountOfQuestions})
                      </option>
                    ))}
                  </select>
                )}
              />

              <div className="label">
                <span className="label-text">Next topic</span>
              </div>
            </label>
          ) : null}

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
