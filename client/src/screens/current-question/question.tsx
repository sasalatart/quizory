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
    handleGetNextQuestion,
    handleRefetchRemainingTopics,
  } = useCurrentQuestion();

  const { handleSubmitAnswer } = useSubmitAnswer({
    onSubmit: async (submissionFeedback) => {
      await handleRefetchRemainingTopics();
      setFeedback(submissionFeedback);
    },
  });

  if (isLoading) {
    return <CenteredSpinner />;
  }

  if (!question) {
    return <NoQuestionsLeftCard />;
  }

  if (feedback) {
    return (
      <QuestionFeedbackCard
        question={question}
        feedback={feedback}
        remainingTopics={remainingTopics ?? []}
        onNext={async (topic) => {
          await handleGetNextQuestion(topic);
          setFeedback(undefined);
        }}
      />
    );
  }

  return <QuestionFormCard question={question} onSubmit={handleSubmitAnswer} />;
}
