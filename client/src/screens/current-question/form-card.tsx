import { Controller, useForm } from 'react-hook-form';
import startCase from 'lodash/startCase';
import { RemainingTopic, UnansweredQuestion } from '@/generated/api';
import { InlineSpinner, QuestionBadges } from '@/layout';
import { HintButton } from './hint-button';

interface Props {
  question: UnansweredQuestion;
  remainingTopics: RemainingTopic[];
  onChangeTopic: (topic: string) => unknown;
  onSubmit: ({ choiceId }: { choiceId: string }) => unknown;
}

interface Form {
  choiceId: string;
  topic: string;
}

export function QuestionFormCard({ question, remainingTopics, onChangeTopic, onSubmit }: Props): JSX.Element {
  const { control, register, handleSubmit, formState, resetField } = useForm<Form>();

  return (
    <form onSubmit={handleSubmit((data) => onSubmit(data))} className="card bg-neutral shadow-xl">
      <div className="card-body">
        <div className="flex justify-between">
          <QuestionBadges topic={question.topic} difficulty={question.difficulty} />
          <HintButton hint={question.hint} />
        </div>
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

        <div className="card-actions justify-center">
          <label className="form-control w-full sm:max-w-64">
            <Controller
              name="topic"
              control={control}
              rules={{ required: 'Topic is required' }}
              defaultValue={question.topic}
              render={({ field }) => (
                <select
                  {...field}
                  onChange={e => {
                    field.onChange(e)
                    resetField('choiceId')
                    onChangeTopic(e.target.value)
                  }}
                  className="select select-bordered"
                >
                  {remainingTopics.map(({ topic, amountOfQuestions }) => (
                    <option key={topic} value={topic}>
                      {startCase(topic)} ({amountOfQuestions})
                    </option>
                  ))}
                </select>
              )}
            />
          </label>
          
          <button
            type="submit"
            disabled={formState.isSubmitting || !formState.isValid}
            className="btn btn-primary btn-block sm:btn-wide"
          >
            {formState.isSubmitting ? <InlineSpinner /> : null}
            Submit
          </button>
        </div>
      </div>
    </form>
  );
}
