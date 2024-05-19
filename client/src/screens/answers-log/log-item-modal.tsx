import { RefObject } from 'react';
import { clsx } from 'clsx';
import { type AnswersLogItem } from '@/generated/api';
import { QuestionBadges } from '@/layout';

interface Props {
  logItem: AnswersLogItem | undefined;
  onClose: () => unknown;
  modalRef: RefObject<HTMLDialogElement>;
}

export function LogItemModal({ logItem, onClose, modalRef }: Props): JSX.Element | null {
  if (!logItem) {
    return null;
  }

  const { question } = logItem;
  return (
    <dialog ref={modalRef} className="modal modal-bottom sm:modal-middle">
      <div className="modal-box">
        <h3 className="font-bold text-lg">{question.question}</h3>

        <div className="mt-4">
          <QuestionBadges topic={question.topic} difficulty={question.difficulty} />
        </div>

        <div className="mt-4">
          {question.choices.map(({ id, choice, isCorrect }) => {
            const isChecked = id === logItem.choiceId;
            return (
              <div key={id} className="form-control">
                <label className="flex justify-start cursor-pointer label">
                  <input
                    name="choice"
                    type="radio"
                    value={id}
                    className="radio"
                    disabled
                    checked={isChecked}
                  />
                  <span
                    className={clsx(
                      'label-text',
                      'ml-4',
                      isCorrect && 'text-success',
                      isChecked && !isCorrect && 'text-error',
                    )}
                  >
                    {choice}
                  </span>
                </label>
              </div>
            );
          })}
        </div>

        <p className="italic mt-4">
          <span className="font-bold">Hint: </span>
          {question.hint}
        </p>
        <p className="italic mt-4">
          <span className="font-bold">More info: </span>
          {question.moreInfo}
        </p>

        <div className="modal-action">
          <button onClick={onClose} className="btn">
            Close
          </button>
        </div>
      </div>
    </dialog>
  );
}
