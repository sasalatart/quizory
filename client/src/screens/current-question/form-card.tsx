import { useContext } from 'react';
import { useMutation } from 'react-query';
import { useForm } from 'react-hook-form';
import { UnansweredQuestion } from '@/generated/api';
import { InlineSpinner, QuestionBadges } from '@/layout';
import { QueryClientContext } from '@/providers';
import { HintButton } from './hint-button';
import { Feedback } from './feedback-card';

interface Props {
  question: UnansweredQuestion;
  onSubmit: (feedback: Feedback) => unknown;
}

interface Form {
  choiceId: string;
}

export function QuestionFormCard({ question, onSubmit }: Props): JSX.Element {
  const apiClient = useContext(QueryClientContext);

  const { mutateAsync: submitAnswer } = useMutation(
    ({ choiceId }: Form) => apiClient.submitAnswer({ submitAnswerRequest: { choiceId } }),
    {
      onSuccess: ({ correctChoiceId, moreInfo }, { choiceId }) => {
        onSubmit({ correctChoiceId, selectedChoiceId: choiceId, moreInfo });
      },
    },
  );

  const { register, handleSubmit, formState } = useForm<Form>();
  return (
    <form
      onSubmit={handleSubmit((data) => submitAnswer(data))}
      className="card bg-neutral shadow-xl"
    >
      <div className="card-body">
        <QuestionBadges topic={question.topic} difficulty={question.difficulty} />
        <h2 className="card-title">{question.question}</h2>

        {question.choices.map(({ id, choice }) => (
          <div key={id} className="form-control">
            <label className="flex justify-start cursor-pointer label">
              <input
                type="radio"
                value={id}
                className="radio"
                {...register('choiceId', { required: true })}
              />
              <span className="label-text ml-4">{choice}</span>
            </label>
          </div>
        ))}

        <div className="card-actions justify-end">
          <button
            type="submit"
            disabled={formState.isSubmitting || !formState.isValid}
            className="btn btn-primary"
          >
            {formState.isSubmitting ? <InlineSpinner /> : null}
            Submit
          </button>
          <HintButton hint={question.hint} />
        </div>
      </div>
    </form>
  );
}
