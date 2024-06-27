import { useState } from 'react';
import { CenteredSpinner } from '@/layout';
import { QuestionFormCard } from './form-card';
import { Feedback, QuestionFeedbackCard } from './feedback-card';
import { NoQuestionsLeftCard } from './no-questions-left-card';
import { useCurrentQuestion, useSubmitAnswer } from './hooks';

export function Question(): JSX.Element {
  const [feedback, setFeedback] = useState<Feedback>();

  const {
    question,
    remainingTopics,
    isLoading,
    handleChangeTopic,
    handleRefetchCurrentQuestion,
    handleRefetchRemainingTopics,
  } = useCurrentQuestion();

  const { handleSubmitAnswer } = useSubmitAnswer({
    question,
    onSubmit: async (submissionFeedback) => {
      await handleRefetchRemainingTopics();
      setFeedback(submissionFeedback);
    },
  });

  if (isLoading) {
    return <CenteredSpinner />;
  }

  if (feedback) {
    return (
      <QuestionFeedbackCard
        feedback={feedback}
        onNext={async () => {
          await handleRefetchCurrentQuestion();
          setFeedback(undefined);
        }}
      />
    );
  }

  if (!question) {
    return <NoQuestionsLeftCard />;
  }

  return (
    <QuestionFormCard
      question={question}
      onChangeTopic={handleChangeTopic}
      onSubmit={handleSubmitAnswer}
      remainingTopics={remainingTopics ?? []}
    />
  );
}
